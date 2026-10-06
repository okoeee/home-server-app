FROM golang:1.27.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# 実行環境を空のイメージから作成するため、Cgoを使用せずにビルド
RUN CGO_ENABLED=0 go build -o /server .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /server /server

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/server"]
