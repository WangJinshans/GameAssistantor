package main

import (
	"fmt"
	"sort"
)

type Person struct {
	LastName  string
	FirstName string
}

// ByLastFirst 用于对Person切片进行排序的类型，实现了sort.Interface
type ByLastFirst []Person

func (a ByLastFirst) Len() int      { return len(a) }
func (a ByLastFirst) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// 复杂对象多字段排序: 先根据LastName进行排序, 后根据FirstName进行排序
func (a ByLastFirst) Less(i, j int) bool {
	p := a[i]
	q := a[j]
	if p.LastName == q.LastName {
		return p.FirstName < q.FirstName
	}
	return p.LastName < q.LastName
}

func main() {
	people := []Person{
		{"Smith", "John"},
		{"Doe", "Jane"},
		{"Smith", "Jane"},
	}

	fmt.Println("Before sorting:", people)

	// 使用sort.SliceStable对people进行多字段排序
	sort.SliceStable(people, func(i, j int) bool {
		pi, pj := people[i], people[j]
		return pi.LastName < pj.LastName || (pi.LastName == pj.LastName && pi.FirstName < pj.FirstName)
	})

	fmt.Println("After sorting:", people)

	// 或者使用ByLastFirst类型进行排序
	sort.Stable(ByLastFirst(people))
	fmt.Println("After sorting with ByLastFirst:", people)
}

// 冒泡 O(n2):
// 比较相邻的元素, 如果第一个比第二个大, 就交换它们两个
// 对每一对相邻元素作同样的工作, 从开始第一对到结尾的最后一对, 这样在最后的元素应该会是最大的数

// 选择 O(n2):
// 在未排序序列中找到最小(大)元素, 存放到排序序列的起始位置
// 从剩余未排序元素中继续寻找最小(大)元素, 然后放到已排序序列的末尾

// 插入O(n2):
// 把待排序的数组分成已排序和未排序两部分, 初始的时候把第一个元素认为是已排好序的
// 从第二个元素开始, 在已排好序的子数组中寻找到该元素合适的位置并插入该位置

// 快排 O(n㏒n):
// 从数列中挑出一个元素, 称为"基准"(pivot)
// 重新排序数列, 所有比基准值小的元素摆放在基准前面, 所有比基准值大的元素摆在基准后面(相同的数可以到任何一边), 在这个分区结束之后, 该基准就处于数列的中间位置, 这个称为分区(partition)操作
// 递归地(recursively)把小于基准值元素的子数列和大于基准值元素的子数列排序

// 归并 O(n㏒n):
// 分治思想
