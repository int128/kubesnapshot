package options

import (
	"fmt"
	"io"

	flags "github.com/jessevdk/go-flags"
)

// Options represents command options from arguments or environment variables.
type Options struct {
	DryRun          bool   `long:"dry-run" env:"DRY_RUN" description:"Dry-run flag"`
	ClusterName     string `long:"kube-cluster-name" env:"KUBE_CLUSTER_NAME" required:"1" description:"Kubernetes cluster name"`
	RetainSnapshots int    `long:"retain-snapshots" env:"RETAIN_SNAPSHOTS" default:"5" description:"Number of snapshots to retain"`
}

// Parse returns an Options.
func Parse() (*Options, error) {
	var opts Options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	args, err := parser.Parse()
	switch {
	case err != nil:
		return nil, err
	case len(args) > 0:
		return nil, fmt.Errorf("No argument required")
	case opts.RetainSnapshots < 1:
		return nil, fmt.Errorf("RetainSnapshots must be 1 or more but %d", opts.RetainSnapshots)
	}
	return &opts, nil
}

// IsErrHelp returns true if --help or -h is set.
func IsErrHelp(err error) bool {
	e, ok := err.(*flags.Error)
	if !ok {
		return false
	}
	return e.Type == flags.ErrHelp
}

// WriteHelp writes the help for options.
func WriteHelp(w io.Writer) {
	var opts Options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	parser.WriteHelp(w)
}
