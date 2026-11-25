# Uy vazifasi javoblari uchun: homework/$VAZIFA paketini yigʻadi (Qiyin, Guru)
FROM golang:1.26.8 AS build
ARG VAZIFA
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY homework ./homework
RUN CGO_ENABLED=0 go build -o /server ./homework/${VAZIFA}

FROM scratch
COPY --from=build /server /server
USER 65532:65532
CMD ["/server"]
