package domain

import (
	"errors"
	"regexp"
)

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var validCategoryName = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("category name cannot be empty")
	}

	if !validCategoryName.MatchString(c.Name) {
		return errors.New("category name can only contain letters, numbers, and spaces")
	}

	return nil
}
