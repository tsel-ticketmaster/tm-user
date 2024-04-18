# Image Builder
FROM golang:1.22-alpine3.19 AS go-builder

LABEL maintainer="patrick_m_sangian@telkomsel.co.id"

# Set Working Directory
WORKDIR /usr/src/app

# Copy Source Code
COPY . ./

# Dependencies installation and binary file builder
RUN apk add --no-cache make

RUN make install \
  && make build


# Final Image
# ---------------------------------------------------
FROM alpine:3.19

RUN apk add --no-cache curl tzdata

RUN /bin/ls /usr/share/zoneinfo
RUN /bin/cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
RUN echo "Asia/Jakarta" >  /etc/timezone

# Set Working Directory
WORKDIR /usr/src/app

# Copy Anything The Application Needs
COPY --from=go-builder /tmp/app ./

# Expose Application Port
EXPOSE 9000

# Run The Application
CMD ["./app"]