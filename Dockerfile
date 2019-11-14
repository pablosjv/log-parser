FROM golang:1.12 AS builder

ARG VERSION
LABEL VERSION=$VERSION

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /go/bin/log-parser


FROM scratch
COPY --from=builder /go/bin/log-parser /log-parser
ENTRYPOINT ["/log-parser"]
