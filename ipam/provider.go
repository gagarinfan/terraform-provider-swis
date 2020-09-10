package ipam

import (
	"github.com/mrxinu/gosolar"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname of solarwinds ipam instance",
				DefaultFunc: schema.EnvDefaultFunc("SOLAR_HOST", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username to solarwindws ipam instance",
				DefaultFunc: schema.EnvDefaultFunc("SOLAR_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password to solarwinds ipam instance",
				DefaultFunc: schema.EnvDefaultFunc("SOLAR_PASSWORD", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"ipam": resourceIp(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	
	return gosolar.NewClient(host, username, password, true), nil
}
