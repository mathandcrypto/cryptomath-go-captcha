ARG IMG_GO=golang:1.17-alpine
ARG IMG_ALPINE=alpine:3.15

FROM $IMG_GO AS builder

RUN apk add --no-cache --update git \
    build-base

WORKDIR /app

COPY ./go.mod ./go.sum ./Makefile ./

RUN make install-deps

COPY ./ .

RUN make vendor && make build-seed && make build-migrate

FROM $IMG_ALPINE as image-captcha-seed

WORKDIR /app

COPY --from=builder /app/out/bin/cryptomath-captcha-seed ./

VOLUME ["/app/configs"]

RUN chmod +x /app/cryptomath-captcha-seed

ENTRYPOINT ["./cryptomath-captcha-seed"]

FROM $IMG_ALPINE as image-captcha-migrate

WORKDIR /app

COPY --from=builder /app/out/bin/cryptomath-captcha-migrate ./
COPY --from=builder /app/migrations ./

RUN chmod +x /app/cryptomath-captcha-migrate

ENTRYPOINT ["./cryptomath-captcha-migrate"]