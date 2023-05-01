package pq

import (
	"reflect"
	"testing"
)

func TestFilterMap(t *testing.T) {
	type args[S comparable, T Number] struct {
		m   map[S]T
		top int
	}
	type testCase[S comparable, T Number] struct {
		name string
		args args[S, T]
		want map[S]T
	}
	tests := []testCase[string, float64]{
		{
			name: "TestFilterMap",
			args: args[string, float64]{
				m: map[string]float64{
					"one":   1,
					"two":   2,
					"three": 3,
					"four":  4,
					"five":  5,
				},
				top: 3,
			},
			want: map[string]float64{
				"five":  5,
				"four":  4,
				"three": 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterMap(tt.args.m, tt.args.top); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
