package main

func Binary[V Number](arr []V, target V) int {
	var low, high int = 0, len(arr) - 1
	for low <= high {
		var mid int = low + (high-low)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}
