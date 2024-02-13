FROM golang

WORKDIR /app

COPY /src .
COPY go.mod .
COPY go.sum .

RUN go mod tidy

RUN go build -o rinha ./*.go

EXPOSE 8080

ENTRYPOINT [ "./rinha" ]
