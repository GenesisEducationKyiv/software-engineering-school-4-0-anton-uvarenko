FROM golang:1.22.2 as builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o ./app ./cmd/app/main.go

FROM ubuntu:oracular
RUN apt-get install libc6
COPY --from=builder /app/app /home/app
CMD /home/app 
