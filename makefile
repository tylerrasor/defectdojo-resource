test:
	go test ./... -failfast -cover

coverage: cover
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
