FROM golang:latest as build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd

FROM scratch

COPY --from=build app/main main

EXPOSE 80

CMD [ "./main" ]