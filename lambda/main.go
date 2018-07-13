package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/int128/kubesnapshot/awsk8s"
	"github.com/int128/kubesnapshot/backup"
	"github.com/int128/kubesnapshot/options"
)

func handler(ctx context.Context) error {
	opts, err := options.Parse()
	if err != nil {
		options.WriteHelp(os.Stderr)
		return err
	}
	svc, err := awsk8s.New()
	if err != nil {
		return err
	}
	b := &backup.Backup{
		DryRun:          opts.DryRun,
		ClusterName:     awsk8s.ClusterName(opts.ClusterName),
		RetainSnapshots: opts.RetainSnapshots,
	}
	log.Printf("Backup the cluster %+v", b)
	return b.Do(svc)
}

func main() {
	lambda.Start(handler)
}
