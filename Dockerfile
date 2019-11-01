FROM docker.io/golang:latest as golang

WORKDIR /build

COPY . .

RUN mkdir -p /out

# Cache builds without version info
RUN go build -mod readonly -o /out/veidemann-reset -ldflags "-s -w"

ARG VERSION
RUN go build -mod readonly -o /out/veidemann-reset \
    -ldflags "-s -w -X github.com/nlnwa/veidemann-reset/pkg/version.Version=${VERSION:-$(git describe --tags --always)}"


FROM gcr.io/distroless/base

COPY LICENSE /LICENSE

COPY --from=golang /out /out

ENTRYPOINT ["/out/veidemann-reset"]
