package main

import (
	"bytes"
	"math"
	"math/cmplx"
)

func getCharByIndex(str string, idx int) rune { return []rune(str)[idx] }

func getStringBySliceOfIndexes(str string, indexes []int) string {
	var b bytes.Buffer

	for _, i := range indexes {
		b.WriteRune([]rune(str)[i])
	}
	return b.String()
}

func addPointers(ptr1, ptr2 *int) *int {
	if ptr1 == nil || ptr2 == nil {
		return nil
	}
	*ptr1 = *ptr1 + *ptr2
	return ptr1
}

func isComplexEqual(a, b complex128) bool {
	eps := 1e-6

	realP := math.Abs(real(a)-real(b)) < eps
	imagP := math.Abs(imag(a)-imag(b)) < eps
	return realP && imagP
}

func getRootsOfQuadraticEquation(a, b, c float64) (complex128, complex128) {

	underRoot := b*b - 4*a*c
	root := cmplx.Sqrt(complex(underRoot, 0))

	left := complex(-b/(2*a), 0)
	right := root / complex(2*a, 0)

	return left - right, left + right
}

func mergeSort(s []int) []int {
	length := len(s)

	switch {
	case length == 1:
		return s
	case length == 2:
		if s[0] > s[1] {
			swapPointers(&s[0], &s[1])
		}
	case true:
		mid := length / 2

		lSorted, rSorted := make([]int, mid), make([]int, length-mid)

		copy(lSorted, s[0:mid])
		copy(rSorted, s[mid:length])

		lArray, rArray := mergeSort(lSorted), mergeSort(rSorted)
		lLength, rLength := len(lArray), len(rArray)

		lPivot, rPivot := 0, 0

		for i := 0; i < length; i++ {
			if lPivot < lLength && (rPivot == rLength || lArray[lPivot] <= rArray[rPivot]) {
				s[i] = lArray[lPivot]
				lPivot++
			} else {
				s[i] = rArray[rPivot]
				rPivot++
			}
		}
	}
	return s
}

func reverseSliceOne(s []int) {
	for i := 0; i < len(s)/2; i++ {
		swapPointers(&s[i], &s[len(s)-1-i])
	}
}

func reverseSliceTwo(s []int) []int {
	newSlice := make([]int, len(s))
	copy(newSlice, s)

	reverseSliceOne(newSlice)
	return newSlice
}

func swapPointers(a, b *int) {
	x := *a
	*a = *b
	*b = x
}

func isSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true

}

func deleteByIndex(s []int, idx int) []int {
	return append(s[:idx], s[idx+1:]...)
}
