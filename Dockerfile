FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o social ./cmd


FROM alpine

WORKDIR /root/

COPY --from=builder /build/social .
COPY --from=builder /build/config.yaml .
COPY --from=builder /build/wait-for-postgres.sh .

RUN apk update
RUN apk --no-cache add postgresql-client

RUN chmod +x wait-for-postgres.sh 

RUN pwd

CMD ["./social"]

EXPOSE 8080
