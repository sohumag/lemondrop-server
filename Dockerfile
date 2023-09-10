
FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY  . .
RUN go mod download

ENV MONGO_URI="mongodb+srv://agrawalsohum:KgGM2cF55EBXLmih@cluster0.jvkmpsh.mongodb.net/?retryWrites=true&w=majority"

RUN go build -o bin/ ./...
EXPOSE 8080 80


# Run
CMD ["./bin/cmd"]