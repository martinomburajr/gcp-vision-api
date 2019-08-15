package router

import "testing"

func TestExtractGCSPath(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    string
		wantErr bool
	}{
		{"err-nil", "", "", true},
		{"err-notGCSscheme", "sometext", "", true},
		{"err-onlyGCSScheme", "gs://", "", true},
		{"err-onlyGCSScheme", "gs://dddd", "", false},
		{"ok-path", "gs://dddd/somepath", "somepath", false},
		{"ok-deepPath", "gs://dddd/somepath/deeperpath", "somepath/deeperpath", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractGCSPath(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractGCSPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractGCSPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
