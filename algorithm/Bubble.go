package main

import (
	"fmt"
)

func msort(l []int) {
	// 遍历列表中的每个元素
	for i := 0; i < len(l); i++ {
		// 开始实际排序列中的元素， 第一次排序就把最大的元素或最小元素的放到了最右边
		for j := 1; j < len(l)-i; j++ {
			// 如果当前元素比上一个元素小的话
			if l[j] < l[j-1] {
				// 就将它们互换
				l[j], l[j-1] = l[j-1], l[j]
			}
		}
	}
}

func main() {
	l := []int{3, 5, 1, 4, 2}
	msort(l)
	fmt.Println(l)
}
