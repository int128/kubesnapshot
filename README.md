# kubesnapshot

A command to create snapshots for EBS volumes owned by the Kubernetes cluster on AWS.
Written in Go and deployable to Lambda.


## Purpose

You can backup the following resources by kubesnapshot:

- Dynamic provisioned volumes
- etcd volumes (kops)

A resource owned by a Kubernetes cluster should have [specific tags](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/aws/tags.go).
When you create a dynamic provisioned volume by `PersistentVolumeClaim`, the provisioner will create an EBS volume with tags for the cluster.

[kops](https://github.com/kubernetes/kops) creates 2 EBS volumes for etcd.

kubesnapshot creates an snapshot of each EBS volume.


## How it works

kubesnapshot does the following steps:

1.  It finds EBS volumes owned by the Kubernetes cluster.
    They should have a tag of key `kubernetes.io/cluster/NAME` and value `owned`.
2.  It creates an snapshot of each EBS volume owned by the cluster.
    It also copies the tags of the EBS volume to the snapshot.
3.  It keeps only specified number of snapshots (5 by default) of each EBS volume.

For example, there are the following EBS volumes and snapshots,

ID       | Tags                                 | Taken
---------|--------------------------------------|------
`vol-1`  | `Name=hello.k8s.local-dynamic-pvc-1` | -
`snap-1` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-06-30
`snap-2` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-01
`snap-3` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-02

and you run `kubesnapshot --retain-snapshots=2`,
there will be the following snapshots after run.

ID       | Tags                                 | Taken
---------|--------------------------------------|------
`vol-1`  | `Name=hello.k8s.local-dynamic-pvc-1` | -
`snap-3` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-02
`snap-4` | `Name=hello.k8s.local-dynamic-pvc-1` | 2018-07-03


## Run locally

Install:

```sh
go install github.com/int128/kubesnapshot
```

Run:

```sh
kubesnapshot --dry-run --kube-cluster-name=hello.k8s.local
```

It accepts the following options:

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

Here is an example:

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


## Deploy to Lambda

kubesnapshot can be run as a Lambda function.
See [lambda/README.md](lambda/README.md).


## Contributions

This is an open source software.
Please feel free to open issues or pull requests.
