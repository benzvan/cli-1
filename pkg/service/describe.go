package service

import (
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/fastly"
)

// DescribeCommand calls the Fastly API to describe a service.
type DescribeCommand struct {
	common.Base
	manifest manifest.Data
	Input    fastly.GetServiceInput
	name     string
}

// NewDescribeCommand returns a usable command registered under the parent.
func NewDescribeCommand(parent common.Registerer, globals *config.Data) *DescribeCommand {
	var c DescribeCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)
	c.CmdClause = parent.Command("describe", "Show detailed information about a Fastly service").Alias("get")
	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("name", "Service name (ignored if SERVICE-ID provided)").Short('n').StringVar(&c.name)
	return &c
}

// Exec invokes the application logic for the command.
func (c *DescribeCommand) Exec(in io.Reader, out io.Writer) error {
	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		if c.name == "" {
			return errors.ErrNoServiceIDOrName
		}
		// search for service ID
		var err error
		serviceID, err = getServiceIDFromServiceName(c.name, c.Globals)
		if err != nil {
			return err
		}
	}
	c.Input.ID = serviceID

	service, err := c.Globals.Client.GetServiceDetails(&c.Input)
	if err != nil {
		return err
	}

	text.PrintServiceDetail(out, "", service)
	return nil
}
