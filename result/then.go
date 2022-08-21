package result

func Then[T any, S any, E Optional](self Result[T, E], f func(T) Result[S, E]) Result[S, E] {
	if self.IsOk() {
		return f(self.value)
	}

	return Err[S](self.err)
}

type When[T any, S any, E Optional] Result[T, E]

func (self When[T, S, E]) Then(fn func(T) Result[S, E]) Result[S, E] {
	if Result[T, E](self).IsOk() {
		return fn(self.value)
	}

	return Err[S](self.err)
}
