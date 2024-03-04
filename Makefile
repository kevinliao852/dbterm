clean:
				@find . -type f -name dbterm | xargs rm -f

build:
				@go build -ldflags "-s -w" -o dbterm cmd/main.go
