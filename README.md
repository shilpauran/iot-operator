# IoT Service Operator

The Operator is designed to manage one or more IoT service instances. This Operator is based on the [operator-sdk](https://github.com/operator-framework/operator-sdk). To learn more about the writing an Operator in Go, see the [user guide](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md).

The Operator watches objects of custom resource definition [IoTService](deploy/crds/cp_v1alpha1_iotservice_crd.yml) inside the namespace where the Operator is deployed in. The name of the Operator (the namespace the Operator runs within) represents the "IoT landscape", here `devawsk8s`.

In order to watch the resources in these namespaces the Operator has to run as [cluster-scoped](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#operator-scope) Operator (setting the `WATCH_NAMESPACE` environment variable to empty string). To still focus the scope of the Operator to a namespace, an additional environment variable `OPERATOR_NAME` has to be set to match the namespace the Operator should watch.  For the typical development setup you would have the following:
```
export WATCH_NAMESPACE=
export OPERATOR_NAME="devawsk8s"
```

![namespace overview](../docs/Namespace-Overview.png "Namespace Overview")


Inside the namespace of the Operator, there have to exist some resources:
- Secret: image-pull-secret - used for the Operator pod itself, additionally is copied by the Operator to each IoT service instance namespace
- ServiceAccount - the context, in which the Operator process accesses the Kubernetes cluster

Furthermore for [RBAC] there have to exist some cluster-scoped resources:
- ClusterRole - the role allowing all operations the Operator needs to perform
- ClusterRoleBinding - binds the ClusterRole to the ServiceAccount

For each IoT service instance the Operator creates a new namespace named `<operatorname>-<instanceid>`, where `operatorname` is the name of the Operator and `instanceid` is the instance id of the IoT service instance, specified as [name](deploy/crds/example-iotservice.yml#L4) of the `IoTService` object. Within that namespace the Operator creates all necessary resources. The following activity diagram depicts this process.

![deployment via operator](../docs/K8s-Deployment-via-Operator.png "Deployment via Operator")

## Progress Monitoring
The Operator implements monitoring logic to check the state of all resources in the context of an IoT service instance. Out of these states the Operator computes the state for the IoT service instance and stores it within the `status` subresource of the `IoTService` object.

The Operator manages the fields:
- `.status.phase` - one of `Available`, `Processing`, `Failed`
- `.status.observedGeneration` - always set to `metadata.generation` as soon as the Operator touches the `IoTService` object
- `.status.reason` - only populated in case `status.phase` is `Failed`, the reason for the failure
- `.status.pods` - an array of the pod names currently running

## Development environment

### Prerequisites
- [go](https://golang.org/doc/install) version v1.10+.
- [dep](https://golang.github.io/dep/docs/installation.html) version v0.5.0+.
- git
- [docker](https://docs.docker.com/install/) version 17.03+.
  - Assuming you are using the Linux Subsystem for Windows (WSL) then you can run Docker Desktop for Windows
  - This installs a Kubernetes enabled version of Docker
  - This also means you must configure WSL to use your Windows version of Docker
  - [This](https://nickjanetakis.com/blog/setting-up-docker-for-windows-and-wsl-to-work-flawlessly) guide does a good job of outlining the steps you need to take
    - Docker for Windows should "Expose daemon on tcp://localost:2375 without TLS"
    - After installing Python and Pip and docker-compose
    - `echo "export DOCKER_HOST=tcp://localhost:2375" >> ~/.bashrc && source ~/.bashrc`
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.11.0+.
- Access to a Kubernetes v.1.11.0+ cluster.
- IDE, e.g. [GoLand](https://www.jetbrains.com/go/) 

### Install the Operator SDK CLI
The Operator SDK has a CLI tool that helps the developer to create, build, and deploy an Operator project.

Checkout the desired release tag and install the SDK CLI tool:

```sh
$ mkdir -p $GOPATH/src/github.com/operator-framework
$ cd $GOPATH/src/github.com/operator-framework
$ git clone https://github.com/operator-framework/operator-sdk
$ cd operator-sdk
$ git checkout <RELEASE TAG>
$ make dep
$ make install
```
Note that we work currently with release tag v0.4.0 for the operator-sdk.

### Hint on repository organization
As Go projects should all be located on `$GOPATH/src/...` and the use of symbolic links are explicitly discouraged, the developer will most likely end up in having two `iotservices-k8s` Git repositories, one in the default `~/git/` path, one in `$GOPATH/src/github.wdf.sap.corp/iotservices/iotservices-k8s`.

In order to avoid the confusions which result from that, it is advisable to only use the one repo in `$GOPATH/src/github.wdf.sap.corp/iotservices/iotservices-k8s`. If necessary, the developer could create a symlink in `~/git/`, which references the real repository.
```sh
$ mkdir -p $GOPATH/src/github.wdf.sap.corp/iotservices
$ cd $GOPATH/src/github.wdf.sap.corp/iotservices
$ git clone https://github.wdf.sap.corp/iotservices/iotservices-k8s
$ cd iotservices-k8s/iot-operator
```
For developers working with Windows, it is even a bit more complicated. As not all tools in the K8s context work under Windows, it may be necessary to work with the [Windows Subsystem for Linux](https://docs.microsoft.com/de-de/windows/wsl/install-win10) (WSL). In order to access the same repository also from WSL, it is advisable to create only a `GOPATH` in the Windows file system, e.g. in `%USERPROFILE%/go/` and create in the Linux file system a symbolic link `~/go`, pointing to the `GOPATH` in the Windows file system.
  ```
  # Assume you ahve installed Go in your Windows %USERPROFILE%\go
  # In the Linux Subsystem
  # ln -s /mnt/c/Users/I#/go ~/go
  # ~/.bashcrc
  export GOPATH=~/go
  export GOBIN=$GOPATH/bin
  export PATH=$PATH:$GOBIN
  # We do not want a WATCH_NAMESPACE defined for cluster-scoped Operators
  export WATCH_NAMESPACE=
  # There should be one operator per landscape
  export OPERATOR_NAME="devawsk8s"
  ```
    
### Install dependencies
After initial repository checkout and regularly upon changes in the project imports run:
```sh
$ cd $GOPATH/src/github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator
$ dep ensure -v
```

### Run Operator locally
The activity diagram above shows how to separate productive landscape Operator and locally deployed Operators in development mode: Each Operator connected with the Kubernetes cluster gets informed about changes to resources, however only resources with matching namespace are processed.
For local development of the Operator, therefore configure to watch for a different namespace than the productive landscape Operator to avoid conflicts.

Set the following environment variables:
* `OPERATOR_NAME` - name of the operator/operator namespace
* `WATCH_NAMESPACE` - has to be empty string to be able to watch for secondary resources

Further make sure that a namespace with same name as `OPERATOR_NAME` exist, and that the Secret `image-pull-secret` exists in there as described above (check [deploy/namespace.yml](deploy/namespace.yml) how to create such namespace). For local runs of the Operator the ServiceAccount, ClusterRole and ClusterRoleBinding do not have to exist necessarily.

Before running the Operator, the custom resource definition (CRD) of `IoTService` must be registered with the Kubernetes apiserver. However, this may only necessary on a fresh cluster, where no productive operator has been deployed yet. To do so:

```sh
$ kubectl apply -f deploy/crds/cp_v1alpha1_iotservice_crd.yml
```

#### Run Operator locally from IDE (GoLand)
To configure a launch configuration for the Operator, use a "Go Build" of "Run kind: File" and choose [cmd/manager/main.go](cmd/manager/main.go) under "Files".
Set environment variables `WATCH_NAMESPACE` and `OPERATOR_NAME` as described above. 

#### Run Operator locally with Operator-SDK
Set environment variable `OPERATOR_NAME` as described above using `export`. `WATCH_NAMESPACE` cannot be set directly as environment variable. Instead pass empty string to `namespace` parameter of the start command:

```sh
$ export OPERATOR_NAME=<dev-operator-name>
$ operator-sdk up local --namespace=
```

This starts the Operator locally with the default kubernetes config file present at `$HOME/.kube/config`. You can use a specific kubeconfig via the flag `--kubeconfig=<path/to/kubeconfig>`.

## Deploy and run the Operator in the cluster

Eventually the Operator should be deployed in the cluster. For this there are two ways. Both of them require first to build a docker image.

### Build Docker image
To have a sensible version, we will generate the version out of the last git commit, which will look like this:

```
format:  yyyyMMdd-HHmmss.abbreviatedCommitHash
example: 20160702-152019.75c54f5
```

Build the iot-operator image using the Operator-SDK CLI:

```
$ cd src/github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator
$ operator-sdk build iot-service-dev.docker.repositories.sap.ondemand.com/iot-operator:$(git log -1 --format=%cd.%h --date=format:%Y%m%d-%H%M%S)
```

Push it to a Docker registry:
NOTE: You must first login / prepare for the docker registry.  As a minimum perform the Docker login from [here](../components/iot/core-spring/README.md).

```
$ cd src/github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator
$ docker push iot-service-dev.docker.repositories.sap.ondemand.com/iot-operator:$(git log -1 --format=%cd.%h --date=format:%Y%m%d-%H%M%S)
```

Note down the version which was generated in the former steps.

### Deploying the Operator with `kubectl`
Before running the Operator, the custom resource definition (CRD) of `IoTService` must be registered with the Kubernetes apiserver, this holds true also for running the Operator locally:

```sh
$ kubectl apply -f deploy/crds/cp_v1alpha1_iotservice_crd.yml
```

The Deployment manifest is generated at `deploy/operator.yml`. Be sure to update the deployment image since the default is just a placeholder, e.g.

`
$ sed -i 's|REPLACE_VERSION|20160702-152019.75c54f5|g' deploy/operator.yml
`

or under OSX

`
$ sed -i "" 's|REPLACE_VERSION|20160702-152019.75c54f5|g' deploy/operator.yml
`

The Operator has to be deployed in the namespace where `IoTService` objects have to be created in. Make sure you specify that namespace correctly in each of the manifests used in the following steps. Fill in the image-pull-secret in `deploy/namespace.yml` as Base64 encoded string.
Setup RBAC and deploy the iot-operator:

```sh
$ kubectl apply -f deploy/namespace.yml
$ kubectl apply -f deploy/service_account.yml
$ kubectl apply -f deploy/role.yml
$ kubectl apply -f deploy/role_binding.yml
$ kubectl apply -f deploy/operator.yml
```

Verify that the iot-operator is up and running:

```sh
$ kubectl --namespace <operator-namespace> get deployments
NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
iot-operator       1         1         1            1           1m
```
### Deploying the Operator with the Helm chart

See [here](../helm/iot-operator/README.md).

For deploying the Operator in the cluster with the according Helm chart, see [here](../helm/iot-operator/README.md).


## Create an IoT service instance via the Operator
Create the example `IoTService` custom object that was generated at deploy/crds/example-iotservice.yml

```sh
$ kubectl apply -f deploy/crds/example-iotservice.yml

$ kubectl --namespace <dev-operator-name> get iotservice
NAME                  CREATED AT
example-iotservice   7s
```

## Run the tests
To run the tests locally, run:

```
operator-sdk test local ./test/e2e --namespaced-manifest deploy
```

TODO: Currently there is only one test...

## Project Scaffolding Layout
| File/Folders       | Purpose                           |
| :---               | :--- |
| cmd                | Contains `manager/main.go` which is the main program of the Operator. This instantiates a new manager which registers all custom resource definitions under `pkg/apis/...` and starts all controllers under `pkg/controllers/...`  . |
| pkg/apis           | Contains the directory tree that defines the APIs of the Custom Resource Definitions(CRD). The API for the primary resource type `IoTService` is defined in `pkg/apis/cp/v1alpha1/iotservices_types.go`. |
| pkg/configmap      | Contains the [ConfigMap][ConfigMap] which is mounted as at `/etc/sap/config.yml` in the pod's. |
| pkg/controller     | This pkg contains the controller implementations. The reconcile logic for handling resource type  `IoTService` is implemented in `pkg/controller/iotservice/iotservice_controller.go`. |
| pkg/deployment     | Contains the [Deployment][Deployment] objects |
| pkg/namespace      | Contains the [Namespace][Namespace] objects |
| pkg/network_policy | Contains the [NetworkPolicy][NetworkPolicy] objects. Pods are only accept traffic from the same iotservice instance |
| pkg/secret         | Contains the [Secret][Secret] objects |
| pkg/service        | Contains the [Service][Service] objects |
| pkg/stateful_set   | Contains the [StatefulSet][StatefulSet] objects |
| build | Contains the `Dockerfile` and build scripts used to build the Operator. |
| deploy | Contains various YAML manifests for registering CRDs, setting up [RBAC][RBAC], and deploying the Operator as a Deployment.
| Gopkg.toml Gopkg.lock | The [Go Dep][dep] manifests that describe the external dependencies of this Operator. |
| vendor | The golang [vendor][Vendor] folder that contains the local copies of the external dependencies that satisfy the imports of this project. [Go Dep][dep] manages the vendor directly. |

[ConfigMap]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/
[Deployment]: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
[Namespace]: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
[NetworkPolicy]: https://kubernetes.io/docs/concepts/services-networking/network-policies/
[Secret]: https://kubernetes.io/docs/concepts/configuration/secret/
[Service]: https://kubernetes.io/docs/concepts/services-networking/service/
[StatefulSet]: https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/
[RBAC]: https://kubernetes.io/docs/reference/access-authn-authz/rbac/
[Vendor]: https://golang.org/cmd/go/#hdr-Vendor_Directories
[dep]: https://github.com/golang/dep
