FROM golang:1.18.2-bullseye

LABEL maintainer="Chihiro Hasegawa <encry1024@gmail.com>"

RUN go install -v golang.org/x/tools/gopls@latest && \
    go install -v github.com/ramya-rao-a/go-outline@latest &&\
    go install -v golang.org/x/tools/cmd/goimports@latest
ENTRYPOINT ["/bin/bash"]
