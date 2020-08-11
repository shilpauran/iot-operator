The Operator itself should be deployed via [Helm chart](../../helm/iot-operator) in the cluster.

For development of the cluster however the developer needs an own namespace for the Operator. This can be setup using [namespace.yml](namespace.yml). Make sure the placeholder is resolved there before applying with:
```sh
$ kubectl apply -f deploy/namespace.yml
```

In case the developer works on a fresh cluster with no productive Operator installed, it may be necessary to register the custom resource definition (CRD) of `IoTService`:

```sh
$ kubectl apply -f deploy/crds/cp_v1alpha1_iotservice_crd.yml
```