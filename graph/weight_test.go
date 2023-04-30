package graph

import (
	"reflect"
	"testing"
	"time"
)

func Test_newWeight(t *testing.T) {
	tests := []struct {
		name string
		want *weight
	}{
		{
			name: "newWeight",
			want: &weight{
				values: make([]weightValue, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newWeight(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_weightValue_expired(t *testing.T) {
	type fields struct {
		value float64
		ttl   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "weightValue_expired",
			fields: fields{
				value: 1,
				ttl:   time.Now().Add(-time.Minute),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := weightValue{
				value: tt.fields.value,
				ttl:   tt.fields.ttl,
			}
			if got := w.expired(); got != tt.want {
				t.Errorf("expired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_weight_add(t *testing.T) {
	type fields struct {
		values []weightValue
	}
	type args struct {
		value float64
		ttl   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "weight_add",
			fields: fields{
				values: make([]weightValue, 0),
			},
			args: args{
				value: 1,
				ttl:   time.Minute,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &weight{
				values: tt.fields.values,
			}
			w.add(tt.args.value, tt.args.ttl)
		})
	}
}

func Test_weight_value(t *testing.T) {
	type fields struct {
		values []weightValue
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "weight_Value",
			fields: fields{
				values: []weightValue{
					{
						value: 1,
						ttl:   time.Now().Add(time.Minute),
					},
					{
						value: 1,
						ttl:   time.Now().Add(-time.Minute),
					},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &weight{
				values: tt.fields.values,
			}
			if got := w.value(); got != tt.want {
				t.Errorf("value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_weight_isZero(t *testing.T) {
	type fields struct {
		values []weightValue
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "weight_isZero",
			fields: fields{
				values: []weightValue{
					{
						value: 1,
						ttl:   time.Now().Add(time.Minute),
					},
				},
			},
			want: false,
		},
		{
			name: "weight_isZero",
			fields: fields{
				values: []weightValue{
					{
						value: 1,
						ttl:   time.Now().Add(-time.Minute),
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &weight{
				values: tt.fields.values,
			}
			if got := w.isZero(); got != tt.want {
				t.Errorf("isZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
