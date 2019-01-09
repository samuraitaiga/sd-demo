# Requirements
* golang 1.11
* debian (64bit)
* GCP billable project

# Setup

## Mac
```
$ brew update
$ brew install goenv
$ goenv init -
$ goenv install 1.11.2
$ goenv global 1.11.2

# optional add this line into ~/.bash_profile
eval "$(goenv init -)"
```

## CloudShell
```
nothing to do
```

# Setup

## Enable GCP APIs
enable APIs used in this handson

* Cloud Build
* Cloud Source Repository
* Cloud Container Registry
* Google Kubernetes Engine
* Stackdriver Trace, Profiler, Debugger, Logging

## GCP IAM
1. create IAM user "dohandson" and assign "Project Editor" access.
1. download "handson" user's credential to your CloudShell or Local machine.
1. assign "Source Repository Reader" and "Kubernetes Engine Admin" to PROJECT_NUMBER@cloudbuild.gserviceaccount.com account(Cloud Build default service account).

## Create GKE Cluster

Don't forget to set these two options

* Allow full access to all Cloud APIs
* Enable VPC-native (using alias IP)

```
gcloud beta container --project "samuraitaiga-demo" clusters create "k8s-devops-handson" --zone "asia-northeast1-c" --username "admin" --cluster-version "1.11.4-gke.8" --machine-type "n1-standard-1" --image-type "COS" --disk-type "pd-standard" --disk-size "100" --scopes "https://www.googleapis.com/auth/cloud-platform" --num-nodes "3" --enable-cloud-logging --enable-cloud-monitoring --enable-ip-alias --network "projects/samuraitaiga-demo/global/networks/default" --subnetwork "projects/samuraitaiga-demo/regions/asia-northeast1/subnetworks/default" --default-max-pods-per-node "110" --addons HorizontalPodAutoscaling,HttpLoadBalancing --enable-autoupgrade --enable-autorepair
```

## Beego
```
$ go get -u github.com/astaxie/beego
$ go get -u github.com/astaxie/beego/context
$ go get -u github.com/beego/bee
$ GOPATH=$HOME/go ~/go/bin/bee new devops-handson
$ export GOOGLE_CLOUD_PROJECT=FIX-ME # replace FIX-ME to your GCP project id
$ export GOOGLE_APPLICATION_CREDENTIALS=/PATH/TO/IAM-CREDENTIAL.json # replace /PATH/TO/IAM-CREDENTIAL.json to actual path of iam credential.
```

## Stackdriver Trace

Trace - https://cloud.google.com/trace/docs/setup/go

```
$ go get -u go.opencensus.io/trace
$ go get -u contrib.go.opencensus.io/exporter/stackdriver
```

## Stackdriver Profiler

https://cloud.google.com/profiler/docs/profiling-go
https://cloud.google.com/profiler/docs/profiling-external

```
$ go get -u cloud.google.com/go/profiler
```


## Stackdriver Debugger

Only works 64 bit Debian

```
$ go get -u cloud.google.com/go/cmd/go-cloud-debug-agent
```

## Stacdriver Logging

Customize Known Fields
https://cloud.google.com/logging/docs/agent/configuration#special-fields

## Build

gcflags must be assigned due to restriction of Stackdriver Debugger.

```
$ go build -gcflags=all='-N -l' main.go
```

GOOGLE_APPLICATION_CREDENTIALS=./gcp-credentials/auth.json go-cloud-debug-agent -sourcecontext=./source-context.json -appmodule=devops-handson -appversion=1.0 -- ./main


# Memo
## Cloud build
triggers from GSR automatically clone files into /workspace/. user doesn't need to clone by themselves
