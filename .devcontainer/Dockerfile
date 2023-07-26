FROM kjconroy/sqlc:latest AS sqlc

FROM qmcgaw/godevcontainer:alpine

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH \
    CGO_ENABLED=1 \
    GO111MODULE=on

COPY --from=sqlc /workspace/sqlc /usr/local/bin/sqlc
