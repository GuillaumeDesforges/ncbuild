package builder

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
)

type IBuilder interface {
	Build(ctx context.Context, recipe Recipe) (string, error)
}

type Builder struct {
	DockerClient   *client.Client
	Store          IStore
	KeepContainers bool
}

func (b *Builder) Build(ctx context.Context, recipe Recipe) (string, error) {
	var err error

	recipeHash := recipe.Hash()
	outputPath := b.Store.GetOutputPath(recipe)
	logrus.Infof("Building %s\n", recipe.Name)
	logrus.Infof("Output path: %s\n", outputPath)

	containerArgs := make([]string, len(recipe.Args)+1)
	containerArgs[0] = recipe.Executable
	copy(containerArgs[1:], recipe.Args)

	logrus.Infof("Pulling image %s\n", recipe.BuildDockerImage)
	imagePullOut, err := b.DockerClient.ImagePull(ctx, recipe.BuildDockerImage, image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer imagePullOut.Close()
	io.Copy(os.Stderr, imagePullOut)

	mounts := make([]mount.Mount, len(recipe.Inputs))
	for iInput, input := range recipe.Inputs {
		mounts[iInput] = mount.Mount{
			Type:     mount.TypeBind,
			ReadOnly: true,
			Source:   input,
			Target:   input,
		}
	}

	buildContainer, err := b.DockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:        recipe.BuildDockerImage,
			Cmd:          containerArgs,
			AttachStdout: true,
			AttachStderr: true,
			Env: []string{
				fmt.Sprintf("out=%s", outputPath),
			},
		},
		&container.HostConfig{
			Mounts: mounts,
		},
		nil,
		nil,
		fmt.Sprintf("ncbuild-%s", recipeHash),
	)
	if err != nil {
		return "", err
	}
	logrus.Debugf("Created container %s\n", buildContainer.ID)

	defer func() {
		if err == nil || (err != nil && !b.KeepContainers) {
			logrus.Debug("Removing container")
			err := b.DockerClient.ContainerRemove(
				ctx,
				buildContainer.ID,
				container.RemoveOptions{
					Force: true,
				},
			)
			if err != nil {
				logrus.Errorf("Failed to remove container %s: %s", buildContainer.ID, err)
			}
		} else {
			logrus.Warnf("Keeping container %s", buildContainer.ID)
		}
	}()

	// run the container
	err = b.DockerClient.ContainerStart(
		ctx,
		buildContainer.ID,
		container.StartOptions{},
	)
	if err != nil {
		return "", err
	}
	responseCh, errCh := b.DockerClient.ContainerWait(
		ctx,
		buildContainer.ID,
		container.WaitConditionNotRunning,
	)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case response := <-responseCh:
		if response.StatusCode != 0 {
			err = fmt.Errorf("Container exited with status %d", response.StatusCode)
			return "", err
		}
	}

	// copy the output from the container to host's store
	output, outStat, err := b.DockerClient.CopyFromContainer(
		ctx,
		buildContainer.ID,
		outputPath,
	)
	if err != nil {
		return "", err
	}
	defer output.Close()

	os.MkdirAll(b.Store.GetStorePath(), 0755)
	err = archive.CopyTo(output, archive.CopyInfo{
		Path:   "/out",
		IsDir:  outStat.Mode.IsDir(),
		Exists: true,
	}, outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
