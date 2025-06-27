{{/* k8s-controller Helm helpers */}}

{{- define "k8s-controller.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "k8s-controller.name" -}}
{{- .Chart.Name -}}
{{- end }}
