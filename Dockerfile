FROM golang:1.18.3


WORKDIR /app


COPY . .

RUN go mod tidy
RUN go build -o borderlessHQ_service


EXPOSE 9091


CMD ["./borderlessHQ_service"]
