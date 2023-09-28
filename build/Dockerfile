FROM golang:1.21 AS build-env

WORKDIR /go/src/app
COPY ./ .

RUN ls -al

RUN go get -d -v ./...
RUN go install -v ./...


FROM alpine:3

LABEL maintainer "Sebastian Plattner"
LABEL org.opencontainers.image.title "Unfolded Cirrle Remote Two Integraton Driber"
LABEL org.opencontainers.image.description "Golang based integration driver for a Unfolded Circle Remote Two"
LABEL org.opencontainers.image.authors "Sebastian Plattner"
LABEL org.opencontainers.image.source https://github.com/splattner/goucrt
LABEL org.opencontainers.image.licenses MPL-2.0

RUN apk --no-cache add libc6-compat

COPY --from=build-env /go/bin/ucrt /app/

WORKDIR /app

RUN adduser -D nonroot && mkdir /app/ucconfig && chown nonroot /app/ucconfig
USER nonroot

ENTRYPOINT [ "/app/ucrt" ]
CMD [ "help" ]
