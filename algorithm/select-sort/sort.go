package select_sort

func SelectSort(data []int) []int {
	index := 0
	for i := 0; i < len(data); i++ {
		index = i
		for j := i+1; j < len(data); j++ {
			if data[j] < data[index] {
				index = j
			}
		}
		data[i],data[index] = data[index],data[i]
	}

	return data
}