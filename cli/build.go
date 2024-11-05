package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	"github.com/GuillaumeDesforges/ncbuild/builder"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build",
		Long:  "Build the project",
		Run: func(cmd *cobra.Command, args []string) {
			BuildCwd()
		},
	}
	return cmd
}

func BuildCwd() {
	ctx := context.Background()

	var recipe builder.Recipe
	recipeFilepath := "ncbuild.json"
	recipeFile, err := os.Open(recipeFilepath)
	if err != nil {
		logrus.Fatal(err)
	}
	defer recipeFile.Close()
	recipeFileBytes, err := os.ReadFile(recipeFilepath)
	err = json.Unmarshal(recipeFileBytes, &recipe)
	if err != nil {
		logrus.Fatal(err)
	}

	currentUser, err := user.Current()
	if err != nil {
		logrus.Fatal(err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Fatal(err)
	}
	defer dockerClient.Close()

	builder := builder.Builder{
		Store: &builder.Store{
			StorePath: "/tmp/ncbuild/",
			User:      currentUser.Gid,
		},
		DockerClient: dockerClient,
	}

	outputPath, err := builder.Build(ctx, recipe)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("%s", outputPath)
	fmt.Println()
}
