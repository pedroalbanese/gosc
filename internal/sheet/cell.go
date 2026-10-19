package sheet

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/pedroalbanese/sc/internal/display"
	"github.com/pedroalbanese/sc/internal/evaler"
	"github.com/pedroalbanese/sc/internal/sheet/align"
)

type Cell struct {
	value      string
	stringType bool
	alignment  align.Align

	forwardRefs map[Address]struct{} // Cells that are required for any formula

	// cached computed value (only for numeric cells)
	cachedValue string
	cacheValid  bool
}

// Creates a new cell.
func NewCell(value string, alignment align.Align, stringType bool) *Cell {
	return &Cell{
		value:       value,
		alignment:   alignment,
		stringType:  stringType,
		forwardRefs: make(map[Address]struct{}),
		cacheValid:  false,
	}
}

// Creates a copy of the cell, altering any formula that is contained based on where the
// cell is being moved. It uses the relative change between the old address and new address
// to figure out the differences to apply to every cell referenced in the formula.
func (c *Cell) Copy(oldAddress, newAddress Address) *Cell {
	val := c.value
	if !c.stringType {
		rowDiff := newAddress.Row() - oldAddress.Row()
		colDiff := newAddress.Column() - oldAddress.Column()
		for forRef, _ := range c.forwardRefs {
			refRow, refCol := forRef.RowCol()
			newRef := NewAddress(refRow+rowDiff, refCol+colDiff)
			val = strings.Replace(val, string(forRef), string(newRef), -1)
		}
	}
	return NewCell(val, c.alignment, c.stringType)
}

// getDisplayValue returns the value to display as a string for numeric cells.
// This is the value that shows up on cell in the sheet display.
// It is the result of any calculation that may have
// been required to be performed.
//
// Cell is assumed not to be string type when this function is called.
// For string type cells, simply use c.value.
func (c *Cell) getDisplayValue(s *Sheet, address Address) string {
	if c.cacheValid {
		return c.cachedValue
	}

	postfix := evaler.GetPostfix(c.value)
	for idx, token := range postfix {
		if evaler.IsCellRange(token) {
			postfix[idx] = s.expandRange(token)
		} else if evaler.IsCellAddr(token) {
			refCellVal := ""
			if cell, err := s.GetCell(Address(token)); err == nil {
				refCellVal = cell.getDisplayValue(s, Address(token))
			}
			postfix[idx] = refCellVal
		}
	}

	result := ""
	if rat, err := evaler.EvaluatePostfix(postfix); err == nil {
		result = rat.FloatString(s.getPrecision(address))
	}

	c.cachedValue = result
	c.cacheValid = true
	return result
}

// Displays cell value using termbox.
func (c *Cell) display(s *Sheet, address Address, row, colStart, colEnd int, selected bool) {
	dispVal := c.getDisplay(s, address)

	display.DisplayValue(dispVal, row, colStart, colEnd, c.alignment, selected)
}

// getDisplay gets what this cell should display
func (c *Cell) getDisplay(s *Sheet, addr Address) string {
	if c.stringType {
		return c.value
	}

	return c.getDisplayValue(s, addr)
}

// Gets the raw value in a format the also specifies any alignment defined in cell.
func (c *Cell) StatusBarVal() string {
	if c.stringType {
		modifier := ""
		switch c.alignment {
		case align.AlignLeft:
			modifier = "<"
		case align.AlignCenter:
			modifier = "|"
		default:
			modifier = ">"
		}
		return fmt.Sprintf("%s\"%s\"", modifier, c.value)
	} else {
		return c.value
	}
}

// expandRange expands a range like A2:A4 into the sum of the values
// of all cells in the rectangular region.
//
// The result is returned as a string in a format that evaler can parse.
func (s *Sheet) expandRange(rng string) string {
	parts := strings.SplitN(rng, ":", 2)
	if len(parts) != 2 {
		return "0"
	}
	start := Address(parts[0])
	end := Address(parts[1])

	r1, c1 := start.RowCol()
	r2, c2 := end.RowCol()

	sum := new(big.Rat)
	for r := r1; r <= r2; r++ {
		for c := c1; c <= c2; c++ {
			addr := NewAddress(r, c)
			cell, err := s.GetCell(addr)
			if err != nil || cell == nil {
				continue
			}
			if cell.stringType {
				// Strings count as zero in a numeric sum.
				continue
			}
			valStr := cell.getDisplayValue(s, addr)
			if valStr == "" {
				continue
			}
			v := new(big.Rat)
			if _, err := fmt.Sscan(valStr, v); err == nil {
				sum.Add(sum, v)
			}
		}
	}
	return sum.FloatString(10)
}
