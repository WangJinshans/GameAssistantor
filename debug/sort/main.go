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
