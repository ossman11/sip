FROM travisci/ci-garnet:packer-1515445631-7dfb2e1

ENV KUBERNETES_VERSION=v1.10.0
ENV CHANGE_MINIKUBE_NONE_USER=true

COPY . .

RUN ./containers/ci/minikube

USER travis
