package store

import "testing"

func TestPaginate(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	tests := []struct {
		name      string
		items     []int
		page      int
		limit     int
		wantSlice []int
		wantTotal int
	}{
		{"empty input", nil, 1, 10, []int{}, 0},
		{"first page", items, 1, 2, []int{1, 2}, 5},
		{"middle page", items, 2, 2, []int{3, 4}, 5},
		{"last partial page", items, 3, 2, []int{5}, 5},
		{"page beyond total", items, 4, 2, []int{}, 5},
		{"limit exceeds total", items, 1, 100, []int{1, 2, 3, 4, 5}, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, total := paginate(tc.items, tc.page, tc.limit)
			if total != tc.wantTotal {
				t.Errorf("total = %d, want %d", total, tc.wantTotal)
			}
			if len(got) != len(tc.wantSlice) {
				t.Fatalf("len = %d (%v), want %d (%v)", len(got), got, len(tc.wantSlice), tc.wantSlice)
			}
			for i := range got {
				if got[i] != tc.wantSlice[i] {
					t.Errorf("got[%d] = %d, want %d", i, got[i], tc.wantSlice[i])
				}
			}
		})
	}
}

// paginate must never return a nil slice for an out-of-range page, so JSON
// encoders emit [] rather than null.
func TestPaginate_NonNilOnEmpty(t *testing.T) {
	got, _ := paginate([]int{1, 2}, 5, 2)
	if got == nil {
		t.Error("expected non-nil empty slice for out-of-range page")
	}
}
