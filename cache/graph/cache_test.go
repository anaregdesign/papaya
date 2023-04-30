package graph

import (
	"context"
	"github.com/anaregdesign/papaya/cache"
	"reflect"
	"testing"
	"time"
)

func TestGraph_AddEdge(t *testing.T) {
	e := newEdgeCache[string](context.Background(), time.Minute)
	type args[S comparable] struct {
		tail S
		head S
		w    float64
	}
	type testCase[S comparable, T any] struct {
		name string
		g    GraphCache[S, T]
		args args[S]
	}
	tests := []testCase[string, string]{
		{
			name: "TestGraph_AddEdge",
			g: GraphCache[string, string]{
				edges: e,
			},
			args: args[string]{
				tail: "tail",
				head: "head",
				w:    1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.AddEdge(tt.args.tail, tt.args.head, tt.args.w)
		})
	}
}

func TestGraph_AddEdgeWithTTL(t *testing.T) {
	e := newEdgeCache[string](context.Background(), time.Minute)
	type args[S comparable] struct {
		tail S
		head S
		w    float64
		ttl  time.Duration
	}
	type testCase[S comparable, T any] struct {
		name string
		g    GraphCache[S, T]
		args args[S]
	}
	tests := []testCase[string, string]{
		{
			name: "TestGraph_AddEdgeWithTTL",
			g: GraphCache[string, string]{
				edges: e,
			},
			args: args[string]{
				tail: "tail",
				head: "head",
				w:    1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.AddEdgeWithTTL(tt.args.tail, tt.args.head, tt.args.w, tt.args.ttl)
		})
	}
}

func TestGraph_AddVertex(t *testing.T) {
	type args[S comparable, T any] struct {
		key   S
		value T
	}
	type testCase[S comparable, T any] struct {
		name string
		g    GraphCache[S, T]
		args args[S, T]
	}
	tests := []testCase[string, string]{
		{
			name: "TestGraph_AddVertex",
			g: GraphCache[string, string]{
				vertices: cache.NewCache[string, string](context.Background(), time.Minute),
			},
			args: args[string, string]{key: "key", value: "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.AddVertex(tt.args.key, tt.args.value)
		})
	}
}

func TestGraph_AddVertexWithTTL(t *testing.T) {
	type args[S comparable, T any] struct {
		key   S
		value T
		ttl   time.Duration
	}
	type testCase[S comparable, T any] struct {
		name string
		g    GraphCache[S, T]
		args args[S, T]
	}
	tests := []testCase[string, string]{
		{
			name: "TestGraph_AddVertexWithTTL",
			g: GraphCache[string, string]{
				vertices: cache.NewCache[string, string](context.Background(), time.Minute),
			},
			args: args[string, string]{key: "key", value: "value", ttl: time.Minute},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.AddVertexWithTTL(tt.args.key, tt.args.value, tt.args.ttl)
		})
	}
}

func TestGraph_GetVertex(t *testing.T) {
	v := cache.NewCache[string, string](context.Background(), time.Minute)
	v.Set("key", "value")

	type args[S comparable] struct {
		key S
	}
	type testCase[S comparable, T any] struct {
		name  string
		g     GraphCache[S, T]
		args  args[S]
		want  T
		want1 bool
	}
	tests := []testCase[string, string]{
		{
			name: "TestGraph_GetVertex",
			g: GraphCache[string, string]{
				vertices: v,
			},
			args:  args[string]{key: "key"},
			want:  "value",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.g.GetVertex(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVertex() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetVertex() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGraph_getWeight(t *testing.T) {
	e := newEdgeCache[string](context.Background(), time.Minute)

	type args[S comparable] struct {
		tail S
		head S
	}
	type testCase[S comparable, T any] struct {
		name string
		g    GraphCache[S, T]
		args args[S]
		want float64
	}
	tests := []testCase[string, string]{
		{
			name: "TestGraph_getWeight",
			g: GraphCache[string, string]{
				edges: e,
			},
			args: args[string]{tail: "tail", head: "head"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.getWeight(tt.args.tail, tt.args.head); got != tt.want {
				t.Errorf("getWeight() = %v, want %v", got, tt.want)
			}
		})
	}
}
