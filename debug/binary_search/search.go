package main

import "github.com/rs/zerolog/log"

func main() {

	result := BinarySearchIndex()
	log.Info().Msgf("result is: %d", result)

	result = BinarySearchFirstIndexWithRepeatableElement()
	log.Info().Msgf("result is: %d", result)

	result = BinarySearchLastIndexWithRepeatableElement()
	log.Info().Msgf("result is: %d", result)
}

// 二分查找
func BinarySearchIndex() (index int) {
	var dataList []int = []int{1, 2, 3, 4, 5, 6, 6, 8, 9, 10, 13, 15}
	var start = 0
	var end int = len(dataList) - 1
	var mid int

	var target = 4

	for start <= end {
		mid = start + (end-start)/2
		if dataList[mid] == target {
			index = mid
			return
		}
		if dataList[mid] < target {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	index = -1
	return
}

// 搜索可重复对象(左边第一个)
func BinarySearchFirstIndexWithRepeatableElement() (index int) {
	var dataList []int = []int{1, 2, 3, 4, 5, 6, 6, 8, 9, 10, 13, 15}
	var start = 0
	var end int = len(dataList) - 1
	var mid int

	var target = 6

	for start <= end {
		mid = start + (end-start)/2
		if dataList[mid] == target {
			end = mid - 1
		} else if dataList[mid] < target {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}

	// 因为循环跳出条件为low = high + 1, 当target比所有元素都大的话, 此时high就是最后一个元素, low会越界, 所以要判断一下
	// 当target比所有元素都小的话, 此时low不会越界, 但要判断一下nums[low]和target相不相等, 如nums数组为[1], 不判断的话会返回0
	if start >= len(dataList) || dataList[start] != target {
		return -1
	}

	// 因为不管怎么样high一定是 mid-1
	// 而循环跳出条件为low = high + 1，所以返回low
	index = start
	return
}

// 搜索可重复对象(右边边第一个)
func BinarySearchLastIndexWithRepeatableElement() (index int) {

	var dataList []int = []int{1, 2, 3, 4, 5, 6, 6, 8, 9, 10, 13, 15}
	var start = 0
	var end int = len(dataList) - 1
	var mid int

	var target = 6

	for start <= end {
		mid = start + (end-start)/2
		if dataList[mid] == target {
			start = mid + 1
		} else if dataList[mid] < target {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}

	//判断是否越界
	if end <= 0 || dataList[end] != target {
		return -1
	}

	// 因为不管怎么样low一定是 mid+1
	// 而循环跳出条件为low = high + 1, 也就是high = low - 1, 所以返回high
	index = end
	return
}
