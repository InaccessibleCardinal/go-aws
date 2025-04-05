dev:
	go run ./cmd/...

test:
	go test ./... -coverprofile cover.out && go tool cover -html=cover.out