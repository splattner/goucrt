FROM golang:1.20

WORKDIR /go/src/app
COPY ./ .

RUN ls -al

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT [ "/go/bin/ucrt" ]
CMD ["server"]