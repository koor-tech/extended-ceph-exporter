FROM quay.io/prometheus/busybox:latest

ARG BUILD_DATE="N/A"
ARG REVISION="N/A"

LABEL org.opencontainers.image.authors="Alexander Trost <alexander@koor-tech>" \
    org.opencontainers.image.created="${BUILD_DATE}" \
    org.opencontainers.image.title="koor-tech/extended-ceph-exporter" \
    org.opencontainers.image.description="A Prometheus exporter to provide \"extended\" metrics about a Ceph cluster's running components (e.g., RGW)." \
    org.opencontainers.image.documentation="https://github.com/koor-tech/extended-ceph-exporter/blob/main/README.md" \
    org.opencontainers.image.url="https://github.com/koor-tech/extended-ceph-exporter" \
    org.opencontainers.image.source="https://github.com/koor-tech/extended-ceph-exporter" \
    org.opencontainers.image.revision="${REVISION}" \
    org.opencontainers.image.vendor="koor-tech" \
    org.opencontainers.image.version="N/A"

ADD .build/linux-amd64/extended-ceph-exporter /bin/extended-ceph-exporter

ENTRYPOINT ["/bin/extended-ceph-exporter"]
