# extended-ceph-exporter

A Prometheus exporter to provide "extended" metrics about a Ceph cluster's running components (e.g., RGW).

[![Ceph - RGW Bucket Usage Overview Grafana Dashboard Screenshot](grafana/ceph-rgw-bucket-usage-overview.png)](grafana/)

## Requirements

* Needs a Ceph cluster up and running.

* Needs an admin user

    ```
    radosgw-admin user create --uid extended-ceph-exporter --display-name "extended-ceph-exporter admin user" --caps "buckets=read;users=read;usage=read;metadata=read;zone=read"
    # Access key / "Username"
    radosgw-admin user info --uid extended-ceph-exporter | jq '.keys[0].access_key'
    # Secret key / "Password
    radosgw-admin user info --uid extended-ceph-exporter | jq '.keys[0].secret_key'
    ```

## Quickstart

* Clone the repository:
  ```console
  git clone https://github.com/koor-tech/extended-ceph-exporter
  cd extended-ceph-exporter
  ```

* Create a copy of the `.env.example` file and name it `.env`. Configure your RGW credentials and endpoint in the `.env` file.

* Configure Prometheus to collect metrics from the exporter from `:9138/metrics` endpoint using a static configuration, here's a sample scrape job from the `prometheus.yml`:

  ```yaml
  # For more information on Prometheus scrape_configs:
  # https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config
  scrape_configs:

    - job_name: "extended-ceph-metrics"

      # Override the global default and scrape targets from this job every 30 seconds.
      scrape_interval: 30s

      static_configs:
        # Please change the ip address `127.0.0.1` to target the exporter is running
        - targets: ['127.0.0.1:9138']
  ```

* To run the exporter locally, run `go run .`

* Should you have Grafana running for metrics visulization, check out the [Grafana dashboards](grafana/).

### Helm

To install the exporter to Kubernetes using Helm, check out the [extended-ceph-exporter Helm Chart](charts/extended-ceph-exporter/).

## Collectors

There is varying support for collectors. The tables
below list all existing collectors and the required Ceph components.

### Enabled by default

| Name             |                            Description                            | Ceph Component |
| :--------------- | :---------------------------------------------------------------: | -------------- |
| `rgw_buckets`    | Exposes RGW Bucket Usage and Quota metrics from the Ceph cluster. | RGW            |
| `rgw_user_quota` |       Exposes RGW User Quota metrics from the Ceph cluster.       | RGW            |

## Development

### Requirements

* Golang 1.19
* Depending on the module requirements, a Ceph cluster with the respective Ceph components.

### Debugging

A VSCode debug config is available to run and debug the project.

To make the exporter talk with a Ceph RGW S3 endpoint, create a copy of the `.env.example` file and name it `.env`.
Be sure ot add your Ceph RGW S3 endpoint and credentials in it.

### Create a new Release

To create a new release (example is for release `v0.1.2`):

1. Increase the version according to Semantic Versioning in the [`VERSION` file](VERSION).
2. Add a new entry to the [`CHANGELOG.md`](CHANGELOG.md) with the changes and improvements listed in it.
3. Set the new version, which will be the new container image tag, in [the `values.yaml` of the Helm chart here](https://github.com/koor-tech/extended-ceph-exporter/blob/main/charts/extended-ceph-exporter/Chart.yaml#L24) (`appVersion:` field).
4. Check out a new branch, which will be used for the pull request to update the version: `git checkout -b BRANCH_NAME`
5. Commit these changes now using `git commit -s -S`.
6. Push the branch using `git push -u origin BRANCH_NAME` with these changes and create a pull request on [GitHub](https://github.com/koor-tech/extended-ceph-exporter).
7. Wait for pull request to be approved and merge it (if you have access to do so).
8. Create the new tag using `git tag v0.1.2` and then run `git push -u origin v0.1.2`
9. In a few minutes, the CI should have built and published a draft of the release here [GitHub - Releases List](https://github.com/koor-tech/extended-ceph-exporter/releases).
10. Now edit the release and use the green button to publish the release.
11. Congratulations! The release is now fully published.
