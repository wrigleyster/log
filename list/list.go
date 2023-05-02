package list

func El[T any](arr []T, i int) T {
	if i < 0 {
		i += len(arr)
	}
	return arr[i]
}

func Sl[T any](arr []T, i ...int) []T {
	if len(i) > 2 {
		panic("more than two indices not supported")
	}
	if len(i) == 0 {
		return arr
	}
	a := i[0]
	if a < 0 {
		a += len(arr)
	}
	if len(i) == 1 {
		return arr[a:]
	}
	b := i[1]
	if b < 0 {
		b += len(arr)
	}
	if a > b {
		return nil
	}
	return arr[a:b]
}

func Reverse[C any](arr []C) {
	for i, j := 0, len(arr)-1; i < j; i++ {
		arr[i], arr[j] = arr[j], arr[i]
		j--
	}
}
