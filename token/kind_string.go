// Code generated by "stringer -type Kind"; DO NOT EDIT

package token

import "fmt"

const _Kind_name = "EOFErrorCommentliteralStartIdentIntLitCharLitliteralEndoperatorStartAddSubMulDivAssignEqNeLtLeGtGeLandNotLparenRparenLbracketRbracketLbraceRbraceCommaSemicolonoperatorEndkeywordStartKwElseKwIfKwReturnKwWhilekeywordEnd"

var _Kind_index = [...]uint8{0, 3, 8, 15, 27, 32, 38, 45, 55, 68, 71, 74, 77, 80, 86, 88, 90, 92, 94, 96, 98, 102, 105, 111, 117, 125, 133, 139, 145, 150, 159, 170, 182, 188, 192, 200, 207, 217}

func (kind Kind) GoString() string {
	if kind >= Kind(len(_Kind_index)-1) {
		return fmt.Sprintf("Kind(%d)", kind)
	}
	return _Kind_name[_Kind_index[kind]:_Kind_index[kind+1]]
}
