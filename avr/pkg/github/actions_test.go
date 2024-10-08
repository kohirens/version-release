package github

import (
	"os"
	"testing"
)

func TestAddOutputVar(t *testing.T) {
	outputFixture := tmpDir + "/test1.env"
	os.Setenv(gaOutputFile, outputFixture)
	cases := []struct {
		name     string
		varName  string
		workflow string
		wantErr  bool
		want     string
	}{
		{"success", "workflow", "publish", false, "workflow=publish\n"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := AddOutputVar(c.varName, c.workflow); (err != nil) != c.wantErr {
				t.Errorf("AddOutputVar() error = %v, wantErr %v", err, c.wantErr)
			}

			b, e1 := os.ReadFile(outputFixture)
			if e1 != nil {
				t.Errorf("AddOutputVar() error reading %v", outputFixture)
			}

			if string(b) != c.want {
				t.Errorf("AddOutputVar() error %v, want %v", string(b), c.want)
			}
		})
	}
}
