FROM golang:1.16-alpine as builder

RUN mkdir -p /backend/

WORKDIR /backend

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -o bin/main svc/backend/*.go

FROM alpine:3.13

LABEL maintainer="trconley"

RUN addgroup -S app \
    && adduser -S -G app app \
    && apk --no-cache add \
    ca-certificates curl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /backend/bin/main .
RUN chown -R app:app ./

USER app

CMD ["./main"]/
