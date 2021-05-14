package helper

import (
	"testing"
)

func TestFindOffset(t *testing.T) {
	tests := []struct {
		name         string
		fileText     string
		line, column int
		expected     int
	}{
		{
			name:     "Initial Character Offset",
			fileText: "first\nsecond",
			line:     1,
			column:   1,
			expected: 0,
		},
		{
			name:     "Middle of first line",
			fileText: "first\nsecond",
			line:     1,
			column:   3,
			expected: 2,
		},
		{
			name:     "Requested Column is Zero",
			fileText: "first\nsecond",
			line:     2,
			column:   0,
			expected: 6,
		},
		{
			name:     "At the end of the first line",
			fileText: "first\nsecond",
			line:     1,
			column:   6,
			expected: 5,
		},
		{
			name:     "Beyond the end of the first line",
			fileText: "first\nsecond",
			line:     1,
			column:   8,
			expected: -1,
		},
		{
			name:     "Middle of second line",
			fileText: "first\nsecond",
			line:     2,
			column:   3,
			expected: 8,
		},
		{
			name:     "Beyond the final line",
			fileText: "first\nsecond",
			line:     5,
			column:   3,
			expected: -1,
		},
		{
			name:     "Negative line and column",
			fileText: "first\nsecond",
			line:     -2,
			column:   -2,
			expected: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if res := FindOffset(test.fileText, test.line, test.column); res != test.expected {
				t.Errorf("Offset does not match expected. Got %d, Expected %d", res, test.expected)
			}
		})
	}
}

func BenchmarkFindOffset(b *testing.B) {
	// To compare benchmark, add it to this array
	benchmarks := []struct {
		name string
		f    func(string, int, int) int
	}{
		{
			name: "FindOffset",
			f:    FindOffset,
		},
	}

	var sampleInput string
	for i := 0; i < 100; i++ {
		sampleInput += `
resource "type" "name" {
    first  = "1"
    second = "2"
    third  = "3"
}
`
	}

	for _, bench := range benchmarks {
		bench.name = "FindOffset"
		b.Run(bench.name, func(b *testing.B) {
			b.Run("Null Inputs", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bench.f("", 0, 0)
				}
			})
			b.Run("Longer Inputs", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bench.f(sampleInput, 65, 8)
				}
			})
			b.Run("Early but beyond end of line", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bench.f(sampleInput, 5, 200)
				}
			})
		})
	}
}
