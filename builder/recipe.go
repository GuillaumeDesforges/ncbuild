package builder

import (
	"fmt"

	"github.com/gohugoio/hashstructure"
)

type Recipe struct {
	// The name of the recipe
	Name string `json:"name"`
	// The Docker image to use to build the output within
	BuildDockerImage string `json:"build_docker_image"`
	// Inputs
	Inputs []string `json:"inputs"`
	// The executable to run within the Docker container
	Executable string `json:"executable"`
	// The arguments to pass to the executable
	Args []string `json:"args"`
}

func (r Recipe) Hash() string {
	hash, err := hashstructure.Hash(r, nil)
	if err != nil {
		panic(err)
	}

	strHash := fmt.Sprintf("%d", hash)
	return strHash
}
