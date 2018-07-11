package cluster

import (
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EBSVolume represents an EBS volume.
type EBSVolume struct {
	ID                   string
	Name                 string
	PersistentVolumeTags PersistentVolumeTags
	tags                 Tags
}

func awsEBSVolume(v *ec2.Volume) *EBSVolume {
	return &EBSVolume{
		ID:                   aws.StringValue(v.VolumeId),
		Name:                 Tags(v.Tags).FindName(),
		PersistentVolumeTags: Tags(v.Tags).PersistentVolume(),
		tags:                 v.Tags,
	}
}

// EBSVolumes represents EBS volumes.
type EBSVolumes []*EBSVolume

// EBSSnapshot represents an EBS snapshot.
type EBSSnapshot struct {
	ID                   string
	Name                 string
	StartTime            time.Time
	PersistentVolumeTags PersistentVolumeTags
	tags                 Tags
}

func awsEBSSnapshot(s *ec2.Snapshot) *EBSSnapshot {
	return &EBSSnapshot{
		ID:                   aws.StringValue(s.SnapshotId),
		Name:                 Tags(s.Tags).FindName(),
		StartTime:            aws.TimeValue(s.StartTime),
		PersistentVolumeTags: Tags(s.Tags).PersistentVolume(),
		tags:                 s.Tags,
	}
}

// EBSSnapshots represents EBS snapshots.
type EBSSnapshots []*EBSSnapshot

// FindByName returns EBS snapshots of the name.
func (snapshots EBSSnapshots) FindByName(name string) EBSSnapshots {
	m := make(EBSSnapshots, 0)
	for _, s := range snapshots {
		if s.Name == name {
			m = append(m, s)
		}
	}
	return m
}

// SortByLatest returns a slice sorted by the StartTime descending.
func (snapshots EBSSnapshots) SortByLatest() EBSSnapshotsSortedByLatest {
	a := make(EBSSnapshotsSortedByLatest, len(snapshots))
	for i, s := range snapshots {
		a[i] = s
	}
	sort.Sort(a)
	return a
}

// EBSSnapshotsSortedByLatest is a slice sorted by the StartTime descending.
// This implements the sort interface.
type EBSSnapshotsSortedByLatest EBSSnapshots

func (s EBSSnapshotsSortedByLatest) Len() int      { return len(s) }
func (s EBSSnapshotsSortedByLatest) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s EBSSnapshotsSortedByLatest) Less(i, j int) bool {
	return s[i].StartTime.After(s[j].StartTime)
}

// TrimHead returns a slice without heading n items.
func (s EBSSnapshotsSortedByLatest) TrimHead(n int) EBSSnapshotsSortedByLatest {
	if n > len(s) {
		return EBSSnapshotsSortedByLatest{}
	}
	return s[n:]
}
