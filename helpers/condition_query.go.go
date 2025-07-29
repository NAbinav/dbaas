// helpers/condition_extract.go
package helpers

import (
	"fmt"
	"strings"
)

var operatorMap = map[string]string{
	"eq":   "=",
	"ne":   "<>",
	"neq":  "<>",
	"gt":   ">",
	"lt":   "<",
	"gte":  ">=",
	"lte":  "<=",
	"like": "LIKE",
	"in":   "IN",
	"nin":  "NOT IN",
}

func Condition_extract(query map[string][]string) (string, error) {
	var conditions []string

	for key, values := range query {
		if key == "or" {
			// Handle OR conditions
			orConditions := strings.Split(values[0], ",")
			var orParts []string

			for _, orCond := range orConditions {
				parts := strings.Split(orCond, ":")
				if len(parts) != 2 {
					return "", fmt.Errorf("invalid OR condition format: %s - use format: field_operator:value", orCond)
				}

				fieldOp := parts[0]
				value := parts[1]

				fieldParts := strings.SplitN(fieldOp, "_", 2)
				if len(fieldParts) != 2 {
					return "", fmt.Errorf("invalid OR field format: %s - use format: field_operator", fieldOp)
				}

				column := fieldParts[0]
				opKey := fieldParts[1]
				op, ok := operatorMap[opKey]
				if !ok {
					return "", fmt.Errorf("unsupported operator in OR: %s", opKey)
				}

				if op == "IN" || op == "NOT IN" {
					items := strings.Split(value, "|")
					var quotedItems []string
					for _, item := range items {
						quotedItems = append(quotedItems, fmt.Sprintf("'%s'", strings.TrimSpace(item)))
					}
					orParts = append(orParts, fmt.Sprintf("%s %s (%s)", column, op, strings.Join(quotedItems, ", ")))
				} else {
					orParts = append(orParts, fmt.Sprintf("%s %s '%s'", column, op, value))
				}
			}

			if len(orParts) > 0 {
				conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(orParts, " OR ")))
			}
		} else {
			// Handle regular AND conditions
			parts := strings.SplitN(key, "_", 2)
			if len(parts) != 2 {
				return "", fmt.Errorf("invalid filter key: %s - use format: field_operator", key)
			}

			column := parts[0]
			opKey := parts[1]
			op, ok := operatorMap[opKey]
			if !ok {
				return "", fmt.Errorf("unsupported operator: %s", opKey)
			}

			raw := values[0]
			if op == "IN" || op == "NOT IN" {
				items := strings.Split(raw, ",")
				var quotedItems []string
				for _, item := range items {
					quotedItems = append(quotedItems, fmt.Sprintf("'%s'", strings.TrimSpace(item)))
				}
				conditions = append(conditions, fmt.Sprintf("%s %s (%s)", column, op, strings.Join(quotedItems, ", ")))
			} else {
				conditions = append(conditions, fmt.Sprintf("%s %s '%s'", column, op, raw))
			}
		}
	}

	if len(conditions) == 0 {
		return "", nil
	}
	return "WHERE " + strings.Join(conditions, " AND "), nil
}
