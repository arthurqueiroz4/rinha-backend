FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build cmd/main.go

EXPOSE 8080

ENTRYPOINT [ "./main" ]
