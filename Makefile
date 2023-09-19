aws_binary:
	GOOS=linux GOARCH=arm64 go build -o bin/aws/ ./...

	scp -i ../"ssh-key-pair.pem" bin/aws/cmd ec2-user@ec2-54-153-101-157.us-west-1.compute.amazonaws.com:~/app/bin

binary:
	GOOS=darwin GOARCH=arm64 go build -o bin/ ./...

run_binary:
	go build -o bin ./...
	./bin/cmd

docker:
	docker build -t lemondrop .
	docker run -p 8080:8080 -p 80:80 lemondrop

run:
	go build -o bin ./...
	./bin/cmd
