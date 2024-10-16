package common

import (
	"fmt"
)

func NewGroupsItemsSelector[T any](selectionLabel string, getItemLabel func(T) string) *GroupsItemsSelector[T] {
	return &GroupsItemsSelector[T]{SelectionLabel: selectionLabel, GetItemLabel: getItemLabel, GroupsItems: make(map[string][]T)}
}

type GroupsItemsSelector[T any] struct {
	SelectionLabel string
	GetItemLabel   func(T) string

	Groups      []string
	GroupsItems map[string][]T
	Errs        []error
}

func (o *GroupsItemsSelector[T]) AddGroupItems(group string, items ...T) {
	o.Groups = append(o.Groups, group)
	o.GroupsItems[group] = items
}

func (o *GroupsItemsSelector[T]) GetGroupAndItemByItemNumber(itemIndex int) (group string, item T) {
	groupItemIndexFrom := 0
	groupItemIndexTo := 0
	for _, currenGroup := range o.Groups {
		groupItemIndexFrom = groupItemIndexTo + 1
		groupItemIndexTo = groupItemIndexFrom + len(o.GroupsItems[currenGroup]) - 1

		if itemIndex >= groupItemIndexFrom && itemIndex <= groupItemIndexTo {
			group = currenGroup
			item = o.GroupsItems[currenGroup][itemIndex-groupItemIndexFrom]
			break
		}
	}
	return
}

func (o *GroupsItemsSelector[T]) AddError(err error) {
	o.Errs = append(o.Errs, err)
}

func (o *GroupsItemsSelector[T]) Print() {
	fmt.Printf("\n%v:\n", o.SelectionLabel)

	var currentItemIndex int
	for _, group := range o.Groups {
		fmt.Println()
		fmt.Printf("%s\n", group)
		fmt.Println()
		currentItemIndex = o.PrintGroup(group, currentItemIndex)
	}
}

func (o *GroupsItemsSelector[T]) PrintGroup(group string, itemIndex int) (currentItemIndex int) {
	currentItemIndex = itemIndex
	items := o.GroupsItems[group]
	for _, item := range items {
		currentItemIndex++
		fmt.Printf("\t[%d]\t%s\n", currentItemIndex, o.GetItemLabel(item))
	}
	fmt.Println()
	return
}

func (o *GroupsItemsSelector[T]) GetGroupItems(group string) (items []T) {
	items = o.GroupsItems[group]
	return
}

func (o *GroupsItemsSelector[T]) HasGroup(group string) (ret bool) {
	ret = o.GroupsItems[group] != nil
	return
}

func (o *GroupsItemsSelector[T]) FindGroupsByItemFirst(item T) (ret string) {
	groups := o.FindGroupsByItem(item)
	if len(groups) > 0 {
		ret = groups[0]
	}
	return
}

func (o *GroupsItemsSelector[T]) FindGroupsByItem(item T) (groups []string) {
	for group, items := range o.GroupsItems {
		for _, groupItem := range items {
			if o.GetItemLabel(groupItem) == o.GetItemLabel(item) {
				groups = append(groups, group)
				continue
			}
		}
	}
	return
}
