FROM golang:1.16.7-alpine
WORKDIR /app/urlShortner

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /app/urlShortner/main 
RUN go build -o main.go
CMD ["./main"]
