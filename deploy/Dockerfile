FROM golang:1.16 as builder
WORKDIR /src/
COPY . .
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o adapter .

FROM gcr.io/distroless/static
COPY --from=builder /src/adapter /
EXPOSE 443
ENTRYPOINT ["/adapter"]