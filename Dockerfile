FROM scratch as builder
COPY check /
COPY in /
COPY out /

FROM alpine:edge AS resource
RUN mkdir /opt/resource
COPY --from=builder / /opt/resource
RUN chmod +x /opt/resource/*

# Test binaries exist
RUN stat /opt/resource/check /opt/resource/in /opt/resource/out
