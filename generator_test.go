package go_verification

import (
	"regexp"
	"testing"
)

func TestNumberGeneratorGenerate(t *testing.T) {
	t.Run("Test default configuration", func(t *testing.T) {
		gen := NewNumberGenerator(5, false)
		result := gen.Generate()

		if len(result) != 5 {
			t.Errorf("Expected result length of 5, but got %d", len(result))
		}

		for _, char := range result {
			if char < '0' || char > '9' {
				t.Errorf("Expected only digits in the result, but got %s", result)
			}
		}
	})

	t.Run("Test non-zero start", func(t *testing.T) {
		gen := NewNumberGenerator(5, true)
		result := gen.Generate()

		if len(result) != 5 {
			t.Errorf("Expected result length of 5, but got %d", len(result))
		}

		if result[0] == '0' {
			t.Errorf("Expected result to not start with '0', but got %s", result)
		}

		for _, char := range result[1:] {
			if char < '0' || char > '9' {
				t.Errorf("Expected only digits in the result, but got %s", result)
			}
		}
	})
}

func TestAlphabetGeneratorGenerate(t *testing.T) {
	tests := []struct {
		name           string
		length         int
		allCapital     bool
		allNonCapital  bool
		expectedLength int
	}{
		{"Test Generate with lowercase only", 10, false, true, 10},
		{"Test Generate with uppercase only", 15, true, false, 15},
		{"Test Generate with mixed case", 20, false, false, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewAlphabetGenerator(tt.length, tt.allCapital, tt.allNonCapital)
			result := generator.Generate()

			if len(result) != tt.expectedLength {
				t.Errorf("Expected length %d, but got length %d", tt.expectedLength, len(result))
			}

			for _, char := range result {
				if tt.allCapital && !isCapital(byte(char)) {
					t.Errorf("Expected all capital letters, but found non-capital letter: %c", char)
				}

				if tt.allNonCapital && !isNonCapital(byte(char)) {
					t.Errorf("Expected all non-capital letters, but found capital letter: %c", char)
				}
			}
		})
	}
}

func isCapital(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isNonCapital(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func TestWordGeneratorGenerate(t *testing.T) {
	t.Run("Test default configuration", func(t *testing.T) {
		gen := NewWordGenerator(5)
		result := gen.Generate()

		if len(result) != 5 {
			t.Errorf("Expected result length of 5, but got %d", len(result))
		}
	})
}

func TestRegexGeneratorGenerate(t *testing.T) {
	t.Run("Test default configuration", func(t *testing.T) {
		expectedPattern := `^G-\d{1,2}\d+\w+\d{1}(this|that)?[12]{2}$`
		gen := NewRegexGenerator(expectedPattern)
		result := gen.Generate()

		matched, err := regexp.MatchString(expectedPattern, result)
		if err != nil {
			t.Errorf("Error matching pattern: %v", err)
		}

		if !matched {
			t.Errorf("Generated string does not match the expected pattern: %s", result)
		}
	})
}
