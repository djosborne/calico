FROM quay.io/operator-framework/helm-operator:v0.4.0

COPY charts/ ${HOME}/helm-charts/
COPY manifests/ ${HOME}/manifests/
COPY helm-operator/watches.yaml ${HOME}/watches.yaml
