# TUZOQ: golang image'da gcc bor, shuning uchun cgo yoqiladi va net paketi libc'ga bogʻlanadi
FROM golang:1.26.8 AS build
WORKDIR /src
COPY go.mod main.go ./
RUN go build -o /api .

FROM scratch
COPY --from=build /api /api
ENTRYPOINT ["/api"]
