package cluster

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestEBSSnapshots_FindByName(t *testing.T) {
	s1 := &EBSSnapshot{Name: "a"}
	s2 := &EBSSnapshot{Name: "b"}
	s3 := &EBSSnapshot{Name: "b"}
	snapshots := EBSSnapshots{s1, s2, s3}

	a := snapshots.FindByName("a")
	if len(a) != 1 || a[0] != s1 {
		t.Errorf("FindByName(a) wants [s1] but %+v", a)
	}

	b := snapshots.FindByName("b")
	if len(b) != 2 || b[0] != s2 || b[1] != s3 {
		t.Errorf("FindByName(b) wants [s2, s3] but %+v", b)
	}

	z := snapshots.FindByName("z")
	if len(z) != 0 {
		t.Errorf("FindByName(z) wants an empty slice but %+v", z)
	}
}

func TestEBSSnapshots_TrimHead(t *testing.T) {
	s1 := &EBSSnapshot{Name: "a"}
	s2 := &EBSSnapshot{Name: "b"}
	s3 := &EBSSnapshot{Name: "b"}
	snapshots := EBSSnapshots{s1, s2, s3}

	matrix := []struct {
		n        int
		expected EBSSnapshots
	}{
		{0, snapshots},
		{1, EBSSnapshots{s2, s3}},
		{2, EBSSnapshots{s3}},
		{3, EBSSnapshots{}},
		{4, EBSSnapshots{}},
	}
	for _, m := range matrix {
		t.Run(fmt.Sprintf("n=%d", m.n), func(t *testing.T) {
			actual := snapshots.TrimHead(m.n)
			if !reflect.DeepEqual(m.expected, actual) {
				t.Errorf("TrimHead wants %+v but %+v", m.expected, actual)
			}
		})
	}
}

func TestEBSSnapshots_SortByLatest(t *testing.T) {
	s1 := &EBSSnapshot{StartTime: time.Date(2018, 4, 2, 12, 34, 56, 0, time.UTC)}
	s2 := &EBSSnapshot{StartTime: time.Date(2018, 3, 1, 12, 34, 56, 0, time.UTC)}
	s3 := &EBSSnapshot{StartTime: time.Date(2018, 5, 6, 12, 34, 56, 0, time.UTC)}
	snapshots := EBSSnapshots{s1, s2, s3}
	snapshots.SortByLatest()
	expected := EBSSnapshots{s3, s1, s2}
	if !reflect.DeepEqual(expected, snapshots) {
		t.Errorf("TrimHead wants %+v but %+v", expected, snapshots)
	}
}
