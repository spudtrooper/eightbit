package convert

import (
	"testing"
)

func TestMakeOutput(t *testing.T) {
	var tests = []struct {
		name                        string
		input, outputDir, converter string
		want                        string
	}{
		{
			name:      "no path",
			input:     "foo.png",
			outputDir: "",
			converter: "pixelate",
			want:      "foo-pixelate.png",
		},
		{
			name:      "one dir",
			input:     "path/foo.png",
			outputDir: "",
			converter: "pixelate",
			want:      "path/foo-pixelate.png",
		},
		{
			name:      "multi dirs",
			input:     "path/to/the/foo.png",
			outputDir: "",
			converter: "pixelate",
			want:      "path/to/the/foo-pixelate.png",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if want, got := test.want, makeOutput(test.input, test.outputDir, test.converter); want != got {
				t.Errorf("makeOutput(%q,%q): want %v, got %v", test.input, test.converter, test.want, got)
			}
		})
	}
}
