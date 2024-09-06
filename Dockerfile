FROM golang:1.22.4 AS builder
WORKDIR /build

COPY go.mod .
COPY go.sum .
ENV GOPROXY="https://goproxy.cn,direct"
RUN go mod download

COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" export COMMIT=$(git rev-parse --short HEAD) && \
    export BUILD_TAG=$(git describe --tags --abbrev=6 | sed 's/-/_/g') && \
    export DATETIME=$(date -u '+%Z-%Y-%m-%d_%I:%M:%S') && \
    export FLAGS="-X 'main.buildTime=${DATETIME}' -X 'main.gitCommitID=${COMMIT}' -X 'main.buildTag=${BUILD_TAG}'" && \
    export CGO_ENABLED=0 GOOS=linux GOARCH=amd64 && \
    go build -ldflags "${FLAGS}" ./cmd/webhook


FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.source=https://github.com/Yamu-OSS/external-dns-yamu-webhook
LABEL org.opencontainers.image.description="external-dns-yamu-webhook"
LABEL org.opencontainers.image.licenses=Apache-2.0

USER 8675:8675
COPY --from=builder --chmod=555 /build/webhook /external-dns-yamu-webhook
EXPOSE 8888/tcp
ENTRYPOINT ["/external-dns-yamu-webhook"]
