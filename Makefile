.PHONY=run
run: fmt
	go run cmd/main.go
	
.PHONY=fmt
fmt: 
	go fmt ./...

