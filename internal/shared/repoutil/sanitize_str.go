package repoutil

import "strings"

func SanitizeSort(sortBy, order string) (string, string) {
	allowedSortBy := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"hit":        true,
	}
	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}

	// whitelist order
	order = SanitizeOrder(order)

	return sortBy, order
}

func SanitizeOrder(order string) string {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}
	return order
}
