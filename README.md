# kubesnapshot [![CircleCI](https://circleci.com/gh/int128/kubesnapshot.svg?style=shield)](https://circleci.com/gh/int128/kubesnapshot)

A command to take snapshots of EBS volumes owned by a Kubernetes cluster on AWS.
Written in Go and ready to run on Kubernetes Cron Job or AWS Lambda.


## Example

If there are the following EBS volumes and snapshots,

ID       | Tags                                 | Taken
---------|--------------------------------------|------
`vol-1`  | `Name=hello.k8s.local-dynamic-pvc-1` | -
`snap-1` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-06-30
`snap-2` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-01
`snap-3` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-02

and you run `kubesnapshot --retain-snapshots=2`,
the following snapshots should exist after run.

ID       | Tags                                 | Taken
---------|--------------------------------------|------
`vol-1`  | `Name=hello.k8s.local-dynamic-pvc-1` | -
`snap-3` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-02
`snap-4` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-03


## Getting Started

### Standalone

You can try kubesnapshot on your laptop as follows:

```sh
go install github.com/int128/kubesnapshot

export AWS_PROFILE=example
export AWS_REGION=us-west-2
kubesnapshot --dry-run --kube-cluster-name=hello.k8s.local
```

### Kubernetes Cron Job

You can deploy a daily job with the [`int128/kubesnapshot`](https://hub.docker.com/r/int128/kubesnapshot/) image to your Kubernetes cluster.

Deploy [`cronjob.yaml`](cronjob.yaml) as follows:

```sh
kubectl -n kube-system create configmap kubesnapshot --from-literal=KUBE_CLUSTER_NAME=hello.k8s.local
kubectl -n kube-system apply -f cronjob.yaml
```

You need to attach the following IAM policy to the instance role or [kube2iam](https://github.com/jtblin/kube2iam).

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeVolumes",
                "ec2:DescribeSnapshots"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "ec2:*",
            "Resource": [
                "arn:aws:ec2:*::snapshot/*",
                "arn:aws:ec2:*:*:volume/*"
            ]
        }
    ]
}
```

### AWS Lambda

You can deploy kubesnapshot as a Lambda function.
See [lambda/README.md](lambda/README.md).


## Configuration

kubesnapshot accepts the following arguments and environment variables:

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


## How it works

kubesnapshot takes snapshot of dynamic provisioned volumes.
When you create a dynamic provisioned volume using `PersistentVolumeClaim`, the EBS provisioner will create an EBS volume with [specific tags](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/aws/tags.go).

If you are using [kops](https://github.com/kubernetes/kops), it will create 2 EBS volumes with the tags for etcd.
kubesnapshot takes snapshot of them as well.

kubesnapshot backups your volumes by the following steps:

1.  Find EBS volumes owned by the Kubernetes cluster.
    They should have a tag of key `kubernetes.io/cluster/NAME` and value `owned`.
2.  Create an snapshot of each EBS volume owned by the cluster.
    And copy the tags of the EBS volume to the snapshot.
3.  Keep only specified number of snapshots (5 by default) of each EBS volume.

Here is an example output:

```
% export KUBE_CLUSTER_NAME=hello.k8s.local
% kubesnapshot --dry-run --retain-snapshots=1
2018/07/13 00:51:10 Backup the cluster &{DryRun:true ClusterName:hello.k8s.local RetainSnapshots:1}
2018/07/13 00:51:10 Finding EBS volumes and snaphosts in the cluster hello.k8s.local
2018/07/13 00:51:11 Found 4 volumes owned by the cluster hello.k8s.local
2018/07/13 00:51:11 Found 4 snapshots owned by the cluster hello.k8s.local

Each snapshot of the following 4 volumes will be created:
  vol-0a13b6e452da55927 hello.k8s.local-dynamic-pvc-f6a475fc-15e0-11e8-9a6e-06365716f47a devops jira-atlassian-jira-software
  vol-0e1b13ac39e6dcb10 hello.k8s.local-dynamic-pvc-9ff16cda-35a4-11e8-b48a-06091fbf9444 devops gitbucket
  vol-0d14174c6c11f5e3d a.etcd-events.hello.k8s.local
  vol-0ee8b3f5a6c2286c2 a.etcd-main.hello.k8s.local

The following 4 snapshots will be deleted:
  snap-03ea3c7ad888e9541 [2018-07-11 02:03:58 +0000 UTC] hello.k8s.local-dynamic-pvc-f6a475fc-15e0-11e8-9a6e-06365716f47a devops jira-atlassian-jira-software
  snap-0bec5e9ce8c939ea8 [2018-07-11 02:03:58 +0000 UTC] hello.k8s.local-dynamic-pvc-9ff16cda-35a4-11e8-b48a-06091fbf9444 devops gitbucket
  snap-0afb58168eb0a88c9 [2018-07-11 02:03:59 +0000 UTC] a.etcd-events.hello.k8s.local
  snap-0faba5c1eb7858b2f [2018-07-11 02:03:59 +0000 UTC] a.etcd-main.hello.k8s.local

2018/07/13 00:51:11 Stop due to dry run.
```


## Contributions

This is an open source software licensed under Apache 2.0.
Feel free to open issues and pull requests.
