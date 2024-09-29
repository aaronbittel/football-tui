package algorithm

type ColumnGraphData struct {
	nums   []int
	colors map[int]string
	desc   string
}

func NewColumnGraphData(nums []int, m map[int]string, desc string) ColumnGraphData {
	if m == nil {
		m = make(map[int]string)
	}
	return ColumnGraphData{
		nums:   nums,
		colors: m,
		desc:   desc,
	}
}

func (c ColumnGraphData) Desc() string {
	return c.desc
}

func (c ColumnGraphData) Nums() []int {
	return c.nums
}

func (c ColumnGraphData) Colors() map[int]string {
	return c.colors
}
