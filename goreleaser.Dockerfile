FROM alpine:3.22.1
WORKDIR /data
LABEL org.opencontainers.image.source="https://github.com/gabe565/subcablemap-dl"
COPY subcablemap-dl /usr/bin
ENTRYPOINT ["subcablemap-dl"]
