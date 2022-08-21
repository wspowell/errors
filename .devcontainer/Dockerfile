FROM golang:1.19-alpine

RUN apk add --no-cache build-base git bash curl the_silver_searcher

# Install golang dev tools
ENV PATH="/usr/local/go/bin:/app/bin:${PATH}"
ENV GOBIN=/usr/local/go/bin
RUN mkdir -p /tmp/gotools
WORKDIR /tmp/gotools
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.47.3
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/go/bin v1.47.3
RUN go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@v2.1.2
# The version is correct for go-outline. The tag is not prefixed with 'v'.
RUN go install github.com/ramya-rao-a/go-outline@1.0.0
RUN go install github.com/cweill/gotests/...@v1.6.0
RUN go install github.com/fatih/gomodifytags@v1.16.0
RUN go install github.com/josharian/impl@v1.1.0
RUN go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.9.0
RUN go install golang.org/x/tools/gopls@v0.9.1
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0
RUN go install github.com/vektra/mockery/v2@v2.14.0
RUN go install gotest.tools/gotestsum@v1.8.1
RUN go install mvdan.cc/gofumpt@v0.3.1
RUN go install github.com/daixiang0/gci@v0.6.2
# Clean up
RUN rm -rf /tmp/gotools /root/go

RUN git config --global --add safe.directory /workspace/errors