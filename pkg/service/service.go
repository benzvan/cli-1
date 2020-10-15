package service

import (
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/go-fastly/fastly"
)

func getServiceIDFromServiceName(name string, g *config.Data) (string, error) {
	searchInput := fastly.SearchServiceInput{
		Name: name,
	}
	service, err := g.Client.SearchService(&searchInput)
	if err != nil {
		return "", err
	}

	return service.ID, nil
}
