FROM golang:1.20-alpine3.18

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

ENTRYPOINT ["tail", "-f", "/dev/null"]