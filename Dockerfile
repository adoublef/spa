# (i.e. 'cgr.dev', 'docker.io')
ARG REGISTRY=docker.io

FROM --platform=${BUILDPLATFORM} ${REGISTRY}/chainguard/go AS go
WORKDIR /src

ENV CGO_ENABLED=1

ARG TARGETOS TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

COPY go.* .
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go mod download

FROM go AS build

COPY . .
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build \
    -tags=osusergo,netgo,timetzdata \
    -ldflags="-s -w -extldflags=-static" \
    -o=/usr/local/bin/a.out ./cmd/spa/

FROM ${REGISTRY}/chainguard/cc-dynamic

COPY --from=build /usr/local/bin /usr/local/bin

ENTRYPOINT [ "a.out" ]
CMD [ "serve" ]