// Code generated by "stringer -type=opcode"; DO NOT EDIT.

package vminstance

import "strconv"

const _opcode_name = "haltOp"

var _opcode_index = [...]uint8{0, 6}

func (i opcode) String() string {
	if i >= opcode(len(_opcode_index)-1) {
		return "opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _opcode_name[_opcode_index[i]:_opcode_index[i+1]]
}
