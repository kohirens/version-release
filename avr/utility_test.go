package main

import "testing"

func Test_validateMergeType(runner *testing.T) {
	cases := []struct {
		name    string
		mType   string
		want    string
		wantErr bool
	}{
		{"known", "squash", "squash", false},
		{"unknown", "sasquash", "", true},
	}
	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			got, err := validateMergeType(c.mType)
			if (err != nil) != c.wantErr {
				t.Errorf("validateMergeType() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("validateMergeType() got = %v, want %v", got, c.want)
			}
		})
	}
}
