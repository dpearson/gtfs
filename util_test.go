package gtfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type mockOpener struct {
	rc          io.ReadCloser
	shouldError bool
}

func (m *mockOpener) Open() (io.ReadCloser, error) {
	if m.shouldError {
		return m.rc, fmt.Errorf("mock error opening resource")
	}

	return m.rc, nil
}

func Test_callWithOpenedReader(t *testing.T) {
	successFn := func(r io.Reader) error {
		return nil
	}
	errFn := func(r io.Reader) error {
		return fmt.Errorf("mock error in called function")
	}
	mockErroringOpener := &mockOpener{
		shouldError: true,
	}
	mockSuccessOpener := &mockOpener{
		rc:          ioutil.NopCloser(strings.NewReader("")),
		shouldError: false,
	}

	err := callWithOpenedReader(successFn, mockErroringOpener)
	if err == nil {
		t.Errorf("Expected error with mock erroring opener, but got none")
	}
	err = callWithOpenedReader(errFn, mockSuccessOpener)
	if err == nil {
		t.Errorf("Expected error with mock erroring function, but got none")
	}
	err = callWithOpenedReader(successFn, mockSuccessOpener)
	if err != nil {
		t.Errorf("Expected no error with mock success function, but got %v", err)
	}
}

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
