package select_sort

import (
	"reflect"
	"testing"
)

func TestSelectSort(t *testing.T) {
	type args struct {
		data []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: " [3,1,6,8,2,7,5]",
			args: args{data: []int{3,1,6,8,2,7,5}},
			want: []int{1,2,3,5,6,7,8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SelectSort(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
