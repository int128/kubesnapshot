package backup

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/int128/kubesnapshot/cluster"
)

// Backup represents a backup for the cluster.
type Backup struct {
	DryRun      bool // Dry-run flag
	RetainCount int  // Number of snapshots to retain
}

// Do performs the backup.
func (b *Backup) Do(service *cluster.Service) error {
	volumes, err := service.ListOwnedEBSVolumes()
	if err != nil {
		return err
	}
	log.Printf("Found %d volumes owned by the cluster %s", len(volumes), service.ClusterName)
	snapshots, err := service.ListOwnedEBSSnapshots()
	if err != nil {
		return err
	}
	log.Printf("Found %d snapshots owned by the cluster %s", len(snapshots), service.ClusterName)

	for _, volume := range volumes {
		tags := cluster.Tags(volume.Tags)
		name := tags.FindName()
		if name != "" {
			var takingID string
			if !b.DryRun {
				taking, err := service.CreateEBSSnapshot(volume)
				if err != nil {
					//FIXME
				}
				takingID = aws.StringValue(taking.SnapshotId)
			}
			log.Printf("Creating the snapshot %s from the volume: ID=%s, Name=%s, PV=%s, PVC=%s",
				takingID,
				aws.StringValue(volume.VolumeId), name, tags.FindPVName(), tags.FindPVCName())

			deleting := snapshots.FindByName(name).SortByLatest().TrimHead(b.RetainCount)
			for _, d := range deleting {
				tags := cluster.Tags(d.Tags)
				if !b.DryRun {
					if err := service.DeleteEBSSnapshot(d); err != nil {
						//FIXME
					}
				}
				log.Printf("Deleting the snapshot: ID=%s, StartTime=%s, Name=%s",
					aws.StringValue(d.SnapshotId),
					aws.TimeValue(d.StartTime).String(),
					tags.FindName())
			}
		}
	}
	return nil
}
