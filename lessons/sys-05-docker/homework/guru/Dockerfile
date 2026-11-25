# Guru javobi: server root emas, 65532 foydalanuvchisi nomidan ishlaydi.
# Kontekst: dars papkasi (docker build -f homework/guru/Dockerfile .)
FROM golang:1.26.8 AS build
WORKDIR /src
COPY go.mod ./
COPY server ./server
RUN CGO_ENABLED=0 go build -o /server ./server
# scratch'da mkdir yoʻq: papkani shu bosqichda yaratamiz
RUN mkdir /data

FROM scratch
COPY --from=build /server /server
# /data egasi 65532: yangi (boʻsh) volume ulanganda shu egalikni oladi
COPY --from=build --chown=65532:65532 /data /data
ENV DATA=/data
USER 65532:65532
EXPOSE 8077
CMD ["/server"]
