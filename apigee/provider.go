package apigee

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("APIGEE_USERNAME", nil),
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"password"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("APIGEE_PASSWORD", nil),
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"username"},
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("APIGEE_ACCESS_TOKEN", nil),
				ConflictsWith: []string{"username", "password", "oauth_server"},
			},
			"server": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_SERVER", client.PublicApigeeServer),
			},
			"server_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_SERVER_PATH", client.ServerPath),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("APIGEE_PORT", 443),
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"oauth_server": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("APIGEE_OAUTH_SERVER", nil),
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"username", "password"},
			},
			"oauth_server_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_OAUTH_SERVER_PATH", nil),
			},
			"oauth_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("APIGEE_OAUTH_PORT", 443),
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_ORGANIZATION", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"apigee_user":                       resourceUser(),
			"apigee_role":                       resourceRole(),
			"apigee_user_role":                  resourceUserRole(),
			"apigee_role_permission":            resourceRolePermission(),
			"apigee_cache":                      resourceCache(),
			"apigee_organization_kvm":           resourceOrganizationKVM(),
			"apigee_environment_kvm":            resourceEnvironmentKVM(),
			"apigee_proxy_kvm":                  resourceProxyKVM(),
			"apigee_target_server":              resourceTargetServer(),
			"apigee_virtual_host":               resourceVirtualHost(),
			"apigee_proxy":                      resourceProxy(),
			"apigee_proxy_deployment":           resourceProxyDeployment(),
			"apigee_shared_flow":                resourceSharedFlow(),
			"apigee_shared_flow_deployment":     resourceSharedFlowDeployment(),
			"apigee_developer":                  resourceDeveloper(),
			"apigee_product":                    resourceProduct(),
			"apigee_company":                    resourceCompany(),
			"apigee_company_developer":          resourceCompanyDeveloper(),
			"apigee_developer_app":              resourceDeveloperApp(),
			"apigee_developer_app_credential":   resourceDeveloperAppCredential(),
			"apigee_company_app":                resourceCompanyApp(),
			"apigee_company_app_credential":     resourceCompanyAppCredential(),
			"apigee_organization_resource_file": resourceOrganizationResourceFile(),
			"apigee_environment_resource_file":  resourceEnvironmentResourceFile(),
			"apigee_proxy_resource_file":        resourceProxyResourceFile(),
			"apigee_proxy_policy":               resourceProxyPolicy(),
			"apigee_alias":                      resourceAlias(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"apigee_user": dataSourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	accessToken := d.Get("access_token").(string)
	server := d.Get("server").(string)
	serverPath := d.Get("server_path").(string)
	port := d.Get("port").(int)
	oauthServer := d.Get("oauth_server").(string)
	oauthServerPath := d.Get("oauth_server_path").(string)
	oauthPort := d.Get("oauth_port").(int)
	organization := d.Get("organization").(string)

	//Check for valid authentication
	if (username == "") && (password == "") && (accessToken == "") {
		return nil, diag.Errorf("You must specify either username/password for Basic Authentication, username/password/oauth_server for OAuth Authentication, or access_token")
	}

	var diags diag.Diagnostics
	c, err := client.NewClient(username, password, accessToken, server, serverPath, port, oauthServer, oauthServerPath, oauthPort, organization)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
