FROM golang:alpine as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o main .

FROM alpine

RUN adduser -D -H -S appuser /app

USER appuser

COPY . /go/src/myapp

COPY --from=builder /build/main /app/

ENV PATH="/app:${PATH}"

WORKDIR /app

EXPOSE 3000

CMD ["./main"]
