.PHONY: test cover mock

# 跑全案測試
test:
	go test -v -cover ./...

# 跑特定層級測試
test-svc:
	go test -v ./internal/services/...

# 產生測試覆蓋率報告並開啟網頁
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# 重新生成 Mock
mock:
	mockery --all --recursive --keeptree