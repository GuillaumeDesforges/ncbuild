package builder

import (
	"fmt"
	"path"
)

type IStore interface {
	GetOutputPath(recipe Recipe) string
	GetStorePath() string
}

type Store struct {
	StorePath string
	User      string
}

func (s *Store) GetOutputPath(recipe Recipe) string {
	recipeHash := recipe.Hash()
	outputStoreName := fmt.Sprintf("%s-%s", recipeHash, recipe.Name)
	return path.Join(s.StorePath, outputStoreName)
}

func (s *Store) GetStorePath() string {
	return s.StorePath
}
