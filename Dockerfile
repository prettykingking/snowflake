FROM golang:1.17-alpine3.14 AS build

ARG SRC_DIR=/go/src/github.com/prettykingking/snowflake

RUN set -eux \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk --no-cache --no-progress add bash make \
    && mkdir -p $SRC_DIR

WORKDIR $SRC_DIR

# Download go modules
COPY go.mod $SRC_DIR
COPY go.sum $SRC_DIR
RUN GO111MODULE=on GOPROXY=https://goproxy.cn,direct go mod download

COPY . $SRC_DIR

RUN make


FROM alpine:3.14

ARG SRC_DIR=/go/src/github.com/prettykingking/snowflake

ENV INSTALL_DIR=/opt/snowflake \
    RUN_AS_USER=worker \
    TARGET=snowflake

RUN set -eux \
    && addgroup -g 1000 $RUN_AS_USER \
    && adduser -D -H -u 1000 -G $RUN_AS_USER -g $RUN_AS_USER $RUN_AS_USER \
    && mkdir -p $INSTALL_DIR/bin \
    && chown -R $RUN_AS_USER:$RUN_AS_USER $INSTALL_DIR

COPY --from=build --chown=$RUN_AS_USER $SRC_DIR/dist/$TARGET $INSTALL_DIR/bin/$TARGET

# copy config files
COPY --chown=$RUN_AS_USER snowflake.sample.toml $INSTALL_DIR/
COPY --chown=$RUN_AS_USER snowflake.sample.toml $INSTALL_DIR/snowflake.toml

USER $RUN_AS_USER
WORKDIR $INSTALL_DIR
