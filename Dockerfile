## golang image created 2021-06-24T00:31:06.02014601Z 
FROM docker.io/library/golang@sha256:be99fa59acd78bb22a41bbc1e15ebfab2262498ee0c2e28c3d09bc44d51d1774 AS builder
WORKDIR /go/src/preflight
COPY . .
RUN make build

# ubi-minimal image created 2021-06-01T12:12:46.922866Z
FROM registry.access.redhat.com/ubi8/ubi-minimal@sha256:0ccb9988abbc72d383258d58a7f519a10b637d472f28fbca6eb5fab79ba82a6b

# Add preflight binary
COPY --from=builder /go/src/preflight/preflight /usr/local/bin/preflight

# Install dependencies
RUN microdnf install \
    buildah \
    bzip2 \
    gzip \
    iptables \
    podman \
    skopeo

# Install OpenShift client binary
RUN curl -L https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-client-linux-4.7.18.tar.gz | tar -xzv -C /usr/local/bin --exclude=README.md --exclude=kubectl

# Install Operator SDK binray
RUN curl -Lo /usr/local/bin/operator-sdk https://github.com/operator-framework/operator-sdk/releases/download/v1.9.0/operator-sdk_linux_amd64
RUN chmod 755 /usr/local/bin/operator-sdk

ENTRYPOINT ["/usr/local/bin/preflight"]
CMD ["--help"]
