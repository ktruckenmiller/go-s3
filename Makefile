exec_name=s3

build-linux:
	docker run -it --rm \
	-v ${PWD}:/app \
	-w /app \
	-e GOOS=linux \
	-e GOARCH=386 \
	-e CGO_ENABLED=0 \
 	golang \
	go build -o $(exec_name) main.go

build-mac:
	docker run -it --rm \
	-v ${PWD}:/app \
	-w /app \
	-e GOOS=darwin \
	-e GOARCH=amd64 \
 	golang \
	go build -o $(exec_name) main.go
