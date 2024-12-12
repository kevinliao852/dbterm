help: ## Show help message
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | column -s "##" -t

clean: ## Clean the built binary
	@find . -type f -name dbterm | xargs rm -f

build: ## Build the main binary
	@go build -ldflags "-s -w" -o dbterm cmd/main.go
