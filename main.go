package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/int128/kubesnapshot/awsk8s"
	"github.com/int128/kubesnapshot/backup"
	"github.com/int128/kubesnapshot/options"
)

func main() {
	opts, err := options.Parse()
	if err != nil {
		if !options.IsErrHelp(err) {
			log.Print(err)
		}
		options.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := awsk8s.New(sess)
	b := &backup.Backup{
		DryRun:          opts.DryRun,
		ClusterName:     awsk8s.ClusterName(opts.ClusterName),
		RetainSnapshots: opts.RetainSnapshots,
	}
	log.Printf("Backup the cluster %+v", b)
	if err := b.Do(svc); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
