FROM golang:1.23 AS build

WORKDIR /simplesearch

COPY . .

RUN go mod download
RUN go build -o ./build/simplesearch ./cmd/simplesearch/main.go && \
    go build -o ./build/esmigrator ./cmd/esmigrator/main.go

FROM gcr.io/distroless/base-debian12

LABEL org.opencontainers.image.source https://github.com/xoticdsign/simplesearch
LABEL org.opencontainers.image.description SimpleSearch Service Image
LABEL org.opencontainers.image.licenses MIT

WORKDIR /simplesearch

COPY --from=build simplesearch/build/simplesearch .
COPY --from=build simplesearch/build/esmigrator tools/
COPY --from=build simplesearch/config config/
COPY --from=build simplesearch/migrations migrations/

EXPOSE 8080

ENTRYPOINT [ "./simplesearch" ]