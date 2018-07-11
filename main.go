package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/int128/kubesnapshot/backup"
	"github.com/int128/kubesnapshot/cluster"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	DryRun      bool   `long:"dry-run" env:"DRY_RUN" description:"Dry-run flag"`
	ClusterName string `long:"kube-cluster-name" env:"KUBE_CLUSTER_NAME" required:"1" description:"Kubernetes cluster name"`
	RetainCount int    `long:"retain-count" env:"RETAIN_COUNT" default:"5" description:"Number of snapshots to retain"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	if len(args) > 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	if opts.RetainCount < 1 {
		log.Fatalf("RetainCount must be 1 or more but %d", opts.RetainCount)
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := cluster.New(opts.ClusterName, sess)
	b := &backup.Backup{
		DryRun:      opts.DryRun,
		RetainCount: opts.RetainCount,
	}
	log.Printf("Backup the cluster %s with %+v", opts.ClusterName, b)
	if err := b.Do(svc); err != nil {
		log.Fatalf("Error: %s", err)
	}
}