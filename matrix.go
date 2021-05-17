package main

type Error struct{ string }

var (
	ErrShape = Error{"mat: dimension mismatch"}
)

type Matrix struct {
	row, col int;
	data [][]int64;
}

func (m *Matrix) Dims() (r,c int){
	return m.row, m.col
}

func NewMatrix(r, c int, nums []int64) *Matrix {
	data := make([][]int64, r);
	for i:=0; i<r; i++{
		data[i] = make([]int64, c);
		if nums != nil {
			for j:=0; j<c; j++{
				data[i][j] = nums[i*c + j];
			}
		}
	}
	mat := Matrix{row: r, col: c, data: data};
	return &mat;
}

func (a *Matrix) MulElem(b *Matrix) *Matrix{
    ar, ac := a.Dims()
	br, bc := b.Dims()
    if ar != br || ac != bc {
		panic(ErrShape)
	}
    m := NewMatrix(ar, ac, nil)
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			foo := Multiply(a.At(r, c),b.At(r, c))
			m.Set(r, c, foo)
		}
	}
    return m
}

func (m *Matrix) At(r, c int) int64 {
	return m.data[r][c]
}

func (m *Matrix) Set(r, c int, val int64){
	m.data[r][c] = val
}

func (a *Matrix) Add(b *Matrix) *Matrix {
	ar, ac := a.Dims()
	br, bc := b.Dims()
    if ar != br || ac != bc {
		panic(ErrShape)
	}
    m := NewMatrix(ar, ac, nil)
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			foo := a.At(r, c)+b.At(r, c)
			m.Set(r, c, foo)
		}
	}
    return m
}


func (a *Matrix) Sub(b *Matrix) *Matrix {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		panic(ErrShape)
	}
    m := NewMatrix(ar, ac, nil)
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			foo := a.At(r, c)-b.At(r, c)
			m.Set(r, c, foo)
		}
	}
    return m
}

func (a *Matrix) Product(b *Matrix) *Matrix{
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ac != br {
		panic(ErrShape)
	}
	m := NewMatrix(ar,bc,nil)

	for i := 0; i < ar; i++ {
		for j := 0; j < bc; j++ {
			var sum int64 = 0
			for k := 0; k < ac; k++ {
				sum += Multiply(a.At(i,k),b.At(k,j))
			}
			m.Set(i, j, sum)
		}
	}
    return m
}

func (m *Matrix) addScalar(i int64) *Matrix {
	r, c := m.Dims();
	a := make([]int64, r*c);
	for x := 0; x < r*c; x++ {
		a[x] = i;
	}
	n := NewMatrix(r, c, a);
	return m.Add(n);
}

func (m *Matrix) T() *Matrix{
	var newCol, newRow int = m.Dims();
	newData := make([]int64, newRow*newCol);
	for i:=0; i<newRow; i++ {
		for j:=0; j<newCol; j++{
			newData[i*newCol + j] = m.data[j][i];
		}
	}
	return NewMatrix(newRow,newCol,newData);
}

func (m *Matrix) Scale(c int64) *Matrix{
	r,col := m.Dims();
    o := NewMatrix(r, col, nil)
	for i:=0; i<r; i++{
		for j:=0; j<col; j++{
			foo := Multiply(m.At(i,j),c)
			o.Set(i,j,foo);
		}
	}
    return o
}

func (a *Matrix) Apply(fn func(i, j int, v int64) int64) *Matrix {
	ar, ac := a.Dims();
    m := NewMatrix(ar,ac,nil)
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.Set(r, c, fn(r, c, a.At(r, c)))
		}
	}
    return m
}
