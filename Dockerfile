FROM dhi.io/golang:1.27-alpine-dev AS builder

ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

# Set the current working directory inside the container.
WORKDIR /go/src/github.com/score-spec/score-implementation-sample

# Copy just the module bits
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build it.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w \
        -X github.com/score-spec/score-implementation-sample/internal/version.Version=${VERSION} \
        -X github.com/score-spec/score-implementation-sample/internal/version.GitCommit=${GIT_COMMIT} \
        -X github.com/score-spec/score-implementation-sample/internal/version.BuildDate=${BUILD_DATE}" \
    -o /usr/local/bin/score-implementation-sample ./cmd/score-implementation-sample

# We can use static since we don't rely on any linux libs or state, but we need ca-certificates to connect to https/oci with the init command.
FROM dhi.io/static:20260611-alpine3.24@sha256:93568eb7c673afb3ad79b15cca341469d3e02cf859caae1049aa22fe7fbce90a

# Set the current working directory inside the container.
WORKDIR /score-implementation-sample

# Copy the binary from the builder image.
COPY --from=builder /usr/local/bin/score-implementation-sample /usr/local/bin/score-implementation-sample

# Run the binary.
ENTRYPOINT ["/usr/local/bin/score-implementation-sample"]
