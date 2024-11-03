package builder

import (
	"fmt"

	"github.com/gohugoio/hashstructure"
)

type Recipe struct {
	Name       string   `json:""`
	Executable string   `json:""`
	Args       []string `json:""`
}

func (r Recipe) Hash() string {
	hash, err := hashstructure.Hash(r, nil)
	if err != nil {
		panic(err)
	}

	strHash := fmt.Sprintf("%d", hash)
	return strHash
}
