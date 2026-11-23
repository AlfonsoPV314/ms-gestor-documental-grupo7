# syntax=docker/dockerfile:1.7
FROM node:22.22.0-alpine3.23@sha256:e4bf2a82ad0a4037d28035ae71529873c069b13eb0455466ae0bc13363826e34 AS swagger
WORKDIR /src
RUN npm install --global pnpm@10.13.1 --ignore-scripts
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml .npmrc ./
COPY apps/web/package.json apps/web/package.json
RUN --mount=type=cache,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile --ignore-scripts --filter ms-gestor-documental
FROM --platform=$BUILDPLATFORM golang:1.26.11-alpine3.23@sha256:23f198931d4bd83b4f5a613bd75656257ae8bf832dcbc2db80e1050b5219ca2a AS build
WORKDIR /src
COPY apps/api/go.mod apps/api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY apps/api/ ./
ARG TARGETOS TARGETARCH VERSION=dev COMMAND=api
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -buildvcs=false -ldflags="-s -w -X main.version=$VERSION" -o /out/servicio ./cmd/$COMMAND
FROM gcr.io/distroless/static-debian13:nonroot@sha256:f7f8f729987ad0fdf6b05eeeae94b26e6a0f613bdf46feea7fc40f7bd72953e6
COPY --from=build --chown=nonroot:nonroot /out/servicio /app/servicio
COPY --from=swagger --chown=nonroot:nonroot /src/node_modules/swagger-ui-dist/swagger-ui.css /src/node_modules/swagger-ui-dist/swagger-ui-bundle.js /app/swagger/
ENV SWAGGER_ASSETS_DIR=/app/swagger
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/servicio"]


