## 1.0.3 / 2023-10-18

* [CHORE] Update ceph-go library to 0.24.0
* [FEATURE] helm: add option to use an existing secret for rgw credentials
* [CHORE] Use [helm-docs](https://github.com/norwoodj/helm-docs) to create chart documentation
* [FEATURE] Autodetect the RGW host and autogenerate the RGW access key and secret

## 1.0.2 / 2022-11-14

* [FEATURE] use the dotenv extension to read RGW credentials and endpoint from `.env` file
* [BUGFIX] Increment helm chart version to address documentation changes

## 1.0.1 / 2022-11-14

* [BUGFIX] fix the required flags check to check for the new flag names

## 1.0.0 / 2022-09-26

* [FEATURE] initial release of RGW bucket and user quota metrics module
* [FEATURE] add basic helm chart for deploying the exporter to Kubernetes
