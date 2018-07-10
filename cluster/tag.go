package cluster

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// OwnedTagFilter returns a filter for the cluster owned tag.
func (s *Service) OwnedTagFilter() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:kubernetes.io/cluster/%s", s.ClusterName)),
		Values: aws.StringSlice([]string{"owned"}),
	}
}

// Tags represent tags of AWS resource.
type Tags []*ec2.Tag

// FindName returns the name tag or empty.
func (t Tags) FindName() string {
	for _, tag := range t {
		if aws.StringValue(tag.Key) == "Name" {
			return aws.StringValue(tag.Value)
		}
	}
	return ""
}

// FindPVName returns name of the Persistent Volume or empty.
func (t Tags) FindPVName() string {
	for _, tag := range t {
		if aws.StringValue(tag.Key) == "kubernetes.io/created-for/pv/name" {
			return aws.StringValue(tag.Value)
		}
	}
	return ""
}

// FindPVCName returns name of the Persistent Volume Claim or empty.
func (t Tags) FindPVCName() string {
	for _, tag := range t {
		if aws.StringValue(tag.Key) == "kubernetes.io/created-for/pvc/name" {
			return aws.StringValue(tag.Value)
		}
	}
	return ""
}
