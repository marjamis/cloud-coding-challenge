FROM alpine:3.5

RUN apk --no-cache --update upgrade && \
    apk add --no-cache ca-certificates curl && \
    update-ca-certificates

HEALTHCHECK --interval=30s --timeout=30s --retries=3 \
  CMD curl -fo /dev/null http://localhost:8081/healthcheck || exit 1

WORKDIR /
USER nobody

COPY ./bin/metis.go /entrypoint
COPY ./templates/ /templates/

ENTRYPOINT [ "/entrypoint" ]
