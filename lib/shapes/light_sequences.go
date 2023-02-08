package shapes

import "math/rand"

type Sequence interface {
	Next() float64
}

type jitterer struct {
	i     int
	elems []float64
}

func NewJitterSequence(elems ...float64) Sequence {
	return &jitterer{
		elems: elems,
	}
}

// not concurrent safe. todo atomic?
func (j *jitterer) Next() float64 {
	f := j.elems[j.i]
	j.i++
	if j.i >= len(j.elems) {
		j.i = 0
	}
	return f
}

type randSeq struct{}

func NewRandomSequence() Sequence {
	return &randSeq{}
}

func (r *randSeq) Next() float64 {
	return rand.Float64()
}
