package docker

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Client struct {
	client *client.Client
}

func NewClient() (*Client, error) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return &Client{}, err
	}
	return &Client{client: c}, nil
}

func (c *Client) GetAll() ([]types.Container, error) {
	return c.client.ContainerList(context.Background(), container.ListOptions{All: true})
}

func (c *Client) GetById(id string) (types.Container, error) {
	filters := filters.NewArgs(filters.KeyValuePair{Key: "id", Value: id})
	containers, err := c.client.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return types.Container{}, err
	}
	if len(containers) > 0 {
		return containers[0], nil
	}

	return types.Container{}, errors.New("no container found by that id")
}

func (c *Client) StartContainer(id string) {
	c.client.ContainerStart(context.Background(), id, container.StartOptions{})
}

func (c *Client) StopContainer(id string) {
	c.client.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (c *Client) Close() {
	c.client.Close()
}
