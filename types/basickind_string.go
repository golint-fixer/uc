// Code generated by "stringer -type BasicKind"; DO NOT EDIT

package types

import "fmt"

const _BasicKind_name = "InvalidCharIntVoid"

var _BasicKind_index = [...]uint8{0, 7, 11, 14, 18}

func (i BasicKind) String() string {
	if i < 0 || i >= BasicKind(len(_BasicKind_index)-1) {
		return fmt.Sprintf("BasicKind(%d)", i)
	}
	return _BasicKind_name[_BasicKind_index[i]:_BasicKind_index[i+1]]
}
