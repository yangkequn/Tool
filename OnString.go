package Tool

import (
	"strings"
)

//RemoveItemFromString remove item from "," split array
func RemoveItemFromString(str *string, item string) bool {
	if len(*str) == 0 || len(item) == 0 {
		return false
	}
	ind := strings.Index(*str, item)
	if ind < 0 {
		return false
	}
	part1 := (*str)[0:ind]
	itemEndingIncludingSeparator := ind + len(item) + 1
	if itemEndingIncludingSeparator >= len(*str) {
		*str = part1
	} else {
		*str = part1 + (*str)[itemEndingIncludingSeparator:len(*str)]
	}
	return true
}

//concat two list  to one,using "," as separator
//items in s2 should be skipped if they are in s1
func NonRedundantMerge(list1 *string, list2 string, toFront bool) bool {
	var modified bool = false
	if len(*list1) == 0 {
		modified = *list1 != list2
		*list1 = list2
		return modified
	}
	if len(list2) == 0 {
		return false
	}
	items := strings.Split(list2, ",")
	for _, item := range items {
		if !strings.Contains(*list1, item) {
			if toFront {
				*list1 = item + "," + *list1
			} else {
				*list1 = *list1 + "," + item
			}
			modified = true
		}
	}
	return modified
}
