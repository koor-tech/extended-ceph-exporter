# extended-ceph-exporter

A Helm chart for deploying the extended-ceph-exporter to Kubernetes

![Version: 1.2.1](https://img.shields.io/badge/Version-1.2.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v1.0.2](https://img.shields.io/badge/AppVersion-v1.0.2-informational?style=flat-square)

## Get Repo Info

```console
helm repo add extended-ceph-exporter https://koor-tech.github.io/extended-ceph-exporter
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Install Chart

To install the chart with the release name `my-release`:

```console
helm install --namespace <your-cluster-namespace> my-release extended-ceph-exporter/extended-ceph-exporter
```

The command deploys extended-ceph-exporter on the Kubernetes cluster in the default configuration.

_See [configuration](#configuration) below._

_See [helm install](https://helm.sh/docs/helm/helm_install/) for command documentation._

### Development Build
To deploy from a local build from your development environment:

```console
cd charts/extended-ceph-exporter
helm install --namespace <your-cluster-namespace> my-release . -f values.yaml
```

## Uninstall Chart

To uninstall/delete the my-release deployment:

```console
helm delete --namespace <your-cluster-namespace> my-release
```

This removes all the Kubernetes components associated with the chart and deletes the release.

_See [helm uninstall](https://helm.sh/docs/helm/helm_uninstall/) for command documentation._

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| additionalEnv | object | `{}` | Will be put in a Secret and used as env vars |
| config.rgw.accessKey | string | Randomly generated | RGW admin access key |
| config.rgw.host | string | First detected RGW endpoint | The Ceph RGW endpoint as a URL, e.g. "https://your-ceph-rgw-endpoint-here:8443" |
| config.rgw.secretKey | string | Randomly generated | RGW admin secret key |
| image.tag | string | `""` | Overrides the image tag whose default is the chart appVersion. |
| postInstallJob.enabled | bool | `true` | If enabled,  will create a rgw admin user `extended-ceph-exporter` either on Rook/Ceph cluster pre upgrade(when having extended-ceph-exporter as a helm dependency) or on post install of extended-ceph-exporter(needs existing Rook/Ceph cluster). This user will be used for extended ceph metrics. |
| prometheusRule.additionalLabels | object | `{}` | Additional Labels for the PrometheusRule object |
| prometheusRule.enabled | bool | `false` | Specifies whether a prometheus-operator PrometheusRule should be created |
| prometheusRule.rules | list | `[]` | Checkout the file for example alerts |
| resources | object | `{}` | We usually recommend not to specify default resources and to leave this as a conscious choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, uncomment the following lines, adjust them as necessary, and remove the curly braces after 'resources:'. |
| service.port | int | `9138` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| serviceMonitor.additionalLabels | object | `{}` | Additional Labels for the ServiceMonitor object |
| serviceMonitor.enabled | bool | `false` | Specifies whether a prometheus-operator ServiceMonitor should be created |
| serviceMonitor.namespaceSelector | string | `nil` |  |
| serviceMonitor.scrapeInterval | string | `"30s"` |  |

