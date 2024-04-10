package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

type Query struct {
	Size   bool
	All    bool
	Latest bool
	Since  string
	Before string
	Limit  int
	Name   string
	Id     string
}

func (q *Query) ParseQuery() container.ListOptions {
	filters := filters.NewArgs()
	if q.Name != "" {
		filters.Add("name", q.Name)
	}
	if q.Id != "" {
		filters.Add("id", q.Id)
	}

	return container.ListOptions{
		Size:    q.Size,
		All:     q.All,
		Latest:  q.Latest,
		Since:   q.Since,
		Before:  q.Before,
		Limit:   q.Limit,
		Filters: filters,
	}
}

func QueryAll() Query {
	return Query{All: true}
}

func QueryById(id string) Query {
	return Query{All: true, Id: id}
}

func QueryByName(name string) Query {
	return Query{All: true, Name: name}
}
