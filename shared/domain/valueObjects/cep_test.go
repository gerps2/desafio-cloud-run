package valueObjects

import (
	"testing"
)

func TestNewCep(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{
			name:        "Valid CEP with dash",
			input:       "12345-678",
			expectError: false,
			expected:    "12345-678",
		},
		{
			name:        "Valid CEP without dash",
			input:       "12345678",
			expectError: false,
			expected:    "12345-678",
		},
		{
			name:        "Invalid CEP - too short",
			input:       "1234-567",
			expectError: true,
		},
		{
			name:        "Invalid CEP - too long",
			input:       "123456-789",
			expectError: true,
		},
		{
			name:        "Invalid CEP - contains letters",
			input:       "1234a-678",
			expectError: true,
		},
		{
			name:        "Invalid CEP - empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "Invalid CEP - only spaces",
			input:       "     ",
			expectError: true,
		},
		{
			name:        "Invalid CEP - special characters",
			input:       "12345@678",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cep, err := NewCep(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				return
			}

			if cep.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, cep.String())
			}
		})
	}
}

func TestCepString(t *testing.T) {
	cep, err := NewCep("12345678")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "12345-678"
	if cep.String() != expected {
		t.Errorf("Expected %s, got %s", expected, cep.String())
	}
}

func TestCepValidation(t *testing.T) {
	// Test boundary conditions
	validCeps := []string{
		"00000-000",
		"99999-999",
		"12345-678",
	}

	for _, validCep := range validCeps {
		t.Run("Valid_"+validCep, func(t *testing.T) {
			_, err := NewCep(validCep)
			if err != nil {
				t.Errorf("Expected valid CEP %s to not return error, got: %v", validCep, err)
			}
		})
	}

	invalidCeps := []string{
		"1234-567",   // too short
		"123456-789", // too long
		"abcde-fgh",  // letters
		"12345-67a",  // letter at end
		"1234567890", // too many digits
	}

	for _, invalidCep := range invalidCeps {
		t.Run("Invalid_"+invalidCep, func(t *testing.T) {
			_, err := NewCep(invalidCep)
			if err == nil {
				t.Errorf("Expected invalid CEP %s to return error", invalidCep)
			}
		})
	}
}
