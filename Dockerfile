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

FROM alpine:latest AS final

# install weasyprint dependancy
RUN apk add --no-cache \
    weasyprint \
    cairo \
    pango \
    gdk-pixbuf \
    ttf-dejavu \
    fontconfig

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ARG UID=10001
RUN mkdir -p /home/appuser

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/appuser" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

RUN chown -R appuser:appuser /home/appuser

USER appuser


COPY --from=build /bin/server /bin/


EXPOSE 8000


ENTRYPOINT [ "/bin/server" ]
