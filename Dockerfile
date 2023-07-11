FROM golang:1.20.5-bullseye AS builder

ADD ./ /etc/encoder-cli

RUN cd /etc/encoder-cli && \
 go env -w GOPROXY=https://goproxy.cn,direct && \
 go env -w GOPRIVATE=github.com/ah-its-andy && \
 go mod download 

RUN cd /etc/encoder-cli && go build ./

FROM ubuntu AS base-image

RUN apt-get update && apt-get install ffmpeg mkvtoolnix -y && apt-get clean 

FROM base-image

COPY --from=0 /etc/encoder-cli/encoder-cli /etc/encoder-cli/encoder-cli

ADD conf /etc/encoder-cli/conf

CMD ["/etc/encoder-cli/encoder-cli", "-c", "/etc/encoder-cli/conf", "-t", "/etc/encoder-cli/task/default.yaml"]