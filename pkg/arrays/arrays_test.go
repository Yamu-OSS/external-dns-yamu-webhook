package arrays

import "testing"

func TestContains(t *testing.T) {
	type args struct {
		arr []string
		v   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "contains",
			args: args{
				arr: []string{"1", "2"},
				v:   "1",
			},
			want: true,
		},
		{
			name: "not contains",
			args: args{
				arr: []string{"1", "2"},
				v:   "3",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.arr, tt.args.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
