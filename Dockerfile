FROM docker.io/golang:1.19 as golang

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Cache build without version info
# -trimpath remove file system paths from executable
# -ldflags arguments passed to go tool link:
#   -s disable symbol table
#   -w disable DWARF generation
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod readonly -trimpath -ldflags "-s -w"

ARG VERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod readonly -trimpath \
    -ldflags "-s -w -X github.com/nlnwa/veidemann-reset/internal/version.Version=${VERSION:-$(git describe --tags --always)}"


FROM gcr.io/distroless/base-debian11

COPY LICENSE /LICENSE

COPY --from=golang /build/veidemann-reset /veidemann-reset

ENTRYPOINT ["/veidemann-reset"]
