package ai

import (
	"github.com/danielmiessler/fabric/common"
)

func NewVendorsModels() *VendorsModels {
	return &VendorsModels{GroupsItemsSelector: common.NewGroupsItemsSelector("Available models", func(item string) string {
		return item
	})}
}

type VendorsModels struct {
	*common.GroupsItemsSelector[string]
}
