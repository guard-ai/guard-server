FROM golang:1.22-alpine

WORKDIR /app
COPY . .

ADD go.mod go.sum ./
RUN go mod download
ADD . ./

RUN go build -o main .
EXPOSE 8000

CMD ["./main"]
