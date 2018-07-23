package gtfs

import "testing"

func Test_parseTransferType(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    TransferType
		wantErr bool
	}{
		{
			name:    "Empty",
			val:     "",
			want:    TransferTypeRecommended,
			wantErr: false,
		},
		{
			name:    "Recommended",
			val:     "0",
			want:    TransferTypeRecommended,
			wantErr: false,
		},
		{
			name:    "Timed Transfer Point",
			val:     "1",
			want:    TransferTypeTimed,
			wantErr: false,
		},
		{
			name:    "Minimum Required Time",
			val:     "2",
			want:    TransferTypeMinimumTime,
			wantErr: false,
		},
		{
			name:    "Not Possible",
			val:     "3",
			want:    TransferTypeNone,
			wantErr: false,
		},
		{
			name:    "Invalid",
			val:     "4",
			want:    TransferTypeRecommended,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTransferType(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTransferType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseTransferType() = %v, want %v", got, tt.want)
			}
		})
	}
}
