// Code generated by "stringer -type ServiceOpportunities enums.go"; DO NOT EDIT.

package main

import "strconv"

const _ServiceOpportunities_name = "OutreachService"

var _ServiceOpportunities_index = [...]uint8{0, 8, 15}

func (i ServiceOpportunities) String() string {
	i -= 1
	if i < 0 || i >= ServiceOpportunities(len(_ServiceOpportunities_index)-1) {
		return "ServiceOpportunities(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ServiceOpportunities_name[_ServiceOpportunities_index[i]:_ServiceOpportunities_index[i+1]]
}