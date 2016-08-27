package subjects

type sliceType []int

func (st *sliceType) foo() {}

func (st *sliceType) arrayfunc([]int) []int {
	return nil
}
