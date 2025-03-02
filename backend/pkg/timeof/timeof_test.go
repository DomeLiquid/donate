package timeof

import (
	"testing"
	"time"
)

func TestTimeOf(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      time.Time
		wantValid bool
	}{
		{
			name:      "Empty input",
			input:     "",
			want:      time.Time{},
			wantValid: false,
		},
		{
			name:      "Valid date input",
			input:     "2021-01-01",
			want:      time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local),
			wantValid: true,
		},
		{
			name:      "Valid datetime input",
			input:     "2021-01-01 12:00:00",
			want:      time.Time{},
			wantValid: false,
		},
		{
			name:      "Valid timestamp second input",
			input:     "1609459200",
			want:      time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			wantValid: true,
		},
		{
			name:      "Valid timestamp milisecond input",
			input:     "1707123089123",
			want:      time.Date(2024, time.February, 5, 8, 51, 29, 123000000, time.UTC),
			wantValid: true,
		},
		{
			name:      "Invalid input",
			input:     "invalid",
			want:      time.Time{},
			wantValid: false,
		},
		// 20060102
		{
			name:      "Valid date format 20060102 input",
			input:     "20240101",
			want:      time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local),
			wantValid: true,
		},
		// 2006-01-02
		{
			name:      "Valid date format 2006-01-02 input",
			input:     "2024-01-01",
			want:      time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local),
			wantValid: true,
		},
		// 20060102150405
		{
			name:      "Valid date format 20060102150405 input",
			input:     "20240102150405",
			want:      time.Date(2024, time.January, 2, 15, 4, 5, 0, time.Local),
			wantValid: true,
		},
		// 2006-01-02-15-04-05
		{
			name:      "Valid date format 2006-01-02-15-04-05 input",
			input:     "2024-01-02-15-04-05",
			want:      time.Date(2024, time.January, 2, 15, 4, 5, 0, time.Local),
			wantValid: true,
		},
		// 2006-01-02T15:04:05
		{
			name:      "Valid date format 2006-01-02T15:04:05 input",
			input:     "2024-01-02T15:04:05",
			want:      time.Date(2024, time.January, 2, 15, 4, 5, 0, time.Local),
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotValid := TimeOf(tt.input)

			if gotValid != tt.wantValid {
				t.Errorf("name %s ,TimeOf() gotValid = %v, want %v", tt.name, gotValid, tt.wantValid)
			}

			if gotValid == tt.wantValid && !got.UTC().Equal(tt.want.UTC()) {
				t.Errorf("name %s, TimeOf() got = %v, want %v", tt.name, got.UTC(), tt.want.UTC())
			}
		})
	}
}
