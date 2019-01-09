# set GCP project ID

## Get project name and id
```bash
gcloud projects list
```

## set GCP project id in environment variable
replace FIXME to your GCP Project ID
```bash
export GOOGLE_CLOUD_PROJECT=FIXME
```

## set default GCP project
replace FIXME to your GCP Project ID
```bash
gcloud config set project FIXME
```

# Enable required APIs

```bash
gcloud services enable cloudbuild.googleapis.com sourcerepo.googleapis.com containerregistry.googleapis.com container.googleapis.com cloudtrace.googleapis.com cloudprofiler.googleapis.com logging.googleapis.com
```

# Create ServiceAccount and assign IAM Roles

## Create Service Account for this example
```bash
gcloud iam service-accounts create dohandson --display-name "DevOps HandsOn Service Account"
```

## Create and download Key file for Service Account
```bash
gcloud iam service-accounts keys create auth.json --iam-account=dohandson@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com --key-file-type=json
````

## Add IAM Permission to Service Account

### Cloud Profiler Agent role
```bash
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT  --member serviceAccount:dohandson@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com --role roles/cloudprofiler.agent
```

### Cloud Trace Agent role
```bash
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT  --member serviceAccount:dohandson@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com --role roles/cloudtrace.agent
```

### Cloud Monitoring Metric Writer role
```bash
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT  --member serviceAccount:dohandson@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com --role roles/monitoring.metricWriter
```

### Cloud Monitoring Metadata Writer role
```bash
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT  --member serviceAccount:dohandson@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com --role roles/stackdriver.resourceMetadata.writer
```

### (Optional)  CloudDebugger Agent role
```bash
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT  --member serviceAccount:dohandson@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com --role roles/clouddebugger.agent
```

# Place IAM key file to project directory

```bash
mv auth.json gcp-credentials/auth.json
```

# Create GKE Cluster

```bash
gcloud container clusters create "k8s-devops-handson"  \
--zone "asia-northeast1-c" \
--enable-autoupgrade \
--enable-autorepair \
--username "admin" \
--cluster-version "1.11.4-gke.8" \
--machine-type "n1-standard-1" \
--image-type "COS" \
--disk-type "pd-standard" \
--disk-size "100" \
--scopes "https://www.googleapis.com/auth/cloud-platform" \
--num-nodes "3" \
--enable-cloud-logging --enable-cloud-monitoring \
--enable-ip-alias \
--network "projects/$GOOGLE_CLOUD_PROJECT/global/networks/default" \
--subnetwork "projects/$GOOGLE_CLOUD_PROJECT/regions/asia-northeast1/subnetworks/default" \
--addons HorizontalPodAutoscaling,HttpLoadBalancing
```

# Obtain auth info from GKE cluster

```bash
gcloud container clusters get-credentials k8s-devops-handson --zone asia-northeast1-c --project $GOOGLE_CLOUD_PROJECT
```

# Replace FIXME to your GCP project id

## in exapmle app config file
```bash
sed -i".org" -e "s/FIXME/$GOOGLE_CLOUD_PROJECT/g" conf/app.conf
```

## in k8s config file
```bash
sed -i".org" -e "s/FIXME/$GOOGLE_CLOUD_PROJECT/g" gke-config/deployment.yaml
```

# Build and Deploy Container

## Build container
```bash
docker build -t gcr.io/$GOOGLE_CLOUD_PROJECT/devops-handson:v1 .
```

## Push container into container registry
```bash
docker push gcr.io/$GOOGLE_CLOUD_PROJECT/devops-handson:v1
```

# Deploy container to GKE

```bash
kubectl create -f gke-config/deployment.yaml
```

# Cleanup

## Unset default GCP project
```bash
gcloud config unset project
```
