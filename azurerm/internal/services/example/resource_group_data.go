package example

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/common"
)

type ResourceGroupDataSource struct {
}

func (ResourceGroupDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func (ResourceGroupDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"location": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (ResourceGroupDataSource) Name() string {
	return "azurerm_resource_group_2"
}

type ResourceGroupArguments struct {
	Name string `json:"name"`
}

func (ResourceGroupDataSource) Read(ctx context.Context, config *common.TerraformConfiguration) error {
	client := config.Client.Resource.GroupsClient

	input := ResourceGroupArguments{}
	if err := config.DeserializeIntoType(&input); err != nil {
		return err
	}

	name := input.Name
	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("reading resource group: %+v", err)
	}

	config.Logger.Printf("this is a warning, which wants an abstraction")
	config.ResourceData.SetId(*resp.ID)
	config.ResourceData.Set("location", location.NormalizeNilable(resp.Location))
	return nil
}

func (ResourceGroupDataSource) ReadTimeout() time.Duration {
	return time.Minute * 5
}
