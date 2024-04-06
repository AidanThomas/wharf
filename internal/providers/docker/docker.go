package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type State int

const (
	Running State = iota
	Exited
	Restarting
)

type Client struct {
	client *client.Client
}

type Container struct {
	State  State
	Image  string
	Names  string
	Status string
	Ports  string
	ID     string
}

func NewClient() (*Client, error) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return &Client{}, err
	}
	return &Client{client: c}, nil
}

func (c *Client) GetAll() ([]Container, error) {
	containers, err := c.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var out []Container
	for _, ctr := range containers {
		out = append(out, parseContainer(ctr))
	}

	return out, nil
}

func (c *Client) GetById(id string) (Container, error) {
	filters := filters.NewArgs(filters.KeyValuePair{Key: "id", Value: id})
	containers, err := c.client.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return Container{}, err
	}
	if len(containers) > 0 {
		return parseContainer(containers[0]), nil
	}

	return Container{}, errors.New("no container found by that id")
}

func (c *Client) SearchByName(term string) ([]Container, error) {
	filters := filters.NewArgs()
	filters.Add("name", term)
	containers, err := c.client.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	var out []Container
	for _, ctr := range containers {
		out = append(out, parseContainer(ctr))
	}

	return out, nil
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

func parseContainer(ctr types.Container) Container {
	out := Container{}
	switch ctr.State {
	case "running":
		out.State = Running
	case "exited":
		out.State = Exited
	case "restarting":
		out.State = Restarting
	}

	out.Image = ctr.Image
	out.Names = strings.Join(ctr.Names, ", ")
	out.Status = ctr.Status

	var ports []string
	for _, p := range ctr.Ports {
		var port string
		if p.IP != "" {
			port = fmt.Sprintf("%s:%d->%d/%s", p.IP, p.PublicPort, p.PrivatePort, p.Type)
		} else {
			port = fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
		}
		ports = append(ports, port)

		if len(ports) > 3 {
			ports = ports[:3]
			ports = append(ports, "...")
		}
	}

	out.Ports = strings.Join(ports, ", ")
	out.ID = ctr.ID

	return out
}
