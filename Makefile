MAIN_PATH = cmd/main.go
run:
	go run $(MAIN_PATH)

tidy:
	go mod tidy
