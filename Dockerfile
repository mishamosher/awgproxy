# Start by building the application.
FROM docker.io/golang:1.22.6 as build

WORKDIR /usr/src/awgproxy
COPY . .

RUN make

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build /usr/src/awgproxy/awgproxy /usr/bin/awgproxy

VOLUME [ "/etc/awgproxy"]
ENTRYPOINT [ "/usr/bin/awgproxy" ]
CMD [ "--config", "/etc/awgproxy/config" ]

LABEL org.opencontainers.image.title="awgproxy"
LABEL org.opencontainers.image.description="AmneziaWG client that exposes itself as a socks5 proxy"
LABEL org.opencontainers.image.licenses="ISC"
