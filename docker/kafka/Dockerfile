FROM openjdk:8-jre-alpine

ENV KAFKA_USER=kafka \
    KAFKA_CONF_DIR=/conf \
    KAFKA_DATA_DIR=/data \
    KAFKA_LOG_DIR=/logs \
    KAFKA_PORT=9092 \
    KAFKA_CREATE_TOPICS=choraldatastream:1:1 \
    SCALA_VERSION=2.12 \
    KAFKA_VERSION=0.10.2.1

# Install required packages
RUN apk add --no-cache \
    bash \
    python \
    su-exec \
    net-tools

# Add a user and make dirs
RUN set -x \
    && adduser -D "$KAFKA_USER" \
    && mkdir -p "$KAFKA_CONF_DIR" "$KAFKA_DATA_DIR" "$KAFKA_LOG_DIR" \
    && chown -R "$KAFKA_USER:$KAFKA_USER" "$KAFKA_CONF_DIR" "$KAFKA_DATA_DIR" "$KAFKA_LOG_DIR"

ARG DISTRO_NAME=kafka_"$SCALA_VERSION"-"$KAFKA_VERSION"

# Download Apache Storm, untar and clean up
RUN set -x \
    && apk add --no-cache --virtual .build-deps \
    && wget -q "http://www.apache.org/dist/kafka/$KAFKA_VERSION/$DISTRO_NAME.tgz" \
    && wget -q "http://www.apache.org/dist/kafka/$KAFKA_VERSION/$DISTRO_NAME.tgz.asc" \
    && tar -xzf "$DISTRO_NAME.tgz" \
    && chown -R "$KAFKA_USER:$KAFKA_USER" "$DISTRO_NAME" \
    && rm -r "$DISTRO_NAME.tgz" "$DISTRO_NAME.tgz.asc" \
    && apk del .build-deps

WORKDIR $DISTRO_NAME

ENV PATH $PATH:/$DISTRO_NAME/bin
ENV KAFKA_HOME=/$DISTRO_NAME

COPY docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]