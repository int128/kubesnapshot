package backup

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/int128/kubesnapshot/aws_k8s"
)

func TestBackup_ComputeOperations(t *testing.T) {
	v1 := &aws_k8s.EBSVolume{Name: "v1"}
	v2 := &aws_k8s.EBSVolume{Name: "v2"}
	volumes := aws_k8s.EBSVolumes{v1, v2}

	s11 := &aws_k8s.EBSSnapshot{Name: "v1", StartTime: time.Date(2018, 4, 2, 12, 34, 56, 0, time.UTC)}
	s12 := &aws_k8s.EBSSnapshot{Name: "v1", StartTime: time.Date(2018, 3, 1, 12, 34, 56, 0, time.UTC)}
	s13 := &aws_k8s.EBSSnapshot{Name: "v1", StartTime: time.Date(2018, 5, 6, 12, 34, 56, 0, time.UTC)}
	s21 := &aws_k8s.EBSSnapshot{Name: "v2", StartTime: time.Date(2018, 10, 2, 12, 34, 56, 0, time.UTC)}
	s22 := &aws_k8s.EBSSnapshot{Name: "v2", StartTime: time.Date(2018, 9, 1, 12, 34, 56, 0, time.UTC)}
	snapshots := aws_k8s.EBSSnapshots{s11, s12, s13, s21, s22}

	matrix := []struct {
		backup  Backup
		expects Operations
	}{
		{
			Backup{RetainSnapshots: 1},
			Operations{
				VolumesToCreateSnapshot: volumes,
				SnapshotsToDelete:       aws_k8s.EBSSnapshots{s13, s11, s12, s21, s22},
			},
		}, {
			Backup{RetainSnapshots: 2},
			Operations{
				VolumesToCreateSnapshot: volumes,
				SnapshotsToDelete:       aws_k8s.EBSSnapshots{s11, s12, s22},
			},
		}, {
			Backup{RetainSnapshots: 3},
			Operations{
				VolumesToCreateSnapshot: volumes,
				SnapshotsToDelete:       aws_k8s.EBSSnapshots{s12},
			},
		}, {
			Backup{RetainSnapshots: 4},
			Operations{
				VolumesToCreateSnapshot: volumes,
				SnapshotsToDelete:       nil,
			},
		},
	}
	for _, m := range matrix {
		t.Run(fmt.Sprintf("%+v", m.backup), func(t *testing.T) {
			actual := m.backup.ComputeOperations(volumes, snapshots)
			if !reflect.DeepEqual(&m.expects, actual) {
				t.Errorf("wants %+v but %+v", &m.expects, actual)
			}
		})
	}
}
