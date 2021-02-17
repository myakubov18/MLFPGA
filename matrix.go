package main

type Error struct{ string }

var (
	ErrShape               = Error{"mat: dimension mismatch"}
)

type Matrix struct {
	row, col int;
	data [][]int;
}

func NewMatrix(r, c int, nums [][]int) *Matrix {
	mat := Matrix{row: r, col: c, data: nums};
	return &mat
}

func (m *Matrix) Dims() (r,c int){
	return m.row, m.col
}

func (m *Matrix) set(r, c, val int){
	m.data[r][c] = val
}

func (m *Matrix) At(r, c int) int {
	return m.data[r][c]
}

func (m *Matrix) MulElem(a, b *Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)*b.At(r, c))
		}
	}
}

func (m *Matrix) Add(a, b *Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)+b.At(r, c))
		}
	}
}

func (m *Matrix) Sub(a, b *Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.set(r, c, a.At(r, c)-b.At(r, c))
		}
	}
}
// As long as a or b is not m, this works fine
func (m *Matrix) Product(a, b *Matrix){
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ac != br {
		panic(ErrShape)
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < bc; j++ {
			var sum int = 0
			for k := 0; k < ac; k++ {
				sum += a.data[i][k]*b.data[k][j]
			}
			m.set(i, j, sum)
		}
	}
}
// TODO
// Matrix multplication element by element and dot product FINISHED
// add, subtract matrix FINISHED
// transpose
// scaling
