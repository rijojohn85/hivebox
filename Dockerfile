FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o hivebox-binary

CMD ["./hivebox-binary"]
