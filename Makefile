# 테스트 입구는 make test 하나 (HUS ADR-0081 v1.2) — CI 는 이 타깃을 이름으로 부른다.
# 명령·상한은 여기 한 곳에만 적는다.
# 자원 상한 (ADR-0081 §2): nice + GOMEMLIMIT/GOGC + 순차(-p 1). -count=1 -race 는 CI 가 쓰던 플래그 그대로.
GO_TEST := GOMEMLIMIT=2GiB GOGC=50 nice -n 19 go test -p 1 -count=1 -race

.PHONY: test

test: ## 전체 Go 테스트 (자원 상한 · -race)
	$(GO_TEST) ./...
