# Multi-stage build on Debian trixie-slim (date-pinned for idempotent
# builds), as documented.
# Build: docker build --build-arg VERSION=1.0.0 -t ycast-go .
# Run:   docker run -d -p 80:8080 -v ./stations.yml:/etc/ycast/stations.yml:ro ycast-go

# ---- build ----------------------------------------------------------------
# Digest pins the manifest list, so multi-arch builds resolve both arches.
FROM debian:trixie-20261005-slim@sha256:a29215f6a35e51e22adffa17f89e9d2ef06214e64a2bad10d765c46aea49f11f AS build
ARG VERSION=dev
RUN apt-get update \
    && apt-get install -y --no-install-recommends golang-any ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN go build -mod=vendor -trimpath \
    -ldflags "-X main.version=${VERSION} -s -w" \
    -o /usr/local/bin/ycast-go ./cmd/ycast-go

# ---- runtime ---------------------------------------------------------------
FROM debian:trixie-20261005-slim@sha256:a29215f6a35e51e22adffa17f89e9d2ef06214e64a2bad10d765c46aea49f11f
LABEL org.opencontainers.image.title="ycast-go" \
      org.opencontainers.image.description="Self-hosted vTuner replacement for network receivers" \
      org.opencontainers.image.source="https://github.com/hgomez/ycast-go" \
      org.opencontainers.image.licenses="GPL-3.0-or-later"

COPY --from=build /usr/local/bin/ycast-go /usr/bin/ycast-go
COPY stations.yml.example /etc/ycast/stations.yml
COPY LICENSE /usr/share/licenses/ycast-go/LICENSE

# Unprivileged runtime: the container maps receiver port 80 to 8080.
RUN useradd --system --no-create-home --shell /usr/sbin/nologin ycast \
    && mkdir -p /var/cache/ycast-go/icons \
    && chown -R ycast:ycast /var/cache/ycast-go
USER ycast

EXPOSE 8080
VOLUME ["/var/cache/ycast-go"]

ENTRYPOINT ["/usr/bin/ycast-go"]
CMD ["-c", "/etc/ycast/stations.yml", "-l", "0.0.0.0", "-p", "8080", "-cache", "/var/cache/ycast-go/icons"]

# trixie-slim has no curl/wget: use bash's /dev/tcp for the health check.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["CMD-SHELL", "exec 3<>/dev/tcp/127.0.0.1/${PORT:-8080} && printf 'GET /healthz HTTP/1.0\\r\\n\\r\\n' >&3 && grep -q ok <&3"]
