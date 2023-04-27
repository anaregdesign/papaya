package collection

import (
	"context"
	"testing"
)

func TestForEach(t *testing.T) {
	ctx := context.Background()
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
	ctx := context.Background()
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
