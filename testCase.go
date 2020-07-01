package simply

// TestCase represents a generic test objet which can be batched in a slice of tests for simply.Run()
type TestCase struct {
	// PreHook can execute a dynamic modifier before comparison
	PreHook func(TestCase)

	Message string

	Target interface{}
	Expected interface{}
}