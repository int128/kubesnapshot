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

// FindByName returns EBS snapshots that have the name.
func (snapshots EBSSnapshots) FindByName(name string) EBSSnapshots {
	m := make(EBSSnapshots, 0)
	for _, s := range snapshots {
		if s.Name == name {
			m = append(m, s)
		}
	}
	return m
}

// TrimHead returns a slice without heading n items.
func (snapshots EBSSnapshots) TrimHead(n int) EBSSnapshots {
	if n > len(snapshots) {
		return EBSSnapshots{}
	}
	return snapshots[n:]
}

// SortByLatest returns a slice sorted by the StartTime descending.
func (snapshots EBSSnapshots) SortByLatest() {
	sort.Sort(sort.Reverse(byStartTime{snapshots}))
}

type byStartTime struct {
	EBSSnapshots
}

func (s byStartTime) Len() int {
	return len(s.EBSSnapshots)
}

func (s byStartTime) Swap(i, j int) {
	s.EBSSnapshots[i], s.EBSSnapshots[j] = s.EBSSnapshots[j], s.EBSSnapshots[i]
}

func (s byStartTime) Less(i, j int) bool {
	return s.EBSSnapshots[i].StartTime.Before(s.EBSSnapshots[j].StartTime)
}
