FROM golang:latest

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd

EXPOSE 80

CMD [ "./main" ]