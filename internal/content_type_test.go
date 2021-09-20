package internal

import "testing"

func Test_categorizeContentType(t *testing.T) {
	type args struct {
		headerStr string
	}
	tests := []struct {
		name string
		args args
		want PayloadType
	}{
		{
			name: "json",
			args: args{
				"application/json",
			},
			want: JsonPayload,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := categorizeContentType(tt.args.headerStr); got != tt.want {
				t.Errorf("categorizeContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
