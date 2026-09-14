#!/usr/bin/env bash
# 시크릿이 커밋에 들어가는 것을 막는다.
#
# lefthook pre-commit 의 `secret-scan` command 가 부른다 (2026-09-14 — 그전엔 .githooks/pre-commit 으로
# core.hooksPath 를 잡았는데, 그러면 git 이 .git/hooks 의 lefthook 을 통째로 무시해 gofmt·test 훅이 죽었다.
# 훅 체계는 lefthook 하나다: `brew install lefthook && lefthook install`).
#
# ## 왜 GitHub push protection 이 아니라 여기인가
#
# GitHub 쪽은 **푸시**를 막는다 — 그때는 이미 로컬 히스토리에 있고, 지우려면 히스토리를
# 고쳐야 한다. 커밋에서 막으면 애초에 안 들어간다. 그리고 이쪽은 플랜·조직 권한이 필요
# 없다(비공개 저장소의 시크릿 스캔은 유료 제품이다).
#
# 둘은 대체가 아니라 층이다. GitHub 쪽도 켜는 것이 맞다 — 훅은 `--no-verify` 로 넘길 수
# 있고 남의 클론에는 설치돼 있지 않을 수 있다.
#
# ## 우회
#
#   SKIP_SECRET_CHECK=1 git commit …
#
# **옵트아웃이다.** 옵트인은 그 값을 아는 사람만 보호하고 이 저장소는 그걸로 두 번 데었다
# (TEST_PG_REQUIRED). 넘기려면 적어야 하고, 적는 사람은 대가를 안다.
#
# ## ⚠️ 자기가 못 돌면 통과시키지 않는다
#
# 처음 판은 `grep -v '^\+\+\+'` 로 추가 줄을 골랐다. 이 머신의 `grep` 은 **ugrep** 이라
# 그 이스케이프에서 죽었고, 그러면 대상 목록이 비어 **훅이 조용히 통과했다**(2026-09-13
# 실측 — 가짜 AWS 키가 실제로 커밋됐다). 검사가 "통과" 가 아니라 "못 돌았다" 를 말해야
# 한다. 그래서 추출은 awk 로 하고, 추출이 비었는데 diff 는 있으면 **거부**한다.
set -uo pipefail

[ -n "${SKIP_SECRET_CHECK:-}" ] && exit 0

staged=$(git diff --cached --name-only --diff-filter=ACMR)
[ -z "$staged" ] && exit 0

fail=0
say() { printf '  %s\n' "$*" >&2; }
head_() { if [ $fail -eq 0 ]; then say ""; say "⛔ 커밋 거부 — 시크릿으로 보이는 것이 있다"; fi; fail=1; }

# ── 1. 파일 이름 ─────────────────────────────────────────────
# 내용과 무관하게 들어가면 안 되는 것들. `.env.example` 은 값이 비어 있어야 하는 견본이라 뺀다.
# `.env.<x>.example` 도 견본이다 (common 미러 `.env.common.example`) — 아래 `.env.*` 가 그것까지
# 삼켜서 2026-09-14 에 미러 수정 커밋이 막혔다. 견본 판정을 먼저 둔다.
while IFS= read -r f; do
  b=${f##*/}
  case "$b" in
    .env.example|.env.sample|*.env.example|*.env.sample|.env.*.example|.env.*.sample) continue ;;
    .env|.env.*|*.pem|*.key|*.p12|*.pfx|*.jks|id_rsa*|id_ed25519*) ;;
    *) continue ;;
  esac
  head_; say "  파일 이름: $f"
done <<< "$staged"

# ── 2. 내용 ─────────────────────────────────────────────────
# 추가된 줄만 본다. **awk 로 뽑는다** — grep 구현체마다 `\+` 해석이 달라서다(머리말 참조).
diffout=$(git diff --cached -U0 --diff-filter=ACMR || true)
added=$(printf '%s\n' "$diffout" | awk 'substr($0,1,1)=="+" && substr($0,1,3)!="+++"')

if [ -n "$diffout" ] && [ -z "$added" ]; then
  head_; say "  검사기가 추가 줄을 하나도 못 뽑았다 — 통과가 아니라 **고장**이다"
fi

# **제공자 접두가 붙은 것만** 본다. 느슨하면 오탐이 늘고, 오탐이 늘면 사람이 훅을 끈다 —
# 꺼진 훅은 없는 훅이다.
patterns='(ghp_|gho_|ghu_|ghs_|ghr_|github_pat_)[A-Za-z0-9_]{20,}
AKIA[0-9A-Z]{16}
xox[baprs]-[A-Za-z0-9-]{10,}
sk-[A-Za-z0-9]{20,}
BEGIN [A-Z ]*PRIVATE KEY
AIza[0-9A-Za-z_-]{30,}'

if [ -n "$added" ]; then
  while IFS= read -r p; do
    [ -z "$p" ] && continue
    if printf '%s\n' "$added" | grep -qE -- "$p"; then
      head_; say "  패턴에 걸림: $p"
      # **값은 찍지 않는다.** 진단하려고 시크릿을 터미널·스크롤백에 남기면 막으려던 것을
      # 스스로 한다. 어느 패턴인지까지만 말한다.
    fi
  done <<< "$patterns"
fi

if [ $fail -ne 0 ]; then
  say ""
  say "값은 일부러 출력하지 않는다 — 여기 찍으면 막으려던 것을 스스로 하게 된다."
  say "정말 의도한 것이면: SKIP_SECRET_CHECK=1 git commit …"
  say ""
  exit 1
fi
exit 0
