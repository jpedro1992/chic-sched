package util

import (
	"reflect"
	"testing"
)

func TestNewAllocation(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    *Allocation
		wantErr bool
	}{
		{name: "test1",
			args: args{
				size: 3,
			},
			want: &Allocation{
				x: []int{0, 0, 0},
			},
			wantErr: false,
		},
		{name: "test2",
			args: args{
				size: 0,
			},
			want: &Allocation{
				x: []int{},
			},
			wantErr: false,
		},
		{name: "test3",
			args: args{
				size: -2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAllocation(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAllocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAllocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAllocationCopy(t *testing.T) {
	type args struct {
		value []int
	}
	tests := []struct {
		name    string
		args    args
		want    *Allocation
		wantErr bool
	}{
		{name: "test1",
			args: args{
				value: []int{1, 2, 3},
			},
			want: &Allocation{
				x: []int{1, 2, 3},
			},
			wantErr: false,
		},
		{name: "test2",
			args: args{
				value: []int{},
			},
			want: &Allocation{
				x: []int{},
			},
			wantErr: false,
		},
		{name: "test3",
			args: args{
				value: nil,
			},
			want: &Allocation{
				x: []int{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAllocationCopy(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAllocationCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAllocationCopy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllocation_SetValue(t *testing.T) {
	type fields struct {
		x []int
	}
	type args struct {
		value []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Allocation
	}{
		{name: "test1",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				value: []int{4, 5},
			},
			want: &Allocation{
				x: []int{4, 5},
			},
		},
		{name: "test2",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				value: []int{},
			},
			want: &Allocation{
				x: []int{},
			},
		},
		{name: "test3",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				value: nil,
			},
			want: &Allocation{
				x: []int{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Allocation{
				x: tt.fields.x,
			}
			a.SetValue(tt.args.value)
			got := a
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Allocation_SetValue() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestAllocation_Fit(t *testing.T) {
	type fields struct {
		x []int
	}
	type args struct {
		allocated *Allocation
		capacity  *Allocation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "test1",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{1, 1, 0},
				},
				capacity: &Allocation{
					x: []int{5, 4, 3},
				},
			},
			want: true,
		},
		{name: "test2",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{1, 1, 0},
				},
				capacity: &Allocation{
					x: []int{2, 3, 3},
				},
			},
			want: true,
		},
		{name: "test3",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{1, 1, 0},
				},
				capacity: &Allocation{
					x: []int{5, 2, 3},
				},
			},
			want: false,
		},
		{name: "test4",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{0, 0},
				},
				capacity: &Allocation{
					x: []int{5, 5, 5},
				},
			},
			want: false,
		},
		{name: "test5",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{5, 5, 0},
				},
				capacity: &Allocation{
					x: []int{5, 5, 5},
				},
			},
			want: false,
		},
		{name: "test6",
			fields: fields{
				x: []int{1, 2, 0},
			},
			args: args{
				allocated: &Allocation{
					x: []int{1, 1, 0},
				},
				capacity: &Allocation{
					x: []int{5, 4, 3},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Allocation{
				x: tt.fields.x,
			}
			if got := a.Fit(tt.args.allocated, tt.args.capacity); got != tt.want {
				t.Errorf("Allocation.Fit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllocation_NumberToFit(t *testing.T) {
	type fields struct {
		x []int
	}
	type args struct {
		allocated *Allocation
		capacity  *Allocation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{name: "test1",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{0, 0, 0},
				},
				capacity: &Allocation{
					x: []int{5, 4, 3},
				},
			},
			want: 1,
		},
		{name: "test2",
			fields: fields{
				x: []int{1, 2, 4},
			},
			args: args{
				allocated: &Allocation{
					x: []int{0, 0, 0},
				},
				capacity: &Allocation{
					x: []int{5, 4, 3},
				},
			},
			want: 0,
		},
		{name: "test3",
			fields: fields{
				x: []int{0, 1, 0},
			},
			args: args{
				allocated: &Allocation{
					x: []int{0, 0, 0},
				},
				capacity: &Allocation{
					x: []int{5, 4, 3},
				},
			},
			want: 4,
		},
		{name: "test5",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{1, 1, 0},
				},
				capacity: &Allocation{
					x: []int{5, 4, 3},
				},
			},
			want: 1,
		},
		{name: "test6",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{1, 1, 2},
				},
				capacity: &Allocation{
					x: []int{10, 8, 6},
				},
			},
			want: 1,
		},
		{name: "test7",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				allocated: &Allocation{
					x: []int{1, 1, 2},
				},
				capacity: &Allocation{
					x: []int{10, 8, 9},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Allocation{
				x: tt.fields.x,
			}
			if got := a.NumberToFit(tt.args.allocated, tt.args.capacity); got != tt.want {
				t.Errorf("Allocation.NumberToFit() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestAllocation_StringPretty(t *testing.T) {
	type fields struct {
		x []int
	}
	type args struct {
		resourceNames []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				resourceNames: []string{"cpu", "memory", "gpu"},
			},
			want: "[cpu:1, memory:2, gpu:3]",
		},
		{
			name: "test2",
			fields: fields{
				x: []int{1, 2, 3},
			},
			args: args{
				resourceNames: []string{"cpu", "memory"},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Allocation{
				x: tt.fields.x,
			}
			if got := a.StringPretty(tt.args.resourceNames); got != tt.want {
				t.Errorf("Allocation.StringPretty() = %v, want %v", got, tt.want)
			}
		})
	}
}
