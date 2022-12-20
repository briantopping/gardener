{{- define "extraMounts.gardener.controlPlane" -}}
{{- if .Values.gardener.controlPlane.deployed }}
- hostPath: {{.Values.gardener.templateRoot}}/example/gardener-local/controlplane
  containerPath: /etc/gardener/controlplane
  readOnly: true
{{- end }}
{{- end -}}

{{- define "extraMounts.gardener.seed.backupBuckets" -}}
{{- if .Values.gardener.seed.deployed -}}
- hostPath: {{.Values.gardener.templateRoot}}/dev/local-backupbuckets
  containerPath: /etc/gardener/local-backupbuckets
{{- end -}}
{{- end -}}

{{- define "extraMounts.registry" -}}
{{- if .Values.registry.deployed }}
- hostPath: {{.Values.gardener.templateRoot}}/dev/local-registry
  containerPath: /etc/gardener/local-registry
{{- end }}
{{- end -}}
