{{/*
Expand the name of the chart.
*/}}
{{- define "jlpt5.name" -}}
{{- default .Chart.Name .Values.global.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "jlpt5.fullname" -}}
{{- if .Values.global.fullnameOverride }}
{{- .Values.global.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.global.nameOverride }}
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
{{- define "jlpt5.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "jlpt5.labels" -}}
helm.sh/chart: {{ include "jlpt5.chart" . }}
{{ include "jlpt5.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "jlpt5.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jlpt5.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
PostgreSQL component labels
*/}}
{{- define "jlpt5.postgresql.labels" -}}
helm.sh/chart: {{ include "jlpt5.chart" . }}
{{ include "jlpt5.postgresql.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: database
{{- end }}

{{- define "jlpt5.postgresql.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jlpt5.name" . }}-postgresql
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: database
{{- end }}

{{/*
Backend component labels
*/}}
{{- define "jlpt5.backend.labels" -}}
helm.sh/chart: {{ include "jlpt5.chart" . }}
{{ include "jlpt5.backend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: backend
{{- end }}

{{- define "jlpt5.backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jlpt5.name" . }}-backend
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Frontend component labels
*/}}
{{- define "jlpt5.frontend.labels" -}}
helm.sh/chart: {{ include "jlpt5.chart" . }}
{{ include "jlpt5.frontend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: frontend
{{- end }}

{{- define "jlpt5.frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jlpt5.name" . }}-frontend
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
PostgreSQL service name
*/}}
{{- define "jlpt5.postgresql.serviceName" -}}
{{- printf "%s-postgresql" (include "jlpt5.fullname" .) }}
{{- end }}

{{/*
Backend service name
*/}}
{{- define "jlpt5.backend.serviceName" -}}
{{- printf "%s-backend" (include "jlpt5.fullname" .) }}
{{- end }}

{{/*
Frontend service name
*/}}
{{- define "jlpt5.frontend.serviceName" -}}
{{- printf "%s-frontend" (include "jlpt5.fullname" .) }}
{{- end }}
