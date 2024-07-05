.PHONY: build-daemon
build-best-record:
	@echo "==> Running build command..."
	cd server && GOOS=windows GOARCH=amd64 go build -o best-record.exe ./main.go