package discourse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ManagementConfig struct {
	Server string
	Domain string
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("DISCOURSE_TOKEN", nil),
			},
			"username": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("DISCOURSE_USERNAME", nil),
			},
			"host": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("DISCOURSE_HOST", nil),
			},
		},
		ResourcesMap:   map[string]*schema.Resource{
			"discourse_group": resourceGroup(),
			"discourse_setting": resourceSetting(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	log("configuring Discourse terraform provider")

	// Defining local variables from provider args for ease of reading
	conf := Config{
		Host: "https://"+d.Get("host").(string),
		Username: d.Get("username").(string),
		Token: d.Get("token").(string),
	}

	request := &ApiRequest {
		Method: "GET",
		Endpoint: "/admin/backups.json",
		Config: conf,
	}

	// Note: We pull the active users list just to validate the key. We throw away the body we receive.
	log("Performing provider healthcheck.")

	_, reqErr, ok := request.Call()
	if ! ok {
		log("Healthcheck failed. Configuration invalid.")
    diags = append(diags, reqErr)
    return nil, diags
  }

	log("Healthcheck successful. Configuration validated.")

	return conf, diags
}
