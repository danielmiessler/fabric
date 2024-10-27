package common

import (
	"fmt"
	"github.com/samber/lo"
)

func NewGroupsItemsSelector[I any](selectionLabel string,
	getItemLabel func(I) string) *GroupsItemsSelector[I] {

	return &GroupsItemsSelector[I]{SelectionLabel: selectionLabel,
		GetItemKey:  getItemLabel,
		GroupsItems: make([]*GroupItems[I], 0),
	}
}

type GroupItems[I any] struct {
	Group string
	Items []I
}

func (o *GroupItems[I]) Count() int {
	return len(o.Items)
}

func (o *GroupItems[I]) ContainsItemBy(predicate func(item I) bool) (ret bool) {
	ret = lo.ContainsBy(o.Items, predicate)
	return
}

type GroupsItemsSelector[I any] struct {
	SelectionLabel string
	GetItemKey     func(I) string

	GroupsItems []*GroupItems[I]
}

func (o *GroupsItemsSelector[I]) AddGroupItems(group string, items ...I) {
	o.GroupsItems = append(o.GroupsItems, &GroupItems[I]{group, items})
}

func (o *GroupsItemsSelector[I]) GetGroupAndItemByItemNumber(number int) (group string, item I, err error) {
	var currentItemNumber int
	found := false

	for _, groupItems := range o.GroupsItems {
		if currentItemNumber+groupItems.Count() < number {
			currentItemNumber += groupItems.Count()
			continue
		}

		for _, groupItem := range groupItems.Items {
			currentItemNumber++
			if currentItemNumber == number {
				group = groupItems.Group
				item = groupItem
				found = true
				break
			}
		}
	}

	if !found {
		err = fmt.Errorf("number %d is out of range", number)
	}
	return
}

func (o *GroupsItemsSelector[I]) Print() {
	fmt.Printf("\n%v:\n", o.SelectionLabel)

	var currentItemIndex int
	for _, groupItems := range o.GroupsItems {
		fmt.Println()
		fmt.Printf("%s\n", groupItems.Group)
		fmt.Println()

		for _, item := range groupItems.Items {
			currentItemIndex++
			fmt.Printf("\t[%d]\t%s\n", currentItemIndex, o.GetItemKey(item))

		}
	}
}

func (o *GroupsItemsSelector[I]) HasGroup(group string) (ret bool) {
	for _, groupItems := range o.GroupsItems {
		if ret = groupItems.Group == group; ret {
			break
		}
	}
	return
}

func (o *GroupsItemsSelector[I]) FindGroupsByItemFirst(item I) (ret string) {
	itemKey := o.GetItemKey(item)

	for _, groupItems := range o.GroupsItems {
		if groupItems.ContainsItemBy(func(groupItem I) bool {
			groupItemKey := o.GetItemKey(groupItem)
			return groupItemKey == itemKey
		}) {
			ret = groupItems.Group
			break
		}
	}
	return
}

func (o *GroupsItemsSelector[I]) FindGroupsByItem(item I) (groups []string) {
	itemKey := o.GetItemKey(item)

	for _, groupItems := range o.GroupsItems {
		if groupItems.ContainsItemBy(func(groupItem I) bool {
			groupItemKey := o.GetItemKey(groupItem)
			return groupItemKey == itemKey
		}) {
			groups = append(groups, groupItems.Group)
		}
	}
	return
}

func ReturnItem(item string) string {
	return item
}

func NewGroupsItemsSelectorString(selectionLabel string) *GroupsItemsSelectorString {
	return &GroupsItemsSelectorString{GroupsItemsSelector: NewGroupsItemsSelector(selectionLabel, ReturnItem)}
}

type GroupsItemsSelectorString struct {
	*GroupsItemsSelector[string]
}
