FROM golang:1.22.3

RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN go build -o main ./server
RUN chmod +x /app/main

CMD ["./main"]