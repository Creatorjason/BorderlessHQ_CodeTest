FROM golang:1.18.3


WORKDIR /app


COPY . .

RUN go mod tidy
RUN go build -o myapp


EXPOSE 9091


CMD ["./myapp"]
