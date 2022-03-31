package utils

import (
	"math/rand"
	"time"
)

type Random struct {
	*rand.Rand
}

func NewRandom() *Random {
	random := &Random{}
	random.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return random
}

func (r *Random) RandIntn(num int) int {
	if num == 0 {
		return 0
	}
	return r.Intn(num)
}

func (r *Random) RandInt(a, b int) int {
	if a == b {
		return a
	}
	if a > b {
		return r.Intn(a-b) + b
	} else {
		return r.Intn(b-a) + a
	}
}

func (r *Random) Shuffle(arr []interface{}) {
	arrLen := len(arr)
	if arrLen <= 1 {
		return
	}

	pos := 0
	for i := 1; i < arrLen; i++ {
		idx := arrLen - i
		pos = r.Intn(idx)
		arr[pos], arr[idx] = arr[idx], arr[pos]
	}
}

//==========================================================================
var GRand *Random = NewRandom()

// arr := []int{1, 2, 3, 4, 5}
// utils.GRand.Shuffle(arr)
// fmt.Println(arr)
