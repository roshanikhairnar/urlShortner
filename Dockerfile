FROM golang:1.16.7-alpine
WORKDIR /app/urlShortner
EXPOSE 8080
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main.go
RUN ["chmod", "+x", "./main.go"]

CMD ["./main.go"]

#netstat -tulpn
