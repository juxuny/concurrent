package concurrent

import "fmt"

type CompareResult int8

const (
	Less    = CompareResult(-1)
	Equal   = CompareResult(0)
	Greater = CompareResult(1)
)

type Comparable interface {
	Less(v Comparable) CompareResult
	String() string
	Int64() int64
}

type ID int64

func (t ID) Int64() int64 {
	return int64(t)
}

func (t ID) String() string {
	return fmt.Sprintf("%v", int64(t))
}

func (t ID) Less(a Comparable) CompareResult {
	if t.Int64()-a.Int64() < int64(0) {
		return Less
	} else if t.Int64()-a.Int64() == int64(0) {
		return Equal
	} else {
		return Greater
	}
}
