package notifier_test

type MockNotifier struct {
	Messages [][]byte
	Err      error
}

func (mn *MockNotifier) WriteMessage(value []byte) error {
	mn.Messages = append(mn.Messages, value)
	return mn.Err
}
