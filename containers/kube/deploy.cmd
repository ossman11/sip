kubectl create -f ./containers/kube/pod.json
kubectl expose deployment sip --type=NodePort
kubectl get services sip