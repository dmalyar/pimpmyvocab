FROM golang:1.14.2-alpine as builder
LABEL stage=builder
WORKDIR /pimpmyvocab
COPY . .
RUN CGO_ENABLED=0 go test ./...
RUN go build -o ./out/pmv_bot ./cmd/pmv_bot.go

FROM alpine:3.11.6
RUN apk add --no-cache bash
WORKDIR /pimpmyvocab
COPY --from=builder ["/pimpmyvocab/out/pmv_bot", "/pimpmyvocab/config.yaml", "/pimpmyvocab/wait-for-it.sh", "./"]
COPY --from=builder /pimpmyvocab/repo/migration db/migration
# CMD ["pmv_bot"] # Uncomment for using without docker-compose