package authorization_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccDataSourceAzureRMClientConfig_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_client_config", "current")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	clientId := os.Getenv("ARM_CLIENT_ID")
	tenantId := os.Getenv("ARM_TENANT_ID")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: `
data "azurerm_client_config" "current" { }
`,
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("client_id").HasValue(clientId),
				check.That(data.ResourceName).Key("tenant_id").HasValue(tenantId),
				check.That(data.ResourceName).Key("subscription_id").HasValue(subscriptionId),
				resource.TestMatchResourceAttr(data.ResourceName, "object_id", regexp.MustCompile("^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$")),
			),
		},
	})
}
