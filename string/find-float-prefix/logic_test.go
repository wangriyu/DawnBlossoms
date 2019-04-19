package logic

import (
	"testing"
)

func TestFindFloatPrefix(t *testing.T) {
	tests := []struct {
		input      string
		wantResult string
	}{
		{
			"1.1a",
			"1.1",
		}, {
			"abc",
			"",
		}, {
			"-1.1e3.3",
			"-1.1e3",
		}, {
			"-1.1e.",
			"-1.1",
		}, {
			"1e1",
			"1e1",
		}, {
			"55e",
			"55",
		}, {
			"1",
			"1",
		}, {
			"+1..1",
			"+1",
		}, {
			"-1.3e.10",
			"-1.3",
		}, {
			"-1.3e-2.5",
			"-1.3e-2",
		},
	}
	for _, tt := range tests {
		if gotResult := string(FindFloatPrefix(tt.input)); gotResult != tt.wantResult {
			t.Errorf("FindFloatPrefix() = %v, want %v", gotResult, tt.wantResult)
		}
	}
}
