The installation configuration uses [kpt](https://github.com/GoogleContainerTools/kpt) and
[kustomize](https://github.com/kubernetes-sigs/kustomize).

## Set the name

```shell script
$ kpt cfg set install/ name NAME
```

Example: `kpt cfg set install/ name example`

## Set the image

```shell script
$ kpt cfg set install/ image IMAGE
```

Example: `kpt cfg set install image gcr.io/pwittrock/example:v0.1.0`

Edit the `Makefile` with the updated image.

## Run the aggregated apiserver in a localcluster

```shell script
$ kubectl config get-contexts                                                          
CURRENT   NAME                 CLUSTER          AUTHINFO         NAMESPACE
*         docker-desktop       docker-desktop   docker-desktop   
```

```shell script
$ make install
```

## List the APIs

```shell script
$ kubectl api-versions
$ kubectl api-resources
```
