package demo_provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"os"
)

func Provider() *schema.Provider {
	//可以新建一个，也可以直接使用匿名函数返回
	firstProvider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ACCESS_KEY", os.Getenv("ALIBABACLOUDSTACK_ACCESS_KEY")),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTP",
				Description:  "",
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alibabacloudstack_cr_namespaces": nil,
		},
	}
	return firstProvider
}
