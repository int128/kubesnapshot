package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/int128/kubesnapshot/aws_k8s"
	"github.com/int128/kubesnapshot/backup"
	"github.com/int128/kubesnapshot/options"
)

func handler(ctx context.Context) error {
	opts, err := options.Parse()
	if err != nil {
		options.WriteHelp(os.Stderr)
		return err
	}
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	svc := aws_k8s.New(sess)
	b := &backup.Backup{
		DryRun:          opts.DryRun,
		ClusterName:     aws_k8s.ClusterName(opts.ClusterName),
		RetainSnapshots: opts.RetainSnapshots,
	}
	log.Printf("Backup the cluster %+v", b)
	return b.Do(svc)
}

func main() {
	lambda.Start(handler)
}
