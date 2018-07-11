package cluster

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// OwnedTagFilter returns a filter for the cluster owned tag.
func (k KubernetesClusterName) OwnedTagFilter() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:kubernetes.io/cluster/%s", k)),
		Values: aws.StringSlice([]string{"owned"}),
	}
}

// Tags represent tags of AWS resource.
type Tags []*ec2.Tag

// FindByKey returns the value for key, or empty.
func (t Tags) FindByKey(key string) string {
	for _, tag := range t {
		if aws.StringValue(tag.Key) == key {
			return aws.StringValue(tag.Value)
		}
	}
	return ""
}

// FindName returns the name tag or empty.
func (t Tags) FindName() string {
	return t.FindByKey("Name")
}

// PersistentVolume returns the PersistentVolumeTags.
func (t Tags) PersistentVolume() PersistentVolumeTags {
	return PersistentVolumeTags{
		PersistentVolumeName:           t.FindByKey("kubernetes.io/created-for/pv/name"),
		PersistentVolumeClaimNamespace: t.FindByKey("kubernetes.io/created-for/pvc/namespace"),
		PersistentVolumeClaimName:      t.FindByKey("kubernetes.io/created-for/pvc/name"),
	}
}

// PersistentVolumeTags represents reference to a dynamic provisioned PV and PVC.
type PersistentVolumeTags struct {
	PersistentVolumeName           string
	PersistentVolumeClaimNamespace string
	PersistentVolumeClaimName      string
}
