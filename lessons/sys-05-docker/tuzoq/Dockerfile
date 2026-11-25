# PID 1 tuzogʻi: server shell skript ichidan ishga tushadi
FROM golang:1.26.8 AS build
WORKDIR /src
COPY go.mod ./
COPY server ./server
RUN CGO_ENABLED=0 go build -o /server ./server

FROM alpine:3.24
COPY --from=build /server /server
COPY tuzoq/start.sh tuzoq/start-exec.sh /
CMD ["/start.sh"]
