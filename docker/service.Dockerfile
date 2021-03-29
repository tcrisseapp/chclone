FROM golang:1.16-alpine as builder

ARG SERVICE

RUN mkdir -p /${SERVICE}/

WORKDIR /${SERVICE}

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -o bin/main cmd/${SERVICE}/*.go

FROM alpine:3.13

ARG SERVICE

LABEL maintainer="TRConley"

RUN addgroup -S app \
    && adduser -S -G app app \
    && apk --no-cache add \
    ca-certificates curl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /${SERVICE}/bin/main .

RUN chown -R app:app ./

USER app

CMD ["./main"]
