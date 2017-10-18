package main

import (
	"fmt"
)

func xsort(l []int) {
	// 遍历列表中的每个元素
	for i := 0; i < len(l); i++ {
		// 设定一个最小值
		var min int = i
		// 开始实际排序列中的元素, 选出最小的元素
		for j := i + 1; j < len(l); j++ {
			// 如果当前的最小元素大于其他元素
			if l[min] > l[j] {
				// 那就其他元素就是最小元素
				min = j
			}
		}
		// 将最小的元素往左排
		l[i], l[min] = l[min], l[i]
	}
}

func main() {
	l := []int{3, 5, 1, 4, 2}
	xsort(l)
	fmt.Println(l)
}
