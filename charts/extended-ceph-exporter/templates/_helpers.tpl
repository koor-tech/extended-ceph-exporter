{{/*
Expand the name of the chart.
*/}}
{{- define "extended-ceph-exporter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "extended-ceph-exporter.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "extended-ceph-exporter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "extended-ceph-exporter.labels" -}}
helm.sh/chart: {{ include "extended-ceph-exporter.chart" . }}
{{ include "extended-ceph-exporter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "extended-ceph-exporter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "extended-ceph-exporter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "extended-ceph-exporter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "extended-ceph-exporter.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
RGW Host value
*/}}
{{- define "extended-ceph-exporter.rgwHost" -}}
{{- with .Values.config.rgw.host }}
{{- .Values.config.rgw.host }}
{{- else }}
{{- range (lookup "ceph.rook.io/v1" "CephObjectStore" "" "").items }}
{{- .status.info.endpoint}}
{{- break }}
{{- end }}
{{- end }}
{{- end }}
