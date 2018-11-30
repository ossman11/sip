package index

import "testing"

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			name: "Valid: Unknown",
			t:    Unknown,
			want: "unknown",
		},
		{
			name: "Valid: Local",
			t:    Local,
			want: "local",
		},
		{
			name: "Valid: Redirect",
			t:    Redirect,
			want: "redirect",
		},
		{
			name: "Invalid: unknown (to Low)",
			t:    -10,
			want: "unknown",
		},
		{
			name: "Invalid: unknown (to High)",
			t:    999,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
