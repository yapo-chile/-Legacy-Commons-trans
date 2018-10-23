package domain

// Fibonacci is the datatype for numbers that belong to the Fibonacci Series
type Fibonacci int

// FibonacciPair is the pair of the latest pair of known Fibonacci Numbers
type FibonacciPair struct {
	// Indexes
	IA, IB int
	// Values
	A, B Fibonacci
}

// FibonacciRepository defines a backing storage for Fibonacci Numbers
type FibonacciRepository interface {
	// Get should retrieve the Nth Fibonacci if available
	Get(nth int) (Fibonacci, error)
	// Save must store the Nth Fibonacci only if (N-1) and (N-2) are known
	Save(nth int, x Fibonacci) error
	// LatestPair should retrieve the latest know pair of Fibonacci
	LatestPair() FibonacciPair
}

// Next produces the Fibonacci Number and the index that comes right after f
func (f FibonacciPair) Next() (int, Fibonacci) {
	return f.IB + 1, f.A + f.B
}
