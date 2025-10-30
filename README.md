# Overview

This repo contains artifacts for deploying Odigos on Google Cloud
throught the GCP Marketplace.

# Installation

## Command line instructions

### Prerequisites

#### Set up command line tools

You'll need the following tools in your development environment. If you are
using Cloud Shell, `gcloud`, `kubectl`, Docker, and Git are installed in your
environment by default.

- [gcloud](https://cloud.google.com/sdk/gcloud/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [docker](https://docs.docker.com/install/)

Configure `gcloud` as a Docker credential helper:

```shell
gcloud auth configure-docker
```

#### Create a Google Kubernetes Engine cluster

Create a cluster from the command line. If you already have a cluster that you
want to use, this step is optional.

```shell
export CLUSTER=sample-app-cluster
export ZONE=us-west1-a

gcloud container clusters create "$CLUSTER" --zone "$ZONE"
```

#### Configure kubectl to connect to the cluster

```shell
gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE"
```

#### Clone this repo

Clone this repo:

```shell
git clone --recursive https://github.com/odigos-io/gcp-artifacts.git
```
#### Install the Application resource definition

An Application resource is a collection of individual Kubernetes components,
such as Services, Deployments, and so on, that you can manage as a group.

To set up your cluster to understand Application resources, run the following
command:

```shell
kubectl apply -f "https://raw.githubusercontent.com/GoogleCloudPlatform/marketplace-k8s-app-tools/master/crd/app-crd.yaml"
```

You need to run this command once for each cluster.

The Application resource is defined by the
[Kubernetes SIG-apps](https://github.com/kubernetes/community/tree/master/sig-apps)
community. The source code can be found on
[github.com/kubernetes-sigs/application](https://github.com/kubernetes-sigs/application).

### Install the Application

#### Configure the app with environment variables

Choose an instance name and
[namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
for the app. In most cases, you can use the `default` namespace.

```shell
export APP_INSTANCE_NAME=sample-app-1
export NAMESPACE=default
```

Set the on prem token for enterprise if applicable::

```shell
export ODIGOS_ON_PREM_TOKEN=<your token>
```

Configure the container images:

```shell
export ODIGOS_AUTOSCALER_IMAGE="gcr.io/odigos-public/odigos/odigos-autoscaler:1.8.6"
export ODIGOS_COLLECTOR_IMAGE="gcr.io/odigos-public/odigos/odigos-collector:1.8.6"
export ODIGOS_INIT_CONTAINER_IMAGE="gcr.io/odigos-public/odigos/odigos-init-container:1.8.6"
export ODIGOS_SCHEDULER_IMAGE="gcr.io/odigos-public/odigos/odigos-scheduler:1.8.6"
export ODIGOS_INSTRUMENTOR_IMAGE="gcr.io/odigos-public/odigos/odigos-instrumentor:1.8.6"
export ODIGOS_ODIGLET_IMAGE="gcr.io/odigos-public/odigos/odigos-odiglet:1.8.6"
export ODIGOS_UI_IMAGE="gcr.io/odigos-public/odigos/odigos-ui:1.8.6"
export ODIGOS_ENTERPRISE_ODIGLET_IMAGE="gcr.io/odigos-public/odigos/odigos-enterprise-odiglet:1.8.6"
export ODIGOS_ENTERPRISE_INSTRUMENTOR_IMAGE="gcr.io/odigos-public/odigos/odigos-enterprise-instrumentor:1.8.6"
```
#### Expand the manifest template

Use `envsubst` to expand the template. We recommend that you save the expanded
manifest file for future updates to the application.

```shell
awk 'FNR==1 {print "---"}{print}' manifest/* \
  | envsubst '$APP_INSTANCE_NAME $NAMESPACE $ODIGOS_ON_PREM_TOKEN $ODIGOS_AUTOSCALER_IMAGE $ODIGOS_COLLECTOR_IMAGE $ODIGOS_INIT_CONTAINER_IMAGE $ODIGOS_SCHEDULER_IMAGE $ODIGOS_INSTRUMENTOR_IMAGE $ODIGOS_ODIGLET_IMAGE $ODIGOS_UI_IMAGE $ODIGOS_ENTERPRISE_ODIGLET_IMAGE $ODIGOS_ENTERPRISE_INSTRUMENTOR_IMAGE' \
  > "${APP_INSTANCE_NAME}_manifest.yaml"
```

#### Apply the manifest to your Kubernetes cluster

Use `kubectl` to apply the manifest to your Kubernetes cluster.

```shell
kubectl apply -f "${APP_INSTANCE_NAME}_manifest.yaml" --namespace "${NAMESPACE}"
```

#### View the app in the Google Cloud Console

To get the Console URL for your app, run the following command:

```shell
echo "https://console.cloud.google.com/kubernetes/application/${ZONE}/${CLUSTER}/${NAMESPACE}/${APP_INSTANCE_NAME}"
```

To view the app, open the URL in your browser.

# Using the app

See [the official Odigos documentation](https://docs.odigos.io) for info on using Odigos.
