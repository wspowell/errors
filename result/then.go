package result

func Then[T any, S any](self Result[T], f func(T) Result[S]) Result[S] {
	if self.IsOk() {
		return f(self.value)
	}

	return Err[S](self.err)
}

type When[T any, S any] Result[T]

func (self When[T, S]) Then(fn func(T) Result[S]) Result[S] {
	if Result[T](self).IsOk() {
		return fn(self.value)
	}

	return Err[S](self.err)
}
