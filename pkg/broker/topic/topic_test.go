package topic

import (
	"reflect"
	"testing"
)

func TestNewShareTopic(t *testing.T) {
	var (
		shareTopic = &ShareTopic{
			FullTopic:  "$share/my-group/sensors/#",
			ShareGroup: "my-group",
			Topic:      "sensors/#",
		}
	)
	type args struct {
		fullTopic string
	}
	tests := []struct {
		name string
		args args
		want *ShareTopic
	}{
		{
			name: "",
			args: args{
				fullTopic: "$share/my-group/sensors/#",
			},
			want: shareTopic,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewShareTopic(tt.args.fullTopic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewShareTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
