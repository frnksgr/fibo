package fibo

// Loop calculation of fibonacci number n
func Loop(n uint) uint64 {
	if n <= 1 {
		return uint64(n)
	}

	var n2, n1 uint64 = 0, 1

	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1
}

// Recursive calculation of fibonacci number n (NOTE: exponential grows)
func Recursive(n uint) uint64 {
	var result uint64
	switch n {
	case 0:
		result = 0
	case 1:
		result = 1
	default:
		result = Recursive(n-1) + Recursive(n-2)
	}
	return result
}

// RecursiveSequential calculation of fibonacci number n
func RecursiveSequential(n uint) uint64 {
	_, result := recursiveSequential(n)
	return result
}

func recursiveSequential(n uint) (uint64, uint64) {
	var left, right uint64
	switch n {
	case 0:
		left, right = 0, 0
	case 1:
		left, right = 0, 1
	default:
		left, right = recursiveSequential(n - 1)
		left, right = right, left+right
	}
	return left, right
}
