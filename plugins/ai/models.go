package ai

import (
	"github.com/danielmiessler/fabric/common"
)

func NewVendorsModels() *VendorsModels {
	return &VendorsModels{GroupsItemsSelectorString: common.NewGroupsItemsSelectorString("Available models")}
}

type VendorsModels struct {
	*common.GroupsItemsSelectorString
}
