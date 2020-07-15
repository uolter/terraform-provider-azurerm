package example

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/common"
)

var _ common.DataSource = ResourceGroupDataSource{}

func TestValidateDataSource(t *testing.T) {
	common.ValidateDataSource(t, ResourceGroupDataSource{})
}
