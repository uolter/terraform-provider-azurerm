package example

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/common"
)

func DataSources() []common.DataSource {
	return []common.DataSource{
		ResourceGroupDataSource{},
	}
}
