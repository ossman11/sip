minikube start --vm-driver "hyperv" --hyperv-virtual-switch "Minikube" --disk-size 40g --cpus 4 --memory 8192 --v 9999

kubectl get all

@FOR /f "tokens=*" %%i IN ('minikube docker-env') DO @%%i

SETX DOCKER_TLS_VERIFY %DOCKER_TLS_VERIFY%
SETX DOCKER_HOST %DOCKER_HOST%
SETX DOCKER_CERT_PATH %DOCKER_CERT_PATH%
SETX DOCKER_API_VERSION %DOCKER_API_VERSION%

RefreshEnv
