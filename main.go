package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"terraform-provider-swis/swis"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: swis.Provider})
}
