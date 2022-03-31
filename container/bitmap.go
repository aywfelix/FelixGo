package container

import (
	"fmt"
)

// 一个uint64有64位可以保存64个数字
// length代表bits切片长度, 保存的数字范围[0， length*64-1]
type BitMap struct {
	bits   []uint64
	length int
}

func NewBitMap(len int) *BitMap {
	bitMap := &BitMap{
		bits:   make([]uint64, len),
		length: len,
	}
	return bitMap
}

func (b *BitMap) Set(num int) {
	index := num >> 6
	pos := num & 0x3f
	// 如果不足，自动扩展长度
	if index >= b.length {
		exp := index/b.length + 1
		b.length *= exp
		bits := make([]uint64, b.length)
		copy(bits, b.bits)
		b.bits = bits
	}
	b.bits[index] |= 1 << pos
}

func (b *BitMap) Remove(num int) {
	index := num >> 6
	pos := num & 0x3f
	if index >= b.length {
		return
	}
	b.bits[index] = b.bits[index] & ^(1 << pos)
}

func (b *BitMap) Clear() {
	for i := 0; i < b.length; i++ {
		b.bits[i] = b.bits[i] & 0x40
	}
}

func (b *BitMap) IsExist(num int) bool {
	index := num >> 6
	pos := num & 0x3f
	if index >= b.length {
		return false
	}
	return (b.bits[index]&(1<<pos) != 0)
}

func (b *BitMap) Len() int { return b.length }

func (b *BitMap) Numbers() []int {
	numbers := make([]int, 64)
	for i := 0; i < b.length; i++ {
		for j := 0; j < 64; j++ {
			if b.bits[i]&(1<<j) != 0 {
				numbers = append(numbers, i*64+j)
			}
		}
	}
	return numbers
}

func (b *BitMap) String() string {
	s := ""
	numbers := b.Numbers()
	for i := 0; i < len(numbers); i++ {
		s = fmt.Sprintf("%s %d", s, numbers[i])
	}
	return s
}
