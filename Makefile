aws_binary:
	GOOS=linux GOARCH=arm64 go build -o bin/aws/main cmd/main.go

	scp -i ../"ssh-key-pair.pem" bin/aws/main ec2-user@ec2-54-193-64-134.us-west-1.compute.amazonaws.com:~/app/bin

binary:
	GOOS=darwin GOARCH=arm64 go build -o bin/main cmd/main.go

run_binary:
	go build -o bin ./...
	./bin/cmd