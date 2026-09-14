# units — SDK 레퍼런스

| 문서 | 내용 |
|---|---|
| [go/units.md](go/units.md) | `Power`(내부 W) · `Energy`(내부 kWh) — 단위 생성자 `W`/`KW`/`MW`/`Wh`/`KWh`/`MWh`, 접근자, 차원 교차 `Over`/`Per`, `Validate` · `ErrNaN`/`ErrInfinite` |

RPC 는 없다 — 이 repo 는 서비스가 아니라 zero-dependency Go 라이브러리다. fleet 전체 진입점은
[HUS docs/api-map.md](https://github.com/enclaboratory/HUS/blob/main/docs/api-map.md).

## 부르는 법

**Go** — `go get github.com/enclaboratory/units` (private: `GOPRIVATE=github.com/enclaboratory/*`).

```go
cap := units.MW(1.5)                       // 값은 단위 생성자로 들어온다
e := units.KW(500).Over(2 * time.Hour)     // P × t → Energy (1 MWh)
if err := units.KWh(payload.Supply).Validate(); err != nil { /* errors.Is(err, units.ErrNaN) */ }
```

`Power` 와 `Energy` 는 서로 다른 정의 타입이라 섞으면 컴파일 에러다 — 차원을 넘는 연산은 `Over`/`Per` 뿐이다.
돈·단가·표시 포맷은 이 패키지 밖이다 ([go/units.md](go/units.md) 의 Non-goals).

## 이 문서는 언제 것인가

`main` 의 문서는 `main` 의 코드와 같다 — 소스 주석에서 생성하고 CI 가 어긋남을 막는다.
특정 발행분은 태그로 연다: `https://github.com/enclaboratory/units/blob/<vX.Y.Z>/docs/reference/go/units.md`
(태그 목록 [tags](https://github.com/enclaboratory/units/tags); 이 문서는 v0.2.0 다음 태그부터 실린다). 문서 안 소스 링크는 `main` 기준이다.

## 고치려면

여기가 아니라 주석을 고친다 — godoc 은 루트의 `*.go` (패키지 설명은 `doc.go`).
이 repo 는 pre-commit 자동 재생성이 없으므로 커밋 전에 `scripts/gen-docs.sh` 를 돌려 같이 싣는다.
검사는 `scripts/gen-docs.sh check` (CI docs-sync 와 동일; 절차·도구 핀은 그 스크립트).
