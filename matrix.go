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

const m0 = 0x5555555555555555 // 01010101 ...
const m1 = 0x3333333333333333 // 00110011 ...
const m2 = 0x0f0f0f0f0f0f0f0f // 00001111 ...
const m3 = 0x00ff00ff00ff00ff // etc.
const m4 = 0x0000ffff0000ffff

type Error struct{ string }

var (
	ErrShape               = Error{"mat: dimension mismatch"}
)

type Matrix struct {
	row, col int;
	data [][]int64;
	min, max int64;
}

/*func NewMatrix(r, c int, nums [][]int) *Matrix {
	mat := Matrix{row: r, col: c, data: nums};
	return &mat
}*/

func NewMatrix(r, c int, nums []int64) *Matrix {
	data := make([][]int64, r);
	var max int64 = 0;
	var min int64 = int64(^uint64(0)>>1);
	//fmt.Println(min);
	//fmt.Println(c);
	for i:=0; i<r; i++{
		data[i] = make([]int64, c);
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

func (m *Matrix) Set(r, c int, val int64){
	m.data[r][c] = val
	if(val < m.min) {
		m.min = val
	}
	if(val > m.max) {
		m.max = val
	}
}

func (m *Matrix) At(r, c int) int64 {
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
			foo := MultiplyFixed(a.At(r, c),b.At(r, c))
			m.Set(r, c, foo)
			if(foo < m.min) {
				m.min = foo
			}
			if(foo > m.max) {
				m.max = foo
			}
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
			foo := a.At(r, c)+b.At(r, c)
			m.Set(r, c, foo)
			if(foo < m.min) {
				m.min = foo
			}
			if(foo > m.max) {
				m.max = foo
			}
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
			foo := a.At(r, c)-b.At(r, c)
			m.Set(r, c, foo)
			if(foo < m.min) {
				m.min = foo
			}
			if(foo > m.max) {
				m.max = foo
			}
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
	newData := make([]int64, ar*bc);
	m := NewMatrix(ar,bc,newData)
	for i := 0; i < ar; i++ {
		for j := 0; j < bc; j++ {
			var sum int64 = 0
			for k := 0; k < ac; k++ {
				sum += MultiplyFixed(a.At(i,k),b.At(k,j))
			}
			m.Set(i, j, sum)
			if(sum < m.min) {
				m.min = sum
			}
			if(sum > m.max) {
				m.max = sum
			}
		}
	}
	return m;
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

func (m *Matrix) Scale(c int64){
	r,col := m.Dims();
	for i:=0; i<r; i++{
		for j:=0; j<col; j++{
			foo := MultiplyFixed(m.At(i,j),c)
			m.Set(i,j,foo);
			if(foo < m.min) {
				m.min = foo
			}
			if(foo > m.max) {
				m.max = foo
			}
		}
	}
}
// unused?
func (m *Matrix) ScaleDown(c int64){
	r,col := m.Dims();
	for i:=0; i<r; i++{
		for j:=0; j<col; j++{
			m.Set(i,j,m.At(i,j)/c);
		}
	}
}

func (m *Matrix) Apply(fn func(i, j int, v int64) int64, a *Matrix){
	ar, ac := a.Dims();
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			m.Set(r, c, fn(r, c, a.At(r, c)))
		}
	}
}

//need to account for things like negative numbers properly, fix after i optimize this to not use loops;
func MultiplyFixed(a, b int64) int64{
	var isNegative bool = false;
	if (a < 0) != (b < 0){
		isNegative = true
  }
	a &= 0x7FFFFFFFFFFFFFFF
	b &= 0x7FFFFFFFFFFFFFFF
	bL := b >> 32
	res := a * bL
	bR := int64(Reverse64(uint64(b)) >> 32)
	v:=a>> 1;
	for i:=0; i < 32; i++ {
		if bR% 2 == 1{
			res += v
		}
		v = v >>1;
		bR = bR >> 1;
	}
	if(isNegative && res >= 0) {
    res = -res
	}	else if(!isNegative && res < 0) {
    res = -res
	}
	return res
}

func ReverseBytes64(x uint64) uint64 {
	const m = 1<<64 - 1
	x = x>>8&(m3&m) | x&(m3&m)<<8
	x = x>>16&(m4&m) | x&(m4&m)<<16
	return x>>32 | x<<32
}

func Reverse64(x uint64) uint64 {
	const m = 1<<64 - 1
	x = x>>1&(m0&m) | x&(m0&m)<<1
	x = x>>2&(m1&m) | x&(m1&m)<<2
	x = x>>4&(m2&m) | x&(m2&m)<<4
	return ReverseBytes64(x)
}

// TODO
// Matrix multplication element by element and dot product FINISHED
// add, subtract matrix FINISHED
// transpose FINISHED
// scaling
