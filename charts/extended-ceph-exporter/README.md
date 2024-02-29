# extended-ceph-exporter

A Helm chart for deploying the extended-ceph-exporter to Kubernetes

![Version: 1.3.5](https://img.shields.io/badge/Version-1.3.5-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v1.2.0](https://img.shields.io/badge/AppVersion-v1.2.0-informational?style=flat-square)

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

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| additionalEnv | object | `{}` | Will be put in a Secret and used as env vars |
| affinity | object | `{}` | [Affinity](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity) |
| autoscaling | object | `{"enabled":false,"maxReplicas":100,"minReplicas":1,"targetCPUUtilizationPercentage":80}` | Autoscaling configuration |
| config.rgw.accessKey | string | Randomly generated | RGW admin access key |
| config.rgw.existingSecret | object | `{"keys":{"access":"access","secret":"secret"},"name":""}` | Existing RGW admin credentials secret config |
| config.rgw.existingSecret.keys.access | string | `"access"` | Access key secret key name |
| config.rgw.existingSecret.keys.secret | string | `"secret"` | Secret key secret key name |
| config.rgw.existingSecret.name | string | `""` | Name of the existing RGW admin credentials secret |
| config.rgw.host | string | First detected RGW endpoint | The Ceph RGW endpoint as a URL, e.g. "https://your-ceph-rgw-endpoint-here:8443" |
| config.rgw.secretKey | string | Randomly generated | RGW admin secret key |
| fullnameOverride | string | `""` | Override fully-qualified app name |
| image.tag | string | `""` | Overrides the image tag whose default is the chart appVersion. |
| nameOverride | string | `""` | Override chart name |
| nodeSelector | object | `{}` | [Create a pod that gets scheduled to your chosen node](https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/#create-a-pod-that-gets-scheduled-to-your-chosen-node) |
| podAnnotations | object | `{}` | Annotations to add to the pod |
| podSecurityContext | object | `{}` | [Pod security context](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) |
| postInstallJob.enabled | bool | `true` | If enabled,  will create a rgw admin user `extended-ceph-exporter` either on Rook/Ceph cluster pre upgrade(when having extended-ceph-exporter as a helm dependency) or on post install of extended-ceph-exporter(needs existing Rook/Ceph cluster). This user will be used for extended ceph metrics. |
| prometheusRule.additionalLabels | object | `{}` | Additional Labels for the PrometheusRule object |
| prometheusRule.enabled | bool | `false` | Specifies whether a prometheus-operator PrometheusRule should be created |
| prometheusRule.rules | prometheusrules.monitoring.coreos.com | `[]` |  |
| replicaCount | int | `1` | Number of replicas of the exporter |
| resources | object | `{"limits":{"cpu":"125m","memory":"150Mi"},"requests":{"cpu":"25m","memory":"150Mi"}}` | These are sane defaults for "small" object storages |
| securityContext | object | `{}` | [Security context](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) |
| service.port | int | `9138` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| serviceMonitor.additionalLabels | object | `{}` | Additional Labels for the ServiceMonitor object |
| serviceMonitor.enabled | bool | `false` | Specifies whether a prometheus-operator ServiceMonitor should be created |
| serviceMonitor.namespaceSelector | string | `nil` |  |
| serviceMonitor.scrapeInterval | duration | `"30s"` | Interval at which metrics should be scraped |
| serviceMonitor.scrapeTimeout | duration | `"20s"` | Timeout for scraping |
| tolerations | list | `[]` | [Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) |
