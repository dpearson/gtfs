package gtfs

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_readCSVWithHeadings(t *testing.T) {
	type args struct {
		rc         io.Reader
		fields     map[string]bool
		strictMode bool
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]string
		wantErr bool
	}{
		{
			name: "Well-Formed (non-strict)",
			args: args{
				rc: strings.NewReader("agency_id,agency_name,agency_url,agency_timezone,agency_lang\n1,Test Agency,http://example.com,America/New_York,en"),
				fields: map[string]bool{
					"agency_id":       true,
					"agency_name":     true,
					"agency_url":      true,
					"agency_timezone": true,
					"agency_lang":     true,
				},
				strictMode: false,
			},
			want: []map[string]string{
				{
					"agency_id":       "1",
					"agency_name":     "Test Agency",
					"agency_url":      "http://example.com",
					"agency_timezone": "America/New_York",
					"agency_lang":     "en",
				},
			},
			wantErr: false,
		},
		{
			name: "Well-Formed (strict)",
			args: args{
				rc: strings.NewReader("agency_id,agency_name,agency_url,agency_timezone,agency_lang\n1,Test Agency,http://example.com,America/New_York,en"),
				fields: map[string]bool{
					"agency_id":       true,
					"agency_name":     true,
					"agency_url":      true,
					"agency_timezone": true,
					"agency_lang":     true,
				},
				strictMode: true,
			},
			want: []map[string]string{
				{
					"agency_id":       "1",
					"agency_name":     "Test Agency",
					"agency_url":      "http://example.com",
					"agency_timezone": "America/New_York",
					"agency_lang":     "en",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Field (non-strict)",
			args: args{
				rc: strings.NewReader("agency_id,agency_name,agency_url,agency_timezone\n1,Test Agency,http://example.com,America/New_York"),
				fields: map[string]bool{
					"agency_id":       true,
					"agency_name":     true,
					"agency_url":      true,
					"agency_timezone": true,
					"agency_lang":     true,
				},
				strictMode: false,
			},
			want: []map[string]string{
				{
					"agency_id":       "1",
					"agency_name":     "Test Agency",
					"agency_url":      "http://example.com",
					"agency_timezone": "America/New_York",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Field (strict)",
			args: args{
				rc: strings.NewReader("agency_id,agency_name,agency_url,agency_timezone\n1,Test Agency,http://example.com,America/New_York"),
				fields: map[string]bool{
					"agency_id":       true,
					"agency_name":     true,
					"agency_url":      true,
					"agency_timezone": true,
					"agency_lang":     true,
				},
				strictMode: true,
			},
			want: []map[string]string{
				{
					"agency_id":       "1",
					"agency_name":     "Test Agency",
					"agency_url":      "http://example.com",
					"agency_timezone": "America/New_York",
				},
			},
			wantErr: false,
		},
		{
			name: "Extra Field (non-strict)",
			args: args{
				rc: strings.NewReader("agency_id,agency_name,agency_url,agency_timezone,foo,agency_lang\n1,Test Agency,http://example.com,America/New_York,abc,en"),
				fields: map[string]bool{
					"agency_id":       true,
					"agency_name":     true,
					"agency_url":      true,
					"agency_timezone": true,
					"agency_lang":     true,
				},
				strictMode: false,
			},
			want: []map[string]string{
				{
					"agency_id":       "1",
					"agency_name":     "Test Agency",
					"agency_url":      "http://example.com",
					"agency_timezone": "America/New_York",
					"agency_lang":     "en",
				},
			},
			wantErr: false,
		},
		{
			name: "Extra Field (strict)",
			args: args{
				rc: strings.NewReader("agency_id,agency_name,agency_url,agency_timezone,foo,agency_lang\n1,Test Agency,http://example.com,America/New_York,abc,en"),
				fields: map[string]bool{
					"agency_id":       true,
					"agency_name":     true,
					"agency_url":      true,
					"agency_timezone": true,
					"agency_lang":     true,
				},
				strictMode: true,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCSVWithHeadings(tt.args.rc, tt.args.fields, tt.args.strictMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCSVWithHeadings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readCSVWithHeadings() = %v, want %v", got, tt.want)
			}
		})
	}
}
