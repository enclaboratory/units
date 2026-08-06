# units — Claude 작업 가이드

## 표준 (HUS)

HUS v2.5 를 따른다. 본문 + ADR SSoT (복사하지 않고 참조만):

- 본문: [hakyoung_universal_stack_v2.md](https://github.com/enclaboratory/HUS/blob/main/hakyoung_universal_stack_v2.md) — 특히 §6.1 (SSoT) · §11, 이 repo 의 결정 근거 (사용자 결정 B, 2026-08-06)
- ADR: [adr/](https://github.com/enclaboratory/HUS/tree/main/adr)
- 🗺 **문서 지도**: [HUS/docs/document-map.md](https://github.com/enclaboratory/HUS/blob/main/docs/document-map.md) — 무엇을 읽을지(세션 리딩패스) · 규칙 검증 정본 매트릭스 · ADR 주제군

## 이 repo 의 성격

- **zero-dependency Go 라이브러리** (stdlib 만) — 물리량(Power/Energy)의 컴파일 타임 단위 검증.
- 내부 표현 고정: Power=W, Energy=Wh. 차원 교차는 `Over`/`Per` 명시 연산만.
- **비목표는 doc.go §Non-goals 가 정본** — 돈(→ billing `domain.Money` 정수 minor)·단가(₩/kWh 합성)·자동 단위 표시는 여기 추가 금지.
- 의존성 추가 금지 — leaf 로 유지 (go-errors 와 동일 원칙).

## coordination

세션 시작 시 `HUS/coordination/bin/husctl inbox units`로 확인. 이 repo 는 passive RP (ADR-0044) — cron/자율 loop 없음.

## 세션 시작 시 자율 절차 (ADR-0044 passive + ADR-0071 gateway)

0. **repo-lock**: 작업 시작 시 `husctl repo-lock acquire units --owner <session-id>` (TTL 90분 — 장기 세션은 만료 전 재acquire). 완료/종료 시 `release`. acquire 실패(exit 8) = 타 세션 작업 중 — 같은 repo 동시 수정 금지 (2026-08-05, orchestrator 디스패치 체크와 대칭).
1. `/Users/hakyoung/workspace/HUS/coordination/bin/husctl inbox units` — ledger+legacy 티켓과 remote sync 상태 확인.
2. ledger 티켓 작업 시작: `husctl ticket claim <ticket-id> --owner <session-id>`.
3. 완료: `husctl ticket done <ticket-id> --owner <session-id>`. 발신: `husctl ticket send --from units --to <rp> --file <path>`.

- `request_id: null` legacy 티켓은 read-only로 surface하고 orchestrator migration 전 직접 archive하지 않는다.
- ⛔ HUS working tree/index 직접 수정·commit·rebase 금지. HUS write는 `husctl`, remote sync는 orchestrator 채널.
- ⛔ **self-poll cron 등록 금지** (2026-07-22 폐기, ADR-0044 passive). inbox 주기 감시는 orchestrator 단독.
