package example

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/common"
)

// check complies with the interface
var _ common.DataSource = ResourceGroupDataSource{}

// check the implementation (fields) are valid
func TestValidateDataSource(t *testing.T) {
	common.ValidateDataSource(t, ResourceGroupDataSource{})
}
