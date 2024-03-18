package demo_provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCRNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCRNamespacesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
		},
	}
}

func dataSourceAlibabacloudStackCRNamespacesRead(data *schema.ResourceData, i interface{}) error {

	return nil
}
