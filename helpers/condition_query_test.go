package helpers

import (
	"strings"
	"testing"
)

func checkConditions(t *testing.T, result string, expectedConditions []string, shouldStartWithWhere bool) {
	if shouldStartWithWhere && !strings.HasPrefix(result, "WHERE ") {
		t.Errorf("Result should start with 'WHERE ', got: %q", result)
	}

	for _, condition := range expectedConditions {
		if !strings.Contains(result, condition) {
			t.Errorf("Result should contain %q, got: %q", condition, result)
		}
	}

	if len(expectedConditions) > 1 {
		expectedAndCount := len(expectedConditions) - 1
		actualAndCount := strings.Count(result, " AND ")
		if actualAndCount != expectedAndCount {
			t.Errorf("Expected %d AND operators, got %d in: %q", expectedAndCount, actualAndCount, result)
		}
	}
}

func TestCondition_extract(t *testing.T) {
	tests := []struct {
		name               string
		query              map[string][]string
		expectedConditions []string
		shouldHaveWhere    bool
		wantErr            bool
	}{
		{
			name:               "empty query",
			query:              map[string][]string{},
			expectedConditions: []string{},
			shouldHaveWhere:    false,
			wantErr:            false,
		},
		{
			name: "simple equality",
			query: map[string][]string{
				"name_eq": {"john"},
			},
			expectedConditions: []string{"name = 'john'"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "multiple AND conditions",
			query: map[string][]string{
				"name_eq": {"john"},
				"age_gt":  {"25"},
			},
			expectedConditions: []string{"name = 'john'", "age > '25'"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "IN operator",
			query: map[string][]string{
				"status_in": {"active,pending,inactive"},
			},
			expectedConditions: []string{"status IN ('active', 'pending', 'inactive')"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "NOT IN operator",
			query: map[string][]string{
				"role_nin": {"admin,moderator"},
			},
			expectedConditions: []string{"role NOT IN ('admin', 'moderator')"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "LIKE operator",
			query: map[string][]string{
				"email_like": {"%@gmail.com"},
			},
			expectedConditions: []string{"email LIKE '%@gmail.com'"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "OR conditions",
			query: map[string][]string{
				"or": {"name_eq:john,age_gt:25"},
			},
			expectedConditions: []string{"(name = 'john' OR age > '25')"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "OR with IN operator",
			query: map[string][]string{
				"or": {"status_in:pending|processing,priority_eq:high"},
			},
			expectedConditions: []string{"(status IN ('pending', 'processing') OR priority = 'high')"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "mixed AND and OR",
			query: map[string][]string{
				"category_eq": {"electronics"},
				"or":          {"brand_eq:apple,price_lt:500"},
			},
			expectedConditions: []string{
				"category = 'electronics'",
				"(brand = 'apple' OR price < '500')",
			},
			shouldHaveWhere: true,
			wantErr:         false,
		},
		{
			name: "all comparison operators",
			query: map[string][]string{
				"name_eq":    {"john"},
				"age_gt":     {"25"},
				"salary_gte": {"50000"},
				"score_lt":   {"100"},
				"rating_lte": {"4.5"},
				"status_ne":  {"inactive"},
			},
			expectedConditions: []string{
				"name = 'john'",
				"age > '25'",
				"salary >= '50000'",
				"score < '100'",
				"rating <= '4.5'",
				"status <> 'inactive'",
			},
			shouldHaveWhere: true,
			wantErr:         false,
		},
		{
			name: "neq operator",
			query: map[string][]string{
				"status_neq": {"deleted"},
			},
			expectedConditions: []string{"status <> 'deleted'"},
			shouldHaveWhere:    true,
			wantErr:            false,
		},
		{
			name: "invalid operator",
			query: map[string][]string{
				"name_invalid": {"john"},
			},
			expectedConditions: []string{},
			shouldHaveWhere:    false,
			wantErr:            true,
		},
		{
			name: "invalid format - no underscore",
			query: map[string][]string{
				"invalidformat": {"john"},
			},
			expectedConditions: []string{},
			shouldHaveWhere:    false,
			wantErr:            true,
		},
		{
			name: "invalid OR format",
			query: map[string][]string{
				"or": {"invalid_format"},
			},
			expectedConditions: []string{},
			shouldHaveWhere:    false,
			wantErr:            true,
		},
		{
			name: "invalid OR field format",
			query: map[string][]string{
				"or": {"invalidfield:value"},
			},
			expectedConditions: []string{},
			shouldHaveWhere:    false,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Condition_extract(tt.query)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition_extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result if no error expected
			if !tt.wantErr {
				if len(tt.expectedConditions) == 0 {
					// Empty result expected
					if result != "" {
						t.Errorf("Expected empty result, got: %q", result)
					}
				} else {
					checkConditions(t, result, tt.expectedConditions, tt.shouldHaveWhere)
				}
			}
		})
	}
}

// Test with realistic URL query parameters
func TestCondition_extract_URLParams(t *testing.T) {
	// Simulate what you'd get from parsing URL: /users?name_eq=john&age_gt=25&status_in=active,pending
	query := map[string][]string{
		"name_eq":   {"john"},
		"age_gt":    {"25"},
		"status_in": {"active,pending"},
	}

	result, err := Condition_extract(query)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedConditions := []string{
		"name = 'john'",
		"age > '25'",
		"status IN ('active', 'pending')",
	}

	checkConditions(t, result, expectedConditions, true)
}

// Test edge cases
func TestCondition_extract_EdgeCases(t *testing.T) {
	t.Run("empty string values", func(t *testing.T) {
		query := map[string][]string{
			"name_eq": {""},
		}

		result, err := Condition_extract(query)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "WHERE name = ''"
		if result != expected {
			t.Errorf("Got %q, want %q", result, expected)
		}
	})

	t.Run("values with spaces", func(t *testing.T) {
		query := map[string][]string{
			"name_eq": {"John Doe"},
		}

		result, err := Condition_extract(query)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "WHERE name = 'John Doe'"
		if result != expected {
			t.Errorf("Got %q, want %q", result, expected)
		}
	})

	t.Run("IN with spaces in values", func(t *testing.T) {
		query := map[string][]string{
			"status_in": {"active user, pending approval, inactive"},
		}

		result, err := Condition_extract(query)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expectedConditions := []string{"status IN ('active user', 'pending approval', 'inactive')"}
		checkConditions(t, result, expectedConditions, true)
	})
}
