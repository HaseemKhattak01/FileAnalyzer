FROM golang:1.22.5-alpine
WORKDIR /golang
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 1323
CMD ["./main"]