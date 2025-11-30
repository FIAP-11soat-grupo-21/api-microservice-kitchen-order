package value_objects

import "fmt"

type Slug struct {
	value string
}

func NewSlug(value string) (Slug, error) {
	if len(value) < 3 {
		return Slug{}, fmt.Errorf("name must be at least 3 characters long")
	}

	if len(value) > 100 {
		return Slug{}, fmt.Errorf("name must be at most 100 characters long")
	}

	return Slug{value: value}, nil
}

func (n *Slug) Value() string {
	return n.value
}
