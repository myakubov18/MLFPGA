package main

import(
    "fmt"
)

const m0 = 0x5555555555555555 // 01010101 ...
const m1 = 0x3333333333333333 // 00110011 ...
const m2 = 0x0f0f0f0f0f0f0f0f // 00001111 ...
const m3 = 0x00ff00ff00ff00ff // etc.
const m4 = 0x0000ffff0000ffff

func printFixed(a int64){
  if(a < 0){
    a = -a
    fmt.Print("-")
  }
  fmt.Print(a >> 48)
  br := Reverse64(uint64(a))>>16
  var sum float64 = 0
  var v float64 = 0.5
  for i:=0; i < 48; i++ {
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

func Divide(a, b int64) int64{
    var isNegative bool = false;
    if(a < 0){
        a = -a
        isNegative = !isNegative
    }
    if(b < 0){
        b = -b
        isNegative = !isNegative
    }
    if(b == (b>>48)<<48){
        return a / (b>>48)
    }
    res := (a / (b >> 24)) << 24
    if(isNegative && res >= 0) {
        res = -res
    } else if(!isNegative && res < 0) {
        res = -res
    }
    return res
}

func Multiply(a, b int64) int64{
    var isNegative bool = false
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
    al := a >> 48
    bl := b >> 48
    res := (al * bl) << 48

    ar := a & 0xffffffffffff
    br := b & 0xffffffffffff

    //C*B
    res += (bl * ar)

    //A*D
    res += (al * br)

    //B*D
    res += ((ar>>24) * (br>> 24))
    if(isNegative && res >= 0) {
        res = -res
    } else if(!isNegative && res < 0) {
        res = -res
    }
    //overflow detection


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
