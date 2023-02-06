package result

// import "context"

// func Then[T any, S any, E Optional](ctx context.Context, self Result[T, E], f func(context.Context, T) Result[S, E]) Result[S, E] {
// 	if self.IsOk() {
// 		return f(ctx, self.value)
// 	}

// 	return Err[S](self.err)
// }

// type When[T any, S any, E Optional] Result[T, E]

// func (self When[T, S, E]) Then(ctx context.Context, fn func(context.Context, T) Result[S, E]) Result[S, E] {
// 	if Result[T, E](self).IsOk() {
// 		return fn(ctx, self.value)
// 	}

// 	return Err[S](self.err)
// }
