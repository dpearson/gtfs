package gtfs

import "testing"

func Test_parsePaymentMethod(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    PaymentMethod
		wantErr bool
	}{
		{
			name:    "Empty",
			val:     "",
			want:    PaymentMethodOnBoard,
			wantErr: true,
		},
		{
			name:    "On Board",
			val:     "0",
			want:    PaymentMethodOnBoard,
			wantErr: false,
		},
		{
			name:    "Before Boarding",
			val:     "1",
			want:    PaymentMethodBeforeBoarding,
			wantErr: false,
		},
		{
			name:    "Invalid",
			val:     "2",
			want:    PaymentMethodOnBoard,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePaymentMethod(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePaymentMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePaymentMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
