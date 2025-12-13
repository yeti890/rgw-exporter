# ---------- build stage ----------
FROM golang:1.25.5-alpine3.23 AS build
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy && \
    CGO_ENABLED=0 GOFLAGS="-trimpath" go build -ldflags "-s -w" -o /out/rgw-exporter

# ---------- runtime stage (alpine) ----------
FROM alpine
WORKDIR /app
COPY --from=build /out/rgw-exporter /app/rgw-exporter
EXPOSE 9240
ENV GODEBUG=madvdontneed=1
ENTRYPOINT ["/app/rgw-exporter"]
