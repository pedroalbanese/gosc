package sheet

import (
	"regexp"
	"strconv"
	"strings"
)

func adjustFormula(formula string, rowDel, colDel int) string {
	rx := regexp.MustCompile(`\$?[A-Z]+\$?\d+:\$?[A-Z]+\$?\d+|\$?[A-Z]+\$?\d+`)

	return rx.ReplaceAllStringFunc(formula, func(match string) string {
		if strings.Contains(match, ":") {
			parts := strings.SplitN(match, ":", 2)
			start := parts[0]
			end := parts[1]

			startAddr := Address(start)
			endAddr := Address(end)
			singleCol := startAddr.Column() == endAddr.Column()
			singleRow := startAddr.Row() == endAddr.Row()

			newStart := adjustRef(start, rowDel, colDel, false, true)
			newEnd := adjustRef(end, rowDel, colDel, true, true)

			if colDel >= 0 && singleCol && startAddr.Column() == colDel {
				return "@ERR"
			}
			if rowDel >= 0 && singleRow && startAddr.Row() == rowDel {
				return "@ERR"
			}

			if newStart == "@ERR" || newEnd == "@ERR" {
				return "@ERR"
			}
			return newStart + ":" + newEnd
		}
		// Individual reference: if it's in the deleted column/row, return @ERR.
		return adjustRef(match, rowDel, colDel, false, false)
	})
}

func adjustRef(ref string, rowDel, colDel int, isEnd bool, isRange bool) string {
	colAbs := strings.HasPrefix(ref, "$")
	ref = strings.TrimPrefix(ref, "$")

	colEnd := 0
	for colEnd < len(ref) && ref[colEnd] >= 'A' && ref[colEnd] <= 'Z' {
		colEnd++
	}
	colStr := ref[:colEnd]
	rowStr := ref[colEnd:]
	rowAbs := strings.HasPrefix(rowStr, "$")
	rowStr = strings.TrimPrefix(rowStr, "$")

	row, err := strconv.Atoi(rowStr)
	if err != nil {
		return ref
	}

	addr := Address(colStr + "0")
	col := addr.Column()

	if colDel >= 0 {
		if colAbs {
			// not adjusted
		} else if col == colDel {
			if isRange && !isEnd {
				// start of a range: keep it
			} else if isRange && isEnd {
				col-- // shrink the range
			} else {
				return "@ERR" // individual reference
			}
		} else if col > colDel {
			col--
		}
	}

	if rowDel >= 0 {
		if rowAbs {
			// not adjusted
		} else if row == rowDel {
			if isRange && !isEnd {
				// keep it
			} else if isRange && isEnd {
				row--
			} else {
				return "@ERR"
			}
		} else if row > rowDel {
			row--
		}
	}

	newAddr := NewAddress(row, col)
	result := newAddr.String()
	if colAbs {
		result = "$" + result
	}
	if rowAbs {
		idx := 0
		for idx < len(result) && (result[idx] == '$' || (result[idx] >= 'A' && result[idx] <= 'Z')) {
			idx++
		}
		result = result[:idx] + "$" + result[idx:]
	}
	return result
}
