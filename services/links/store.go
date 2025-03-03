package links

import (
	"github.com/OmprakashD20/refero-api/repository"
)

type Store struct {
	db *repository.Queries
}

func NewStore(db *repository.Queries) *Store {
	return &Store{db}
}
