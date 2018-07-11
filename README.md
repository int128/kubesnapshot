# kubesnapshot

`kubesnapshot` creates snapshots for EBS volumes owned by the Kubernetes cluster.

## Getting Started

Configure your AWS credentials:

```sh
aws configure --profile hello
export AWS_PROFILE=hello
```

Build and run:

```sh
# Build
make

# Run
export AWS_REGION=us-west-2
export KUBE_CLUSTER_NAME=hello.k8s.local
./kubesnapshot --dry-run
```

### Options

```
Usage:
  kubesnapshot [OPTIONS]

Application Options:
      --dry-run            Dry-run flag [$DRY_RUN]
      --kube-cluster-name= Kubernetes cluster name [$KUBE_CLUSTER_NAME]
      --retain-snapshots=  Number of snapshots to retain (default: 5) [$RETAIN_SNAPSHOTS]

Help Options:
  -h, --help               Show this help message
Usage:
  main [OPTIONS]

Application Options:
  --dry-run              Dry-run flag [$DRY_RUN]
  --kube-cluster-name=   Kubernetes cluster name [$KUBE_CLUSTER_NAME]
  --retain-snapshots=    Number of snapshots to retain [$RETAIN_SNAPSHOTS]
```

## Deploy to Lambda

Install the [AWS SAM CLI](https://github.com/awslabs/aws-sam-cli):

```sh
easy_install --user pip
pip install --user aws-sam-cli
```

Deploy:

```sh
export AWS_REGION=us-west-2
cd lambda
make deploy
```

You can change schedule in [`lambda/template.yaml`](lambda/template.yaml).
By default the function will be executed at 00:00 UTC everyday.

## Contributions

This is an open source software.
Please feel free to open issues or pull requests.
