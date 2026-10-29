package main

import "fmt"

func twoSum(nums []int, target int) []int {
	mp := map[int]int{}
	for i1, n := range nums {
		i2, ok := mp[n]
		if ok {
			return []int{i1, i2}
		}
		mp[target-n] = i1
	}
	return []int{}
}

func main() {
	fmt.Println(twoSum([]int{1, 4, 2, 3}, 4))             // [0 3]
	fmt.Println(twoSum([]int{2, 3, 4, 5}, 5))             // [0 1]
	fmt.Println(twoSum([]int{4, 0, 9, 5, 1}, 6))          // [3 4]
	fmt.Println(twoSum([]int{1, 2, 3, 4, 6}, 7))          // [0 4]
	fmt.Println(twoSum([]int{12, 5, 4, 3}, 8))            // [1 3]
	fmt.Println(twoSum([]int{0, 0, 0, 1, 0, 0, 0, 8}, 9)) // [3 7]
	fmt.Println(twoSum([]int{0, 0, 0, 1, 0, 0, 0, 0}, 9)) // []
}
