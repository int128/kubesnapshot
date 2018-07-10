package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/int128/kubesnapshot/backup"
	"github.com/int128/kubesnapshot/cluster"
)

func main() {
	clusterName := os.Getenv("KUBE_CLUSTER_NAME")
	if clusterName == "" {
		log.Fatalf("KUBE_CLUSTER_NAME is not set")
	}
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := cluster.New(clusterName, sess)
	b := &backup.Backup{
		DryRun:      true,
		RetainCount: 5,
	}
	if err := b.Do(svc); err != nil {
		fmt.Printf("Could not take snapshots: %s", err)
	}
}
