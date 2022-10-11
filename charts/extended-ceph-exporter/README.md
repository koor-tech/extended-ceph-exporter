# extended-ceph-exporter

* Installs the extended-ceph-exporter

## Get Repo Info

```console
helm repo add extended-ceph-exporter https://koor-tech.github.io/extended-ceph-exporter
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release extended-ceph-exporter/extended-ceph-exporter \
--set config.rgw.accessKey=$(< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c10) \
--set config.rgw.secretKey=$(< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c10) \
--set config.rgw.host=https://your-ceph-rgw-endpoint-here
-namespace <your-cluster-namespace>
```

or update `values.yaml` before running:

```console
helm install my-release extended-ceph-exporter/extended-ceph-exporter -f values.yaml
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
