package notifier

// Notifier wraps a Writer.
type Notifier struct {
	Writer
}

// New returns a Notifier using Writer `w`.
func New(w Writer) *Notifier {
	return &Notifier{Writer: w}
}

// Writer is the interface that wraps the `WriteMessage` method.
type Writer interface {
	// WriteMessage takes in a message and sends it.
	WriteMessage(value []byte) error
}
