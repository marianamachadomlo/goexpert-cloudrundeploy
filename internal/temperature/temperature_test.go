package temperature

import "testing"

func TestFromCelsius(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		wantC    float64
		wantF    float64
		wantK    float64
	}{
		{
			name:    "example from spec",
			celsius: 28.5,
			wantC:   28.5,
			wantF:   83.3,
			wantK:   301.5,
		},
		{
			name:    "zero celsius",
			celsius: 0,
			wantC:   0,
			wantF:   32,
			wantK:   273,
		},
		{
			name:    "negative celsius",
			celsius: -10,
			wantC:   -10,
			wantF:   14,
			wantK:   263,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromCelsius(tt.celsius)

			if got.TempC != tt.wantC {
				t.Fatalf("TempC = %v, want %v", got.TempC, tt.wantC)
			}
			if got.TempF != tt.wantF {
				t.Fatalf("TempF = %v, want %v", got.TempF, tt.wantF)
			}
			if got.TempK != tt.wantK {
				t.Fatalf("TempK = %v, want %v", got.TempK, tt.wantK)
			}
		})
	}
}

func TestFromCelsiusFormula(t *testing.T) {
	celsius := 25.0
	got := FromCelsius(celsius)

	expectedF := celsius*1.8 + 32
	expectedK := celsius + 273

	if got.TempF != round(expectedF, 1) {
		t.Fatalf("TempF formula failed: got %v, want %v", got.TempF, round(expectedF, 1))
	}
	if got.TempK != round(expectedK, 2) {
		t.Fatalf("TempK formula failed: got %v, want %v", got.TempK, round(expectedK, 2))
	}
}
