package domain

import "errors"

var (
	ErrProductNotFound              = errors.New("product not found")
	ErrSalePriceMustBeLessThanPrice = errors.New("sale_price must be less than price")
)
