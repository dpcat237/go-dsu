package module

import (
	"testing"

	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

func TestDifferences_highestLevel(t *testing.T) {
	tests := []struct {
		name string
		dffs Differences
		want DiffLevel
	}{
		{
			name: "Check is selected DiffWeightMedium level",
			dffs: Differences{
				Difference{Level: DiffWeightLow},
				Difference{Level: DiffWeightMedium},
				Difference{Level: DiffWeightMedium},
				Difference{Level: DiffWeightLow},
			},
			want: DiffWeightMedium,
		},
		{
			name: "Check is selected DiffWeightCritical level",
			dffs: Differences{
				Difference{Level: DiffWeightCritical},
				Difference{Level: DiffWeightMedium},
				Difference{Level: DiffWeightMedium},
				Difference{Level: DiffWeightLow},
			},
			want: DiffWeightCritical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dffs.highestLevel(); got != tt.want {
				t.Errorf("highestLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifferences_vulnerabilityLevel(t *testing.T) {
	type args struct {
		vln vulnerability.Vulnerability
	}
	tests := []struct {
		name string
		dffs Differences
		args args
		want DiffLevel
	}{
		{
			name: "Check vulnerability level DiffWeightLow 1",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 0},
			},
			want: DiffWeightLow,
		},
		{
			name: "Check vulnerability level DiffWeightLow 2",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 3},
			},
			want: DiffWeightLow,
		},
		{
			name: "Check vulnerability level DiffWeightMedium 2",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 4},
			},
			want: DiffWeightMedium,
		},
		{
			name: "Check vulnerability level DiffWeightMedium 2",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 6},
			},
			want: DiffWeightMedium,
		},
		{
			name: "Check vulnerability level DiffWeightHigh 1",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 7},
			},
			want: DiffWeightHigh,
		},
		{
			name: "Check vulnerability level DiffWeightHigh 2",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 8},
			},
			want: DiffWeightHigh,
		},
		{
			name: "Check vulnerability level DiffWeightCritical 1",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 9},
			},
			want: DiffWeightCritical,
		},
		{
			name: "Check vulnerability level DiffWeightCritical 2",
			dffs: Differences{},
			args: args{
				vln: vulnerability.Vulnerability{CvssScore: 10},
			},
			want: DiffWeightCritical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dffs.vulnerabilityLevel(tt.args.vln); got != tt.want {
				t.Errorf("vulnerabilityLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
