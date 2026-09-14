#!/usr/bin/env bash
# gen-docs.sh — SDK 레퍼런스를 SSoT(godoc 주석)에서 생성한다.
#
# 출력 (전부 생성물 — 손으로 고치지 않는다, 정본은 소스 주석):
#   docs/reference/go/<pkg>.md   ← 이 모듈의 패키지 godoc (gomarkdoc). 지금은 루트 패키지 units 하나.
#
# 이 repo 에는 proto 가 없으므로 RPC 절(protoc-gen-doc · buf.gen.docs.yaml)은 두지 않는다 —
# 생기면 audit 의 scripts/gen-docs.sh 를 본으로 그 절을 붙인다.
#
# 왜 커밋하나: 팀원이 GitHub 에서 그대로 읽는다 (private repo 라 pkg.go.dev 불가, Pages 는
# team plan 에서 public). 드리프트는 CI docs-sync(검사)가 막는다 — 이 repo 는 lefthook 을
# 안 쓰므로(.githooks/pre-commit 은 시크릿 검사 전용) 로컬 자동 재생성 훅은 없다. 주석을
# 고쳤으면 커밋 전에 `scripts/gen-docs.sh` 를 손으로 돌린다.
# 생성물에 날짜·SHA 를 박지 않는다 — 박으면 매 실행이 diff 라 검사가 못 선다. "언제 것인가" 는
# git 이 답한다 (파일 history / 태그 permalink). 상세 docs/reference/README.md.
#
# 도구 핀: go.mod 에 tool 을 넣지 않는다 — units 는 zero-dependency 발행 라이브러리라 소비자
# module graph 에 문서 도구를 끌고 가지 않으려고 `go run <mod>@<ver>` 로 고정한다 (아래 한 줄이 정본).
set -euo pipefail

GOMARKDOC=github.com/princjef/gomarkdoc/cmd/gomarkdoc@v1.1.0

cd "$(dirname "$0")/.."

MOD=github.com/enclaboratory/units
OUT=docs/reference/go
MODE=${1:-generate}   # generate | check

# --- Go SDK (godoc → markdown) ---------------------------------------------------------
# 생성 패키지(protobuf 등)는 이 repo 에 없다 — 생기면 여기서 grep -v 로 뺀다 (audit 선례:
# `grep -vE '/(audit/v1|buf)(/|$)'`). 그쪽 문서는 proto 주석이 정본이라 gomarkdoc 대상이 아니다.
pkgs=$(go list ./...)
mkdir -p "$OUT"
for pkg in $pkgs; do
  rel=${pkg#"$MOD"}
  rel=${rel#/}
  if [ -z "$rel" ]; then name=$(go list -f '{{.Name}}' "$pkg"); else name=${rel//\//-}; fi
  args=(
    --repository.url https://github.com/enclaboratory/units
    --repository.default-branch main
    --repository.path /
    --output "$OUT/$name.md"
  )
  [ "$MODE" = check ] && args+=(--check)
  go run "$GOMARKDOC" "${args[@]}" "$pkg"
done
