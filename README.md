# kubesnapshot

A command to create snapshots for EBS volumes owned by the Kubernetes cluster on AWS.
Written in Go.


## Run on Local

```sh
# Configure your credentials
aws configure --profile hello
export AWS_PROFILE=hello

# Install
go install github.com/int128/kubesnapshot

# Run
export AWS_REGION=us-west-2
export KUBE_CLUSTER_NAME=hello.k8s.local
kubesnapshot --dry-run
```

Options:

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


## Deploy to Lambda from Local

```sh
# Install AWS SAM CLI
easy_install --user pip
pip install --user aws-sam-cli

# Deploy
export AWS_REGION=us-west-2
export KUBE_CLUSTER_NAME=hello.k8s.local
cd lambda
make bucket
make deploy
```


## Deploy to Lambda from CodeBuild

Open CodeBuild: https://console.aws.amazon.com/codebuild.

Create a project with the following:

- Name: `kubesnapshot`
- Source Provider: GitHub
- Repository URL: `https://github.com/int128/kubesnapshot`
- Runtime: Golang/1.10
- Environment Variables:
    - `KUBE_CLUSTER_NAME`: Name of the Kubernetes cluster (e.g. `hello.k8s.local`)

Open IAM: https://console.aws.amazon.com/iam/home#/roles/codebuild-kubesnapshot-service-role.

Attach the `AdministratorAccess` policy to the role.

Then start the build in CodeBuild.


## Contributions

This is an open source software.
Please feel free to open issues or pull requests.
