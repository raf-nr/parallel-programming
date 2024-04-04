package stacks

type Stack[T any] interface {
	Push(T) error
	Pop() (T, error)
	Peek() (T, error)
	Len() (int, error)
	IsEmpty() bool
}

const (
	EmptyStackError          = "Stack is already empty."
	StackNilPointerError     = "The consistentStack pointer is nil."
	UnsuccessfulPrimitivePop = "Failed to remove element: trying to find a complementary operation."
)
