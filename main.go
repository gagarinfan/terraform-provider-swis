package main

import (
	"github.com/hashicorp/terraform/plugin"
	"terraform-provider-ipam/ipam"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ipam.Provider})
}
