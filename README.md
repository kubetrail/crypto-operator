# crypto-operator
kubernetes operator to track crypto coin prices and balances

## installation
first download the code, build container image and push
to your container registry.
> please make sure go toolchain and docker are installed
> at relatively newer versions and also update the
> IMG value to point to your registry
```bash
export IMG=docker.io/your-account-name/cypto-operator:0.0.1
make generate
make manifests
make docker-build
make docker-push
```
once the container image is available in your registry you can
deploy the controller.
> please make sure you have cert-manager and prometheus running
> on your cluster

install `cert-manager`
```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml
```

install `prometheus` after creating namespace for it and making sure
your `helm` repos are updated
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm --namespace=prometheus-system upgrade --install \
                prometheus prometheus-community/kube-prometheus-stack \
                --set=grafana.enabled=false \
                --version=27.0.1
```

install `CRD's` and the controller
```bash
make install
make deploy
```

make sure all pods and services are up and running
```bash
kubectl --namespace=crypto-operator-system get pods,svc,configmaps,secrets,servicemonitors
NAME                                                      READY   STATUS    RESTARTS   AGE
pod/crypto-operator-controller-manager-5cd5768d5d-222m9   2/2     Running   0          32s

NAME                                                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/crypto-operator-controller-manager-metrics-service   ClusterIP   00.000.000.000   <none>        8443/TCP   32m
service/crypto-operator-webhook-service                      ClusterIP   00.000.000.000   <none>        443/TCP    32m

NAME                                       DATA   AGE
configmap/crypto-operator-manager-config   1      32m
configmap/ee3af832.kubetrail.io            0      31m
configmap/kube-root-ca.crt                 1      32m

NAME                                                    TYPE                                  DATA   AGE
secret/crypto-operator-controller-manager-token-wrdbn   kubernetes.io/service-account-token   3      32m
secret/default-token-k5z4n                              kubernetes.io/service-account-token   3      32m
secret/webhook-server-cert                              kubernetes.io/tls                     3      32m

NAME                                                                                      AGE
servicemonitor.monitoring.coreos.com/crypto-operator-controller-manager-metrics-monitor   32m
```

## track coins
```bash
kubectl create ns my-coins
kubectl --namespace=my-coins create -f config/samples/crypto_v1beta1_coin.yaml
```
```bash
kubectl --namespace=my-coins get coins.crypto.kubetrail.io 
NAME          STATUS    TICKER   PRICE      NUMCOINS   BALANCE          AGE
coin-sample   running   BTC      41815.68   100        4181568.000000   11s
```
the price updates every minute
