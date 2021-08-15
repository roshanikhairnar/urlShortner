FROM golang:1.16.7-alpine
WORKDIR /app/urlShortner

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

RUN go build -o main.go

