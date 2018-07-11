// Package cluster provides operations on the Kubernetes cluster.
package cluster

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/sync/errgroup"
)

// KubernetesClusterName represents name of a Kubernetes cluster.
type KubernetesClusterName string

// Service provides operations on the Kubernetes cluster.
type Service struct {
	ec2 *ec2.EC2
}

// New returns a new Service.
func New(s *session.Session) *Service {
	return &Service{ec2.New(s)}
}

// ListOwnedEBSVolumes returns EBS volumes owned by the cluster.
func (s *Service) ListOwnedEBSVolumes(name KubernetesClusterName) (EBSVolumes, error) {
	out, err := s.ec2.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{name.OwnedTagFilter()},
	})
	if err != nil {
		return nil, fmt.Errorf("Could not describe EBS volumes owned by the cluster %s: %s", name, err)
	}
	volumes := make(EBSVolumes, len(out.Volumes))
	for i, v := range out.Volumes {
		volumes[i] = awsEBSVolume(v)
	}
	return volumes, nil
}

// ListOwnedEBSSnapshots returns EBS snapshots owned by the cluster.
func (s *Service) ListOwnedEBSSnapshots(name KubernetesClusterName) (EBSSnapshots, error) {
	out, err := s.ec2.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{name.OwnedTagFilter()},
	})
	if err != nil {
		return nil, fmt.Errorf("Could not describe EBS snapshots owned by the cluster %s: %s", name, err)
	}
	snapshots := make(EBSSnapshots, len(out.Snapshots))
	for i, s := range out.Snapshots {
		snapshots[i] = awsEBSSnapshot(s)
	}
	return snapshots, nil
}

// ListOwnedEBSVolumesAndSnapshots performs concurrent requests.
func (s *Service) ListOwnedEBSVolumesAndSnapshots(name KubernetesClusterName) (EBSVolumes, EBSSnapshots, error) {
	var volumes EBSVolumes
	var snapshots EBSSnapshots
	var eg errgroup.Group
	eg.Go(func() error {
		var err error
		volumes, err = s.ListOwnedEBSVolumes(name)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		snapshots, err = s.ListOwnedEBSSnapshots(name)
		if err != nil {
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}
	return volumes, snapshots, nil
}

// CreateEBSSnapshot creates a snapshot.
// This copies tags of the volume to the snapshot.
func (s *Service) CreateEBSSnapshot(volume *EBSVolume) (*EBSSnapshot, error) {
	out, err := s.ec2.CreateSnapshot(&ec2.CreateSnapshotInput{
		VolumeId:    aws.String(volume.ID),
		Description: aws.String("Managed by kubesnapshot"),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				Tags:         volume.tags,
				ResourceType: aws.String("snapshot"),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Could not create a snapshot: %s", err)
	}
	return awsEBSSnapshot(out), nil
}

// DeleteEBSSnapshot deletes the snapshot.
func (s *Service) DeleteEBSSnapshot(snapshot *EBSSnapshot) error {
	_, err := s.ec2.DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshot.ID),
	})
	if err != nil {
		return fmt.Errorf("Could not delete the snapshot %v: %s", snapshot.ID, err)
	}
	return nil
}
