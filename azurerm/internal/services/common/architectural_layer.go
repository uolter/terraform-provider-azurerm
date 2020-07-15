package common

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type TerraformConfiguration struct {
        // this is obviously Azure specific, but it could go elsewhere
	Client       *clients.Client
	ResourceData *schema.ResourceData
	Logger       *log.Logger
}

func (c TerraformConfiguration) DeserializeIntoType(input interface{}) error {
	// TODO: raise an error if it's not a pointer

	// let's assume the internal fields for attributes/schema are exposed here
	attrs := map[string]interface{}{
		"name": c.ResourceData.Get("name").(string),
	}

	// definitely ways to improve this, this is /super/ lazy but it's fine
	serialized, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	return json.Unmarshal(serialized, input)
}
