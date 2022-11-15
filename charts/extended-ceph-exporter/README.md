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
helm install --namespace <your-cluster-namespace> my-release extended-ceph-exporter/extended-ceph-exporter 
```

or update `values.yaml` before running:

```console
helm install --namespace <your-cluster-namespace> my-release extended-ceph-exporter/extended-ceph-exporter -f values.yaml
```

## Configuration

The following tables lists the configurable parameters of the rook-operator chart and their default values.

| Parameter              | Description                                                                       | Default                     |
|------------------------|-----------------------------------------------------------------------------------|-----------------------------|
| `config.rgw.host`      | The Ceph RGW endpoint as a URL, e.g. `"https://your-ceph-rgw-endpoint-here:8443"` | First detected RGW endpoint |
| `config.rgw.accessKey` | The RGW admin access key                                                          | Randomly generated          |
| `config.rgw.secretKey` | The RGW admin secret key                                                          | Randomly generated          |

### Development Build
To deploy from a local build from your development environment:

```console
cd charts/extended-ceph-exporter
helm install --namespace <your-cluster-namespace> my-release . -f values.yaml
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
helm delete --namespace <your-cluster-namespace> my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
