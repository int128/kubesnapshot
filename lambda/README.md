# kubesnapshot on AWS Lambda

kubesnapshot can be run as a Lambda function.


## Deploy from local

Install the [AWS SAM CLI](https://github.com/awslabs/aws-sam-cli):

```sh
easy_install --user pip
pip install --user aws-sam-cli
```

Deploy a Lambda function:

```sh
export AWS_REGION=us-west-2
export KUBE_CLUSTER_NAME=hello.k8s.local
make bucket
make deploy
```


## Deploy from AWS CodeBuild

Open CodeBuild: https://console.aws.amazon.com/codebuild.

Create a project with the following:

- Name: `kubesnapshot`
- Source Provider: GitHub
- Repository URL: `https://github.com/int128/kubesnapshot`
- Runtime: Golang/1.10
- Buildspec name: `lambda/buildspec.yml`
- Environment Variables:
    - `KUBE_CLUSTER_NAME`: Name of the Kubernetes cluster (e.g. `hello.k8s.local`)

Open IAM: https://console.aws.amazon.com/iam/home#/roles/codebuild-kubesnapshot-service-role.

Attach the `AdministratorAccess` policy to the role.

Then start the build in CodeBuild.
