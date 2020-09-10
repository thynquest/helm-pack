# helm-pack

[![Build Status](https://travis-ci.org/thynquest/helm-pack.svg?branch=master)](https://travis-ci.org/thynquest/helm-pack)

This is a helm plugin designed to be able to set properties before packaging. it has the same flags as the helm package command except that you can define properties when packaging (`helm pack myfolder --set myprop=myval`)


# Description

The goal of this helm plugin is to be able to set properties before packaging. Suppose that you have the following pipeline:

* write code for you app 
* push your code to your repo which will build a docker image automatically tagged (suppose the tag is `mytag1234`)
* after building your image you automatically trigger another pipeline associated to the helm chart repo which contains a `values.yaml` with some unset value
````yaml
deployment:
 name: "deployment_name"
 replicas: 1
 version: imageversion  <= UNSET VALUE
  container:
  image: "my.repo.url/path/image"
  configMapRef:
   name: "deployment-config"
  limits:
   memory: "1Gi"
   cpu: "1000m"
  request:
   memory: "500Mi"
   cpu: "500m"
  imagePullPolicy: "Always"  
````
* then in order to create our helm package according to the docker image previously built we need to inject the version value (`mytag1234`) **during the packaging process**. so we won't have to make `helm install mypackage --set deployment.version=mytag1234` during the installation. but instead `helm package . --set deployment.version=mytag1234` so that way we during the installation we just have to execute `helm install mypackage` knowing that the version dependency has been resolved earlier during the packaging process.
  
**IMPORTANT NOTE:** this is still a work in progress but if you have any issues/remarks let me know.

# Install

The plugin will be downloaded from github

````sh
$ helm plugin install https://github.com/thynquest/helm-pack.git
````

# Usage

the same options of the package command apply here with the possibilty to set property

`helm pack myfolder --set myproperty=myvalue`
