package main
import ("fmt")

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
func Product(a, b *Matrix) *Matrix{
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ac != br {
		panic(ErrShape)
	}
	newData := make([][]int, ar);
	for i := range newData{
		newData[i] = make([]int, bc)
	}
	m := NewMatrix(ar,bc,newData)
	for i := 0; i < ar; i++ {
		for j := 0; j < bc; j++ {
			var sum int = 0
			for k := 0; k < ac; k++ {
				sum += a.At(i,k)*b.At(k,j)
			}
			m.set(i, j, sum)
		}
	}
	return m;
}

func (m *Matrix) T() *Matrix{
	var newCol, newRow int = m.Dims();
	newData := make([][]int, newRow);
	for i := range newData{
		newData[i] = make([]int, newCol)
	}
	for i:=0; i<newRow; i++ {
		for j:=0; j<newCol; j++{
			newData[i][j] = m.data[j][i];
		}
	}
	return NewMatrix(newRow,newCol,newData);
}

func main(){
	data  := [][]int{{1, 2, 3},   
   					 {1, 2, 3},  
   					 {1, 2, 3}}
	mat := NewMatrix(len(data),len(data[0]),data);
	transpose := mat.T();
	fmt.Println(mat);
	fmt.Println(transpose);

	a := [][]int{{1,1,1,1,1}};
	b := [][]int{{1},
				 {1},
				 {1},
				 {1},
				 {1}};
    A := NewMatrix(len(a), len(a[0]), a);
    B := NewMatrix(len(b), len(b[0]), b);
    C := Product(A,B);
    fmt.Println(C);
}
// TODO
// Matrix multplication element by element and dot product FINISHED
// add, subtract matrix FINISHED
// transpose
// scaling
