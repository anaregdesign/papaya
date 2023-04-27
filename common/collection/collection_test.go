package collection

import (
	"context"
	"github.com/anaregdesign/papaya/common/model"
	"reflect"
	"testing"
)

var ctx = context.Background()

func TestForEach(t *testing.T) {
	var a [32]int
	for i := 0; i < len(a); i++ {
		a[i] = i
	}

	t.Run("valid case", func(t *testing.T) {
		ForEach(ctx, a[:], func(i int) {
			if a[i] != i {
				t.Errorf("ForEach() = %v, want %v", a[i], i)
			}
		})
	})
}

func TestMap(t *testing.T) {
	var a [32]int
	for i := 0; i < len(a); i++ {
		a[i] = i
	}
	t.Run("valid case", func(t *testing.T) {
		b := Map(ctx, a[:], func(i int) int {
			return i * 2
		})

		for i := 0; i < len(b); i++ {
			if b[i] != a[i]*2 {
				t.Errorf("Map() = %v, want %v", b[i], a[i]*2)
			}
		}
	})
}

func TestReduce(t *testing.T) {
	add := func(a, b int) int { return a + b }
	type args[T any] struct {
		ctx      context.Context
		slice    []T
		operator model.Operator[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "valid case",
			args: args[int]{
				ctx:      ctx,
				slice:    []int{1, 2, 3, 4, 5},
				operator: add,
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.ctx, tt.args.slice, tt.args.operator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
