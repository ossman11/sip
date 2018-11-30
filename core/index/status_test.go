package index

import "testing"

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want string
	}{
		{
			name: "Valid: Idle",
			s:    Idle,
			want: "idle",
		},
		{
			name: "Valid: Scanning",
			s:    Scanning,
			want: "scanning",
		},
		{
			name: "Valid: Indexing",
			s:    Indexing,
			want: "indexing",
		},
		{
			name: "Valid: Indexed",
			s:    Indexed,
			want: "indexed",
		},
		{
			name: "Invalid: Idle",
			s:    -10,
			want: "idle",
		},
		{
			name: "Invalid: unknown",
			s:    999,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
