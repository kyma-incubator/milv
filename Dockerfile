FROM golang:1.17.9-alpine3.15 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/milv
WORKDIR ${BASE_APP_DIR}

COPY ./ ${BASE_APP_DIR}/
RUN go mod vendor
RUN go build -v -o main .
RUN mkdir /app && mv ./main /app/main

FROM alpine:3.15.4
LABEL source = git@github.com:kyma-incubator/milv.git

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && apk add bash

COPY --from=builder /app /app

ENTRYPOINT ["/app/main"]