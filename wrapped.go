package errors

// Error string of the origin error.
// This does not include the internal code.
func (self *wrapped) Error() string {
	return self.err.Error()
}

// Unwrap to get the underlying error.
func (self *wrapped) Unwrap() error {
	return self.err
}
