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
			if debug, _ := cmd.Flags().GetBool("debug"); debug {
				logrus.SetLevel(logrus.DebugLevel)
			}

			keepContainers, _ := cmd.Flags().GetBool("keep-containers")

			BuildCwd(keepContainers)
		},
	}
	cmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	cmd.PersistentFlags().Bool("keep-containers", false, "Keep containers on failure")

	return cmd
}

func BuildCwd(keepContainers bool) {
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
		DockerClient:   dockerClient,
		KeepContainers: keepContainers,
	}

	outputPath, err := builder.Build(ctx, recipe)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("%s", outputPath)
	fmt.Println()
}
