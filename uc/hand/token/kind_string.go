// Code generated by "stringer -type Kind"; DO NOT EDIT

package token

import "fmt"

const _Kind_name = "EOFErrorCommentliteralStartIdentIntCharliteralEndoperatorStartAddSubMulDivAssignEqNeLtLeGtGeLandNotLparenRparenLbracketRbracketLbraceRbraceCommaSemicolonoperatorEndkeywordStartKwElseKwIfKwReturnKwWhilekeywordEnd"

var _Kind_index = [...]uint8{0, 3, 8, 15, 27, 32, 35, 39, 49, 62, 65, 68, 71, 74, 80, 82, 84, 86, 88, 90, 92, 96, 99, 105, 111, 119, 127, 133, 139, 144, 153, 164, 176, 182, 186, 194, 201, 211}

func (i Kind) String() string {
	if i >= Kind(len(_Kind_index)-1) {
		return fmt.Sprintf("Kind(%d)", i)
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
