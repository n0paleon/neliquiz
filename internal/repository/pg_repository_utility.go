package repository

import "strings"

func sanitizeSort(sortBy, order string) (string, string) {
	allowedSortBy := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"hit":        true,
	}
	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}

	// whitelist order
	order = sanitizeOrder(order)

	return sortBy, order
}

func sanitizeOrder(order string) string {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}
	return order
}
