package utils

import "testing"

func TestParseShareTopic(t *testing.T) {
	type args struct {
		shareTopic string
	}
	tests := []struct {
		name           string
		args           args
		wantShareGroup string
		wantSubTopic   string
	}{
		{
			name: "parse",
			args: args{
				shareTopic: "$share/group/topic",
			},
			wantShareGroup: "group",
			wantSubTopic:   "topic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShareGroup, gotSubTopic := ParseShareTopic(tt.args.shareTopic)
			if gotShareGroup != tt.wantShareGroup {
				t.Errorf("ParseShareTopic() gotShareGroup = %v, want %v", gotShareGroup, tt.wantShareGroup)
			}
			if gotSubTopic != tt.wantSubTopic {
				t.Errorf("ParseShareTopic() gotSubTopic = %v, want %v", gotSubTopic, tt.wantSubTopic)
			}
		})
	}
}
