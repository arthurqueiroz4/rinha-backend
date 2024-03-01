FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o rinha ./*.go

EXPOSE 8080

ENTRYPOINT [ "./rinha" ]
