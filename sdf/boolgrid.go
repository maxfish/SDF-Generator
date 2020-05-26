package sdf

// BoolGrid is a grid of booleans used to represent if a pixel of an image is to be considered "inside" the shape
type BoolGrid struct {
	W, H int
	grid []bool
}

func NewBoolGrid(w, h int) *BoolGrid {
	return &BoolGrid{
		grid : make([]bool, w*h, w*h),
		W: w,
		H: h,
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
