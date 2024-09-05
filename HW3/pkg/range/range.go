package rangeI

import "fmt"

type RangeInt interface {
	Length() int
	Intersect(other RangeInt)
	Union(other RangeInt) bool
	IsEmpty() bool
	ContainsInt(i int) bool
	ContainsRange(other RangeInt) bool
	IsIntersect(other RangeInt) bool
	ToSlice() []int
	Minimum() (int, bool)
	Maximum() (int, bool)
	String() string
}

func NewRangeInt(a, b int) *RangeIntImpl {
	return &RangeIntImpl{
		from: a,
		to:   b,
	}
}

type RangeIntImpl struct {
	from int
	to   int
}

func absInt(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func (r *RangeIntImpl) Length() int {
	if r.IsEmpty() {
		return 0
	} else {
		return absInt(r.to - r.from + 1)
	}
}

func (r *RangeIntImpl) Intersect(other RangeInt) {
	o := other.(*RangeIntImpl)
	r.from = max(r.from, o.from)
	r.to = min(r.to, o.to)
}

func (r *RangeIntImpl) Union(other RangeInt) bool {
	o := other.(*RangeIntImpl)
	if r.IsEmpty() && other.IsEmpty() {
		return true
	}
	if r.IsEmpty() {
		r.from = o.from
		r.to = o.to
		return true
	}

	if !((r.to < (o.from - 1)) || (r.from-1) > (o.to)) {
		r.from = min(r.from, o.from)
		r.to = max(r.to, o.to)
		return true
	} else {
		return false
	}
}

func (r *RangeIntImpl) IsEmpty() bool {
	return r.from > r.to
}

func (r *RangeIntImpl) ContainsInt(i int) bool {
	return r.from <= i && i <= r.to
}

func (r *RangeIntImpl) ContainsRange(other RangeInt) bool {
	o := other.(*RangeIntImpl)
	return r.from <= o.from && o.to <= r.to
}

func (r *RangeIntImpl) IsIntersect(other RangeInt) bool {
	o := other.(*RangeIntImpl)

	if r.IsEmpty() && o.IsEmpty() {
		return false
	}
	if r.IsEmpty() || o.IsEmpty() {
		return false
	}
	return !((r.to < (o.from)) || (r.from) > (o.to))
}

func (r *RangeIntImpl) ToSlice() []int {
	result := make([]int, r.Length())
	val := r.from
	for i := 0; i < len(result); i++ {
		result[i] = val
		val++
	}
	return result
}

func (r *RangeIntImpl) Minimum() (int, bool) {
	return r.from, !r.IsEmpty()
}

func (r *RangeIntImpl) Maximum() (int, bool) {
	return r.to, !r.IsEmpty()
}

func (r *RangeIntImpl) String() string {
	if r.IsEmpty() {
		return ""
	} else {
		return fmt.Sprintf("[%d,%d]", r.from, r.to)
	}
}
