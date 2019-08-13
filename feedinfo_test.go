package gtfs

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const testFeedInfoCSVValid = `feed_publisher_name,feed_publisher_url,feed_lang,feed_start_date,feed_end_date,feed_version
Test Publisher,http://feedinfo.example.com,en,,,3`

const testFeedInfoCSVInvalidMultipleRows = `feed_publisher_name,feed_publisher_url,feed_lang,feed_start_date,feed_end_date,feed_version
Test Publisher,http://feedinfo.example.com,en,,,3
Test Publisher,http://feedinfo.example.com,en,,,3`

func TestGTFS_processFeedInfo(t *testing.T) {
	expectedFeedInfo := FeedInfo{
		PublisherName: "Test Publisher",
		PublisherURL:  "http://feedinfo.example.com",
		Lang:          "en",
		StartDate:     "",
		EndDate:       "",
		Version:       "3",
	}
	type fields struct {
		FeedInfo FeedInfo
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantFeedInfo FeedInfo
	}{
		{
			name:   "Valid",
			fields: fields{},
			args: args{
				r: strings.NewReader(testFeedInfoCSVValid),
			},
			wantErr:      false,
			wantFeedInfo: expectedFeedInfo,
		},
		{
			name: "Valid (Existing Feed Info)",
			fields: fields{
				FeedInfo: FeedInfo{
					PublisherName: "Test Publisher Old",
					PublisherURL:  "http://feedinfo-old.example.com",
					Lang:          "de",
					StartDate:     "",
					EndDate:       "",
					Version:       "2",
				},
			},
			args: args{
				r: strings.NewReader(testFeedInfoCSVValid),
			},
			wantErr:      false,
			wantFeedInfo: expectedFeedInfo,
		},
		{
			name:   "Invalid (Multiple Rows)",
			fields: fields{},
			args: args{
				r: strings.NewReader(testFeedInfoCSVInvalidMultipleRows),
			},
			wantErr:      true,
			wantFeedInfo: FeedInfo{},
		},
		{
			name:   "Empty",
			fields: fields{},
			args: args{
				r: strings.NewReader(""),
			},
			wantErr:      true,
			wantFeedInfo: FeedInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				FeedInfo: tt.fields.FeedInfo,
			}
			if err := g.processFeedInfo(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("GTFS.processFeedInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g.FeedInfo, tt.wantFeedInfo) {
				t.Errorf("GTFS.processAgencies() FeedInfo = %v, wantFeedInfo %v", g.FeedInfo, tt.wantFeedInfo)
			}
		})
	}
}
