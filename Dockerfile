FROM golang:1.15 as modules

ADD go.mod go.sum /m/
WORKDIR /m
RUN go mod download

# builder

FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /app
ADD . /app
WORKDIR /app

RUN make test

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/api ${PROJECT}

RUN useradd -u 10001 myapp

# runner

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd

USER myapp

COPY --from=builder /app/bin/api /api

CMD ["/api"]