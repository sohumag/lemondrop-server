aws_binary:
	GOOS=linux GOARCH=arm64 go build -o bin/aws/ ./...

	scp -i ../"ssh-key-pair.pem" bin/aws/cmd ec2-user@ec2-54-153-101-157.us-west-1.compute.amazonaws.com:~/app/bin

binary:
	GOOS=darwin GOARCH=arm64 go build -o bin/main cmd/main.go

run_binary:
	go build -o bin ./...
	./bin/cmd