#!/usr/bin/env bash
# 훅 설치 — `core.hooksPath` 를 이 디렉토리로 돌린다.
#
# 왜 스크립트인가: `.git/hooks/` 는 저장소에 안 실린다. 클론마다 손으로 켜라고 하면 결국
# 안 켜진다 — 그래서 `npm install` 의 postinstall 이 부른다.
#
# ⚠️ **덮어쓰지 않는다.** 다른 값이 이미 설정돼 있으면 그 사실을 말하고 물러난다. 이 트리를
# 여러 세션이 공유하고, 남이 정한 hooksPath 를 조용히 바꾸면 그쪽 훅이 사라진다.
set -euo pipefail
cd "$(dirname "$0")/.."
cur=$(git config --local core.hooksPath || true)
want=.githooks
if [ "$cur" = "$want" ]; then exit 0; fi
if [ -n "$cur" ]; then
  echo "[hooks] core.hooksPath 가 이미 '$cur' 이다 — 건드리지 않는다. 필요하면 직접: git config core.hooksPath $want" >&2
  exit 0
fi
git config --local core.hooksPath "$want"
echo "[hooks] core.hooksPath=$want (시크릿 커밋 차단 활성)"
