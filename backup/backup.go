package backup

import (
	"fmt"
	"log"
	"os"

	"github.com/int128/kubesnapshot/cluster"
)

// Backup represents a backup for the cluster.
type Backup struct {
	DryRun      bool // Dry-run flag
	RetainCount int  // Number of snapshots to retain (1 or more)
}

// Do performs the backup.
func (b *Backup) Do(service *cluster.Service) error {
	log.Printf("Finding EBS volumes and snaphosts in the cluster %s", service.ClusterName)
	volumes, snapshots, err := service.ListOwnedEBSVolumesAndSnapshots()
	if err != nil {
		return fmt.Errorf("Could not get EBS volumes or snaphosts: %s", err)
	}
	log.Printf("Found %d volumes owned by the cluster %s", len(volumes), service.ClusterName)
	log.Printf("Found %d snapshots owned by the cluster %s", len(snapshots), service.ClusterName)

	ops := b.ComputeOperations(volumes, snapshots)
	ops.Print(os.Stdout)
	if b.DryRun {
		log.Printf("Stop due to dry run.")
		return nil
	}

	for _, volume := range ops.VolumesToCreateSnapshot {
		_, err := service.CreateEBSSnapshot(volume)
		if err != nil {
			return err
		}
		log.Printf("Creating a snapshot from the volume %s (%s)", volume.ID, volume.Name)
	}
	for _, snapshot := range ops.SnapshotsToDelete {
		if err := service.DeleteEBSSnapshot(snapshot); err != nil {
			return err
		}
		log.Printf("Deleting the snapshot %s (%s)", snapshot.ID, snapshot.Name)
	}
	return nil
}
