package sdf

// BoolGrid is a grid of booleans used to represent if a pixel of an image is to be considered "inside" the shape
type BoolGrid struct {
	W, H int
	grid []bool
}

func NewBoolGrid(w, h int) *BoolGrid {
	return &BoolGrid{
		grid: make([]bool, w*h, w*h),
		W:    w,
		H:    h,
	}
}

func (m *BoolGrid) Set(value bool, x, y int) {
	if x < 0 || x >= m.W || y < 0 || y >= m.H {
		return
	}
	m.grid[y*m.W+x] = value
}

func (m *BoolGrid) At(x, y int) bool {
	if x < 0 || y < 0 || x >= m.W || y >= m.H {
		return false
	}
	return m.grid[y*m.W+x]
}

func (m *BoolGrid) Crop() *BoolGrid {
	// Scan the grid to find the first non-empty rows/columns from each side
	y1 := -1
y1Loop:
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if m.grid[y*m.W+x] {
				y1 = y
				break y1Loop
			}
		}
	}
	y2 := m.H
y2Loop:
	for y := m.H - 1; y >= 0; y-- {
		for x := 0; x < m.W; x++ {
			if m.grid[y*m.W+x] {
				y2 = y
				break y2Loop
			}
		}
	}
	x1 := -1
x1Loop:
	for x := 0; x < m.W; x++ {
		for y := 0; y < m.H; y++ {
			if m.grid[y*m.W+x] {
				x1 = x
				break x1Loop
			}
		}
	}
	x2 := m.W
x2Loop:
	for x := m.W - 1; x >= 0; x-- {
		for y := 0; y < m.H; y++ {
			if m.grid[y*m.W+x] {
				x2 = x
				break x2Loop
			}
		}
	}

	// Build the new cropped grid
	newW := x2 - x1
	newH := y2 - y1
	bg := NewBoolGrid(newW, newH)
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			bg.Set(m.At(x+x1, y+y1), x, y)
		}
	}
	return bg
}
