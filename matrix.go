package main
import ("fmt"
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

	//fmt.Println(m.min, "\t", m.max);
	m.min = int64(^uint64(0)>>1);
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
    if(a < 0){
    a = -a
    isNegative = !isNegative
  }
        if(b < 0){
    b = -b
    isNegative = !isNegative
  }
  //((A+B) * (C+D) = (A*C) + (C*B) + (A*D) + (B*D)

  //A*C
  al := a >> 32
  bl := b >> 32
  res := (al * bl) << 32

  ar := a & 0xffffffff
  br := b & 0xffffffff

  //C*B
  res += (bl * ar)

  //A*D
  res += (al * br)

  //B*D
  res += ((ar>>16) * (br>> 16))
    if(isNegative && res >= 0) {
    res = -res
    }    else if(!isNegative && res < 0) {
    res = -res
    }
    return res
}

func DivideFixed(a, b int64) int64{
  var isNegative bool = false;
    if(a < 0){
    a = -a
    isNegative = !isNegative
  }
    if(b < 0){
    b = -b
    isNegative = !isNegative
  }
  if(b == (b>>32)<<32){
    return a / (b>>32)
  }
  res := (a / (b >> 16)) << 16
  if(isNegative && res >= 0) {
    res = -res
    }    else if(!isNegative && res < 0) {
    res = -res
    }
    return res
}

func scale_2(x int64, n int64)int64{
  var running int64 = 1
  var i int64
  for i = 0; i < n; i++ {
      running *= 2
  }
  return x * running
}

func abs(x int64)int64{
  if x < 0{
    return -x
  }
  return x
}

func toInt(x int64)int64{
  var isNegative = x < 0
  x = abs(x)
  x = x >> 32
  if isNegative{
    x = -x
  }
  return x
}

const LN2 int64 = 0xB17217F8
const LN2_H int64 = 0xB17217F7
const LN2_L int64 = 0x1

const INV_LN2 int64 = 0x171547653
const INT_LN2_H int64 = 0x171547600
const INT_LN2_L int64 = 0x52

const ONE_HALF int64 = 0x80000000
const ONE int64 = 0x100000000
const TWO int64 = 0x200000000

const P1 int64 = 0x2AAAAAAB
const P2 int64 = -0x00B60B61
const P3 int64 = 0x0004559B
const P4 int64 = -0x00001BBD
const P5 int64 = 0x000000B2

func exp(x int64)int64{
  var hi int64
  var lo int64
  var k int64
  var t int64 = abs(x);
  if(t > LN2 / 2){
    if(t < MultiplyFixed(ONE_HALF + ONE, LN2)){
      hi = t - LN2_H
      lo = LN2_L
      k = 1
    } else {
      k = toInt((MultiplyFixed(INV_LN2, t) + ONE_HALF))
      k_fixed := k << 32
      hi = t - MultiplyFixed(k_fixed, LN2_H)
      lo = MultiplyFixed(k_fixed, LN2_L)
    }
    if(x < 0){
      hi = -hi
      lo = -lo
      k = -k
    }
    x = hi - lo
  } else if(t < 0x10){
    return 0x1 << 32
  } else{
    lo = 0
    hi = 0
    k = 0
  }
  //now x is in primary range.
  t = MultiplyFixed(x, x)
  P4_5 := P4 + MultiplyFixed(t, P5)
  P3_5 := P3 + MultiplyFixed(t, P4_5)
  P2_5 := P2 + MultiplyFixed(t, P3_5)
  P1_5 := P1 + MultiplyFixed(t, P2_5)
  c := x - MultiplyFixed(t, P1_5)
  if k == 0 {
    return ONE - (lo - DivideFixed(MultiplyFixed(x , c), TWO - c) - x)
  }
  y := ONE - (lo - DivideFixed(MultiplyFixed(x , c), TWO - c) - hi) 
  return scale_2(y, k)
}

func intToFixed(a int) int64{
	return (int64(a) << 32);
}

func printFixed(a int64){
  if(a < 0){
    a = -a
    fmt.Print("-")
  }
  fmt.Print(a >> 32)
  br := Reverse64(uint64(a))>>32
  var sum float64 = 0
  var v float64 = 0.5
  for i:=0; i < 32; i++ {
        if br% 2 == 1{
      sum += v
        }
        br = br >> 1;
    v = v/2
    }
  s := fmt.Sprintf("%16.15f", sum)
  fmt.Print(string(s[1:]))
  fmt.Print("\n")
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
