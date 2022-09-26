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
helm install my-release extended-ceph-exporter/extended-ceph-exporter
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
