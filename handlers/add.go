package handlers

import (
	"net/http"

	"github.com/acd19ml/gofolder/data"
)

// swagger:route POST /products products addProduct
// Add a product to the list of products
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}