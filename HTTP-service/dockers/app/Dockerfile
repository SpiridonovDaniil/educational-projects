FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /opt/app ./cmd/firstapp/firstapp.go

CMD ["/opt/app"]