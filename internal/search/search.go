package search

import (
	"github.com/EbumbaE/CountVibe/internal/entities"
)

type search interface {
}

type ResultSearch struct {
	Products []entities.Product
}

func GetSearch() ResultSearch {
	return ResultSearch{[]entities.Product{}}
}
