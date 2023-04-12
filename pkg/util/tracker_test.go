package util

import "testing"

func TestCreateTrackingId(t *testing.T) {
	type args struct {
		salt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "sd",
			args: args{
				salt: "ssss",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateTrackingId(tt.args.salt); got != tt.want {
				t.Errorf("CreateTrackingId() = %v, want %v", got, tt.want)
			}
		})
	}
}
