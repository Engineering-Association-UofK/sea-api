# syntax=docker/dockerfile:1

################################################################################

ARG GO_VERSION=1.26
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/api

################################################################################

FROM debian:bookworm-slim AS final

# Install Inkscape, fontconfig, and core certificates/tzdata
RUN apt-get update && apt-get install -y --no-install-recommends \
    inkscape \
    fontconfig \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY ./resources/fonts /usr/local/share/fonts/custom-fonts

RUN fc-cache -fv

ARG UID=10001
RUN useradd --uid ${UID} --create-home appuser
USER appuser


COPY --from=build /bin/server /bin/
COPY ./db/migrations /app/db/migrations
COPY ./resources/static-assets /app/resources/static-assets
EXPOSE 8000


ENTRYPOINT [ "/bin/server" ]
