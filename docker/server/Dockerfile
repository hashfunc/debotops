FROM golang:1.17 as builder

WORKDIR /workspace

COPY ../../go.mod go.mod
COPY ../../go.sum go.sum

RUN go mod download

COPY ../../cmd/server/main.go main.go
COPY ../../api api/
COPY ../../pkg pkg/
COPY ../../server server/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o out main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/out ./server
USER 65532:65532

ENTRYPOINT ["/server"]
