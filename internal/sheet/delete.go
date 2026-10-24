package sheet

// DeleteRow removes the given row and shifts all rows below it up by one.
func (s *Sheet) DeleteRow(row int) {
	if !s.loading {
		s.PushUndo(NewAddress(row, 0))
	}

	maxRow, maxCol := 0, 0
	for addr := range s.data {
		if addr.Row() > maxRow {
			maxRow = addr.Row()
		}
		if addr.Column() > maxCol {
			maxCol = addr.Column()
		}
	}

	// Adjust formulas.
	for _, cell := range s.data {
		if !cell.stringType {
			cell.value = adjustFormula(cell.value, row, -1)
		}
	}

	// Move cells up.
	for r := row + 1; r <= maxRow; r++ {
		for c := 0; c <= maxCol; c++ {
			oldAddr := NewAddress(r, c)
			if cell, found := s.data[oldAddr]; found {
				delete(s.data, oldAddr)
				newAddr := NewAddress(r-1, c)
				cell.cacheValid = false
				s.data[newAddr] = cell
			}
		}
	}

	for _, cell := range s.data {
		cell.cacheValid = false
	}

	s.rebuildMaximums()

	newMaxRow := 0
	for addr := range s.data {
		if addr.Row() > newMaxRow {
			newMaxRow = addr.Row()
		}
	}
	if s.startRow > newMaxRow {
		s.startRow = newMaxRow
	}

	if s.SelectedCell.Row() >= row {
		r, c := s.SelectedCell.RowCol()
		if r > 0 {
			r--
		}
		s.SelectedCell = NewAddress(r, c)
	}

	s.display()
}

// DeleteColumn removes the given column and shifts all columns to the right left by one.
func (s *Sheet) DeleteColumn(col int) {
	if !s.loading {
		s.PushUndo(NewAddress(0, col))
	}

	maxRow, maxCol := 0, 0
	for addr := range s.data {
		if addr.Row() > maxRow {
			maxRow = addr.Row()
		}
		if addr.Column() > maxCol {
			maxCol = addr.Column()
		}
	}

	// Adjust formulas.
	for _, cell := range s.data {
		if !cell.stringType {
			cell.value = adjustFormula(cell.value, -1, col)
		}
	}

	// Move cells left.
	for c := col + 1; c <= maxCol; c++ {
		for r := 0; r <= maxRow; r++ {
			oldAddr := NewAddress(r, c)
			if cell, found := s.data[oldAddr]; found {
				delete(s.data, oldAddr)
				newAddr := NewAddress(r, c-1)
				cell.cacheValid = false
				s.data[newAddr] = cell
			}
		}
	}

	if maxCol == col {
		for r := 0; r <= maxRow; r++ {
			delete(s.data, NewAddress(r, col))
		}
	}

	for _, cell := range s.data {
		cell.cacheValid = false
	}

	s.rebuildMaximums()

	newMaxCol := 0
	for addr := range s.data {
		if addr.Column() > newMaxCol {
			newMaxCol = addr.Column()
		}
	}
	if s.startCol > newMaxCol {
		s.startCol = newMaxCol
	}

	if s.SelectedCell.Column() >= col {
		r, c := s.SelectedCell.RowCol()
		if c > 0 {
			c--
		}
		s.SelectedCell = NewAddress(r, c)
	}

	s.display()
}

// rebuildMaximums recomputes maxRowForColumn and maxColumnForRow from scratch.
func (s *Sheet) rebuildMaximums() {
	s.maxRowForColumn = make(map[int]int)
	s.maxColumnForRow = make(map[int]int)
	for addr := range s.data {
		row, col := addr.RowCol()
		if row > s.maxRowForColumn[col] {
			s.maxRowForColumn[col] = row
		}
		if col > s.maxColumnForRow[row] {
			s.maxColumnForRow[row] = col
		}
	}
}
