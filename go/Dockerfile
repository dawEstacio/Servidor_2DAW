FROM golang:1.15

ARG MODULE_NAME

WORKDIR /go/src/goApp
RUN rm -f go.mod || true && \
    rm -f go.sum || true && \
    go mod init ${MODULE_NAME} && \
    go mod tidy && \
    go get github.com/pilu/fresh 

EXPOSE 3000
CMD ["fresh"]