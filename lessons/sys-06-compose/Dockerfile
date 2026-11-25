# 1-bosqich: yigʻish (SYS-0.05 dagi Dockerfile, endi modullar bilan)
FROM golang:1.26.8 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY server ./server
RUN CGO_ENABLED=0 go build -o /server ./server

# 2-bosqich: boʻm-boʻsh image, faqat binary
FROM scratch
COPY --from=build /server /server
USER 65532:65532
EXPOSE 8077
CMD ["/server"]
