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

	forwardRefs map[Address]struct{}

	cachedValue string
	cacheValid  bool
	calculating bool
}

func NewCell(value string, alignment align.Align, stringType bool) *Cell {
	return &Cell{
		value:       value,
		alignment:   alignment,
		stringType:  stringType,
		forwardRefs: make(map[Address]struct{}),
		cacheValid:  false,
	}
}

func (c *Cell) Copy(oldAddress, newAddress Address) *Cell {
	val := c.value
	if !c.stringType {
		rowDiff := newAddress.Row() - oldAddress.Row()
		colDiff := newAddress.Column() - oldAddress.Column()
		for forRef := range c.forwardRefs {
			refRow, refCol := forRef.RowCol()
			newRef := NewAddress(refRow+rowDiff, refCol+colDiff)
			val = strings.Replace(val, string(forRef), string(newRef), -1)
		}
	}
	return NewCell(val, c.alignment, c.stringType)
}

func (c *Cell) getDisplayValue(s *Sheet, address Address) string {
	if c.cacheValid {
		return c.cachedValue
	}
	if c.calculating {
		return "circle"
	}
	c.calculating = true
	defer func() { c.calculating = false }()

	postfix := evaler.GetPostfix(c.value)
	for idx, token := range postfix {
		if evaler.IsCellRange(token) {
			v := s.expandRange(token)
			if v == "circle" {
				c.cachedValue = "circle"
				c.cacheValid = true
				return "circle"
			}
			postfix[idx] = v
		} else if evaler.IsCellAddr(token) {
			refCellVal := "0"
			if cell, err := s.GetCell(Address(token)); err == nil && cell != nil {
				v := cell.getDisplayValue(s, Address(token))
				if v == "circle" {
					c.cachedValue = "circle"
					c.cacheValid = true
					return "circle"
				}
				if v != "" {
					refCellVal = v
				}
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

func (c *Cell) display(s *Sheet, address Address, row, colStart, colEnd int, selected bool) {
	dispVal := c.getDisplay(s, address)
	display.DisplayValue(dispVal, row, colStart, colEnd, c.alignment, selected)
}

func (c *Cell) getDisplay(s *Sheet, addr Address) string {
	if c.stringType {
		return c.value
	}
	return c.getDisplayValue(s, addr)
}

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
	}
	return c.value
}

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
			if err != nil || cell == nil || cell.stringType {
				continue
			}
			valStr := cell.getDisplayValue(s, addr)
			if valStr == "circle" {
				return "circle"
			}
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
