package backup

import (
	"io"
	"text/template"

	"github.com/int128/kubesnapshot/awsk8s"
)

// Operations represents a set of operations for backup.
type Operations struct {
	VolumesToCreateSnapshot awsk8s.EBSVolumes
	SnapshotsToDelete       awsk8s.EBSSnapshots
}

// ComputeOperations returns a operations for the backup.
func (b *Backup) ComputeOperations(volumes awsk8s.EBSVolumes, snapshots awsk8s.EBSSnapshots) *Operations {
	ops := Operations{}
	for _, volume := range volumes {
		if volume.Name != "" {
			ops.VolumesToCreateSnapshot = append(ops.VolumesToCreateSnapshot, volume)

			snapshotsOfVolume := snapshots.FindByName(volume.Name)
			snapshotsOfVolume.SortByLatest()
			snapshotsToDelete := snapshotsOfVolume.TrimHead(b.RetainSnapshots - 1)
			ops.SnapshotsToDelete = append(ops.SnapshotsToDelete, snapshotsToDelete...)
		}
	}
	return &ops
}

// Print shows details of operations.
func (o *Operations) Print(w io.Writer) {
	opsTemplate.Execute(w, o)
}

var opsTemplate = template.Must(template.New("ops").Parse(`
Each snapshot of the following {{len .VolumesToCreateSnapshot}} volumes will be created:
{{- range $_, $e := .VolumesToCreateSnapshot}}
  {{$e.ID}}
  {{- " "}}{{$e.Name}}
  {{- " "}}{{$e.PersistentVolumeTags.PersistentVolumeClaimNamespace}}
  {{- " "}}{{$e.PersistentVolumeTags.PersistentVolumeClaimName}}
{{- end}}

The following {{len .SnapshotsToDelete}} snapshots will be deleted:
{{- range $_, $e := .SnapshotsToDelete}}
  {{$e.ID}}
  {{- " "}}[{{$e.StartTime}}]
  {{- " "}}{{$e.Name}}
  {{- " "}}{{$e.PersistentVolumeTags.PersistentVolumeClaimNamespace}}
  {{- " "}}{{$e.PersistentVolumeTags.PersistentVolumeClaimName}}
{{- end}}

`))
