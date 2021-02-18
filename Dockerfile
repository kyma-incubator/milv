FROM golang:1.15-alpine3.12 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/milv
WORKDIR ${BASE_APP_DIR}

COPY ./ ${BASE_APP_DIR}/

RUN go build -v -o main .
RUN mkdir /app && mv ./main /app/main

FROM alpine:3.12
LABEL source = git@github.com:kyma-incubator/milv.git

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && apk add bash

COPY --from=builder /app /app

ENTRYPOINT ["/app/main"]