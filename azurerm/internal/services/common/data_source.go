package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

type DataSource interface {
	Arguments() map[string]*schema.Schema
	Attributes() map[string]*schema.Schema
	Name() string
	Read(ctx context.Context, config *TerraformConfiguration) error
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
			client := meta.(*clients.Client)

			config := &TerraformConfiguration{
				Client:       client,
				ResourceData: data,
				Logger:       log.New(os.Stdout, "HEYO", 1),
			}
			return source.Read(ctx, config)
		},
	}
	return &resource, nil
}

func ValidateDataSource(t *testing.T, input DataSource) {
	for k, v := range input.Attributes() {
		if !v.Computed {
			t.Fatalf("%q must be computed if an Attribute", k)
		}

		if v.Required || v.Optional {
			t.Fatalf("%q cannot be Required/Optional if it's an Attribute", k)
		}
	}

	for k, v := range input.Attributes() {
		if v.Computed && !(v.Required || v.Optional) {
			t.Fatalf("%q cannot be read-only as an Argument - make it an Attribute", k)
		}
	}
}
