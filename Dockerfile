FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o rinha-golang

EXPOSE 8080

ENTRYPOINT [ "./rinha-golang" ]
