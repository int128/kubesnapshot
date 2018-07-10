# kubesnapshot

`kubesnapshot` creates snapshots for EBS volumes owned by the Kubernetes cluster.

```sh
export AWS_PROFILE=hello
export AWS_REGION=us-west-2
export KUBE_CLUSTER_NAME=hello.k8s.local
kubesnapshot
```

## WIP

- [ ] Concurrent ops and handle errors.
- [ ] Provide options by args or env vars.
