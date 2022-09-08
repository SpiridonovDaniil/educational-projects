FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /opt/proxy ./cmd/proxy/proxy.go

CMD ["/opt/proxy"]