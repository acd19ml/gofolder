package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "nico",
		Price: 1.00,
		SKU:   "abc-abc-abc",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
