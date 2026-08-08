package model

import "testing"

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{"valid pending", StatusPending, true},
		{"valid in_progress", StatusInProgress, true},
		{"valid done", StatusDone, true},
		{"invalid random string", TaskStatus("hello"), false},
		{"invalid uppercase", TaskStatus("PENDING"), false},
		{"invalid empty", TaskStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}