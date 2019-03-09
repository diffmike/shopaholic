FROM golang:1.12

WORKDIR /go/src/app
COPY . .

EXPOSE 8080

ENV GO111MODULE=on

COPY ./go.mod ./go.sum ./
RUN go mod download

RUN go build \
    -installsuffix 'static' \
    -o /app .

CMD ["/app", "bot:start"]
