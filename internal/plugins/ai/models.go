package ai

import (
	"github.com/danielmiessler/fabric/internal/common"
)

func NewVendorsModels() *VendorsModels {
	return &VendorsModels{GroupsItemsSelectorString: common.NewGroupsItemsSelectorString("Available models")}
}

type VendorsModels struct {
	*common.GroupsItemsSelectorString
}
