package slice

func Sum(arr []int) int {
	var total int
	for _, v := range arr {
		total += v
	}
	return total
}

func SumAll(nums ...[]int) []int {
	sums := make([]int, len(nums))

	for i, v := range nums {
		sums[i] = Sum(v)
	}

	return sums
}

func SumTails(nums ...[]int) []int {
	var sums []int

	for _, v := range nums {
		if len(v) == 0 {
			sums = append(sums, 0)
		} else {
			sums = append(sums, Sum(v[1:]))
		}
	}
	return sums
}