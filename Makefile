exec_name=s3

build-linux:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(exec_name) main.go

build-mac:
	go build -o $(exec_name) main.go
