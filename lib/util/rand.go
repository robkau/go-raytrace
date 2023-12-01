package util

import (
	"math/rand"
	"sync"
)

var RandomSource RandPool = nil

func init() {
	RandomSource = NewRandPool()
}

type RandPool interface {
	Float64() float64
}

type randPool struct {
	p *sync.Pool
}

func (r *randPool) Float64() float64 {
	rng := r.p.Get().(*rand.Rand)
	defer r.p.Put(rng)
	return rng.Float64()
}

func NewRandPool() RandPool {
	pool := &sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(rand.Int63()))
		},
	}

	return &randPool{p: pool}
}
