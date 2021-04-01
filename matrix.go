package main
import (//"fmt"
		//"bytes"
		//"encoding/binary"
		//"errors"
		//"io"
		//"math"
)

/*type storage struct {
	Version uint32 // Keep this first.
	Form    byte   // [GST]
	Packing byte   // [BPF]
	Uplo    byte   // [AUL]
	Unit    bool
	Rows    int64
	Cols    int64
	KU      int64
	KL      int64
}*/

type Error struct{ string }

var (
	ErrShape               = Error{"mat: dimension mismatch"}
)

type Matrix struct {
	row, col int;
	data [][]int;
	min, max int;
}

/*func NewMatrix(r, c int, nums [][]int) *Matrix {
	mat := Matrix{row: r, col: c, data: nums};
	return &mat
}*/

func NewMatrix(r, c int, nums []int) *Matrix {
	data := make([][]int, r);
	max := 0;
	min := int(^uint(0)>>1);
	//fmt.Println(min);
	//fmt.Println(c);
	for i:=0; i<r; i++{
		data[i] = make([]int, c);
		if nums != nil {
			//fmt.Println(i);
			for j:=0; j<c; j++{
				//fmt.Print(j, " ");
				data[i][j] = nums[i*c + j];
				if(nums[i*c+j]>max){
					//fmt.Println("new max: ", nums[i*c+j]);
					max = nums[i*c+j];
				}
				if(nums[i*c+j] < min){
					//fmt.Println("new min: ", nums[i*c+j]);
					min = nums[i*c+j]
				}
			}
			//fmt.Println();
		}
	}
	//fmt.Println("Max: ", max);
	//fmt.Println("Min: ", min);
	mat := Matrix{row: r, col: c, data: data, min:min, max:max};
	return &mat;
}

func (m *Matrix) getMin() {
	
}

func (m *Matrix) Dims() (r,c int){
	return m.row, m.col
}

func (m *Matrix) Set(r, c, val int){
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
			m.Set(r, c, a.At(r, c)*b.At(r, c))
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
			m.Set(r, c, a.At(r, c)+b.At(r, c))
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
			m.Set(r, c, a.At(r, c)-b.At(r, c))
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
	newData := make([]int, ar*bc);
	m := NewMatrix(ar,bc,newData)
	for i := 0; i < ar; i++ {
		for j := 0; j < bc; j++ {
			var sum int = 0
			for k := 0; k < ac; k++ {
				sum += a.At(i,k)*b.At(k,j)
			}
			m.Set(i, j, sum)
		}
	}
	return m;
}

func (m *Matrix) T() *Matrix{
	var newCol, newRow int = m.Dims();
	newData := make([]int, newRow*newCol);
	for i:=0; i<newRow; i++ {
		for j:=0; j<newCol; j++{
			newData[i*newCol + j] = m.data[j][i];
		}
	}
	return NewMatrix(newRow,newCol,newData);
}

func (m *Matrix) ScaleUp(c int){
	r,col := m.Dims();
	for i:=0; i<r; i++{
		for j:=0; j<col; j++{
			m.Set(i,j,m.At(i,j)*c);
		}
	}
}

func (m *Matrix) ScaleDown(c int){
	r,col := m.Dims();
	for i:=0; i<r; i++{
		for j:=0; j<col; j++{
			m.Set(i,j,m.At(i,j)/c);
		}
	}
}

func (m *Matrix) Apply(fn func(i, j int, v int) int, a *Matrix){
	ar, ac := a.Dims();
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.Set(r, c, fn(r, c, a.At(r, c)))
		}
	}
}

// TODO
// Matrix multplication element by element and dot product FINISHED
// add, subtract matrix FINISHED
// transpose FINISHED
// scaling
