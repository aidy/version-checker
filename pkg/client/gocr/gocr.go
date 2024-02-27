package gocr

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/jetstack/version-checker/pkg/api"
)

type Client struct {
	FetchMetadata bool
}

func New() (*Client, error) {
	return &Client{
		FetchMetadata: true,
	}, nil
}

func (c *Client) Name() string {
	return "go-containerregistry"
}

func (c *Client) Tags(ctx context.Context, host, repo, image string) ([]api.ImageTag, error) {
	src := fmt.Sprintf("%s/%s/%s", host, repo, image)
	rpo, err := name.NewRepository(src)

	if err != nil {
		return []api.ImageTag{}, err
	}

	bareTags, err := remote.List(rpo, remote.WithContext(ctx))

	if err != nil {
		return []api.ImageTag{}, err
	}

	var tags []api.ImageTag
	for _, t := range bareTags {
		tag := api.ImageTag{Tag: t}

		if c.FetchMetadata {
			img, err := name.ParseReference(fmt.Sprintf("%s:%s", src, t))
			if err != nil {
				continue
			}
			image, err := remote.Image(img, remote.WithContext(ctx))
			if err != nil {
				continue
			}
			if digest, err := image.Digest(); err == nil {
				tag.SHA = digest.String()
			}
			if conf, err := image.ConfigFile(); err == nil {
				tag.Architecture = api.Architecture(conf.Architecture)
				tag.OS = api.OS(conf.OS)
				tag.Timestamp = conf.Created.Time
			}
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (c *Client) IsHost(host string) bool {
	return true
}

func (c *Client) RepoImageFromPath(path string) (string, string) {
	split := strings.Split(path, "/")

	lenSplit := len(split)
	if lenSplit == 1 {
		return "", split[0]
	}

	if lenSplit > 1 {
		return split[lenSplit-2], split[lenSplit-1]
	}

	return path, ""
}
