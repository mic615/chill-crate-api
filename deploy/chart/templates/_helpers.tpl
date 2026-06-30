{{- define "chill-crate-api.labels" -}}
app: {{ .Chart.Name }}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "chill-crate-api.selectorLabels" -}}
app: {{ .Chart.Name }}
{{- end -}}
