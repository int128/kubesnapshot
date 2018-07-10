package cluster

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EBSVolumes represents EBS volumes.
type EBSVolumes []*ec2.Volume

// EBSSnapshots represents EBS snapshots.
type EBSSnapshots []*ec2.Snapshot

// FindByName returns EBS snapshots of the name.
func (s EBSSnapshots) FindByName(name string) EBSSnapshots {
	m := make(EBSSnapshots, 0)
	for _, snapshot := range s {
		name := Tags(snapshot.Tags).FindName()
		if name != "" {
			m = append(m, snapshot)
		}
	}
	return m
}

// SortByLatest returns a slice sorted by the StartTime descending.
func (s EBSSnapshots) SortByLatest() EBSSnapshotsSortedByLatest {
	a := make(EBSSnapshotsSortedByLatest, len(s))
	for i, snapshot := range s {
		a[i] = snapshot
	}
	sort.Sort(a)
	return a
}

// EBSSnapshotsSortedByLatest is a slice sorted by the StartTime descending.
// This implements the sort interface.
type EBSSnapshotsSortedByLatest []*ec2.Snapshot

func (s EBSSnapshotsSortedByLatest) Len() int      { return len(s) }
func (s EBSSnapshotsSortedByLatest) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s EBSSnapshotsSortedByLatest) Less(i, j int) bool {
	return aws.TimeValue(s[i].StartTime).After(aws.TimeValue(s[j].StartTime))
}

// TrimHead returns a slice without heading n items.
func (s EBSSnapshotsSortedByLatest) TrimHead(n int) EBSSnapshotsSortedByLatest {
	if n > len(s) {
		return EBSSnapshotsSortedByLatest{}
	}
	return s[n:]
}

// ListOwnedEBSVolumes returns EBS volumes owned by the cluster.
func (s *Service) ListOwnedEBSVolumes() (EBSVolumes, error) {
	out, err := s.ec2.DescribeVolumes(&ec2.DescribeVolumesInput{
		DryRun:  aws.Bool(s.DryRun),
		Filters: []*ec2.Filter{s.OwnedTagFilter()},
	})
	if err != nil {
		return nil, fmt.Errorf("Could not describe EBS volumes: %s", err)
	}
	return out.Volumes, nil
}

// ListOwnedEBSSnapshots returns EBS snapshots owned by the cluster.
func (s *Service) ListOwnedEBSSnapshots() (EBSSnapshots, error) {
	out, err := s.ec2.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		DryRun:  aws.Bool(s.DryRun),
		Filters: []*ec2.Filter{s.OwnedTagFilter()},
	})
	if err != nil {
		return nil, fmt.Errorf("Could not describe EBS snapshots: %s", err)
	}
	return out.Snapshots, nil
}

// CreateEBSSnapshot creates a snapshot.
// This copies the tags of volume to the snapshot.
func (s *Service) CreateEBSSnapshot(volume *ec2.Volume) (*ec2.Snapshot, error) {
	out, err := s.ec2.CreateSnapshot(&ec2.CreateSnapshotInput{
		DryRun:      aws.Bool(s.DryRun),
		VolumeId:    volume.VolumeId,
		Description: aws.String("Managed by kubesnapshot"),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				Tags:         volume.Tags,
				ResourceType: aws.String("snapshot"),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Could not create a snapshot: %s", err)
	}
	return out, nil
}

// DeleteEBSSnapshot deletes the snapshot.
func (s *Service) DeleteEBSSnapshot(snapshot *ec2.Snapshot) error {
	_, err := s.ec2.DeleteSnapshot(&ec2.DeleteSnapshotInput{
		DryRun:     aws.Bool(s.DryRun),
		SnapshotId: snapshot.SnapshotId,
	})
	if err != nil {
		return fmt.Errorf("Could not delete the snapshot %v: %s", snapshot.SnapshotId, err)
	}
	return nil
}
