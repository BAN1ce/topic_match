package utils

import (
	"reflect"
	"testing"
)

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

func TestSplitTopic(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		//{
		//	name: " /a/b",
		//	args: args{
		//		topic: "/a/b",
		//	},
		//	want: []string{
		//		"/",
		//		"a",
		//		"b",
		//	},
		//},
		//{
		//	name: " a/b",
		//	args: args{
		//		topic: "a/b",
		//	},
		//	want: []string{
		//		"a",
		//		"b",
		//	},
		//},
		//{
		//	name: " a",
		//	args: args{
		//		topic: "a",
		//	},
		//	want: []string{
		//		"a",
		//	},
		//},
		//{
		//	name: "/a",
		//	args: args{
		//		topic: "/a",
		//	},
		//	want: []string{
		//		"/",
		//		"a",
		//	},
		//},
		{
			name: "/",
			args: args{
				topic: "/",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitTopic(tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
