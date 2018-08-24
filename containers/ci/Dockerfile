FROM travisci/ci-garnet:packer-1515445631-7dfb2e1

ENV CHANGE_MINIKUBE_NONE_USER=true

USER travis

RUN curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.7.0/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
RUN curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
RUN sudo minikube start --vm-driver=none
RUN minikube update-context
RUN JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1; done
