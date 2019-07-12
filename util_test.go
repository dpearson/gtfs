package gtfs

import "testing"

func Test_parseBool(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "true",
			args: args{
				val: "1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "false",
			args: args{
				val: "0",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "invalid (2)",
			args: args{
				val: "2",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid (foo)",
			args: args{
				val: "foo",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBool(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
