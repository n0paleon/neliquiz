package entities

import (
	"errors"
	"regexp"
)

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var validCategoryName = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)

func NewCategory(name string) (*Category, error) {
	c := &Category{
		Name: name,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	maxLen := 50 // maksimum panjang karakter nama
	if len(name) > maxLen {
		return nil, errors.New("category name too long")
	}

	return c, nil
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("category name cannot be empty")
	}

	if !validCategoryName.MatchString(c.Name) {
		return errors.New("category name can only contain letters, numbers, and spaces")
	}

	return nil
}
