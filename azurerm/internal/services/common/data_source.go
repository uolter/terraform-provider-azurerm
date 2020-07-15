package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

type TerraformConfiguration struct {
	ResourceData *schema.ResourceData
	Logger       *log.Logger
}

func (c TerraformConfiguration) DeserializeIntoType(input interface{}) error {
	// TODO: raise an error if it's not a pointer

	// let's assume the internal field for attributes/schema is exposed here
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

type DataSource interface {
	Arguments() map[string]*schema.Schema
	Attributes() map[string]*schema.Schema
	Name() string
	Read(ctx context.Context, config *TerraformConfiguration, meta interface{}) error
	ReadTimeout() time.Duration
}

func ToSDKDataSources(input []DataSource) (*map[string]*schema.Resource, error) {
	sdkTypes := make(map[string]*schema.Resource, 0)
	for _, dataSource := range input {
		name := dataSource.Name()
		sdkType, err := toSDKDataSource(dataSource)
		if err != nil {
			return nil, fmt.Errorf("mapping %q: %+v", name, err)
		}
		sdkTypes[name] = sdkType
	}

	return &sdkTypes, nil
}

func toSDKDataSource(source DataSource) (*schema.Resource, error) {
	fields := make(map[string]*schema.Schema, 0)
	for k, v := range source.Arguments() {
		_, existing := fields[k]
		if existing {
			return nil, fmt.Errorf("Duplicate field %q", k)
		}

		fields[k] = v
	}
	for k, v := range source.Attributes() {
		_, existing := fields[k]
		if existing {
			return nil, fmt.Errorf("Duplicate field %q", k)
		}

		fields[k] = v
	}

	readTimeout := source.ReadTimeout()
	resource := schema.Resource{
		Schema: fields,
		Timeouts: &schema.ResourceTimeout{
			Read: &readTimeout,
		},
		Read: func(data *schema.ResourceData, meta interface{}) error {
			// TODO: switch out for CreateContext
			ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, data)
			defer cancel()

			config := &TerraformConfiguration{
				ResourceData: data,
				Logger:       log.New(os.Stdout, "HEYO", 1),
			}
			return source.Read(ctx, config, meta)
		},
	}
	return &resource, nil
}

// TODO: methods for validating instances of the Data Source struct at unit test time
