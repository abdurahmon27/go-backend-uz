# 1-bosqich: yigʻish, Go kompilyatori shu yerda qoladi
FROM golang:1.26.8 AS build
WORKDIR /src
COPY go.mod ./
COPY server ./server
RUN CGO_ENABLED=0 go build -o /server ./server

# 2-bosqich: boʻm-boʻsh image, faqat binary
FROM scratch
COPY --from=build /server /server
ENV DATA=/data
EXPOSE 8077
CMD ["/server"]
