package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EncryptionScopeId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	Name               string
}

func NewEncryptionScopeID(subscriptionId, resourceGroup, storageAccountName, name string) EncryptionScopeId {
	return EncryptionScopeId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		Name:               name,
	}
}

func (id EncryptionScopeId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/encryptionScopes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.Name)
}

// EncryptionScopeID parses a EncryptionScope ID into an EncryptionScopeId struct
func EncryptionScopeID(input string) (*EncryptionScopeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EncryptionScopeId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("encryptionScopes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
