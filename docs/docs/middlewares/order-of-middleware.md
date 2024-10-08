---
sidebar_position: 2
---

### Order of Middleware Execution

The middleware system implemented in the code follows a **stack-based execution model**. Middleware functions are executed in the order they are added, which means that the first middleware added is the first one to execute. However, due to the nature of how the middleware system works, the control flow unwinds in reverse order after processing, allowing each middleware to perform actions both **before** and **after** the next middleware in the stack.

#### Step-by-step Flow of Execution

Let's consider a sample code to show 3 middleware added in order: A, B, C which might look something like below.

```go
package main

import "context"

type Middleware[T any] struct {
	Name string
}

func (m Middleware[T]) Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error {
	println("Middleware " + m.Name + " -> before logic")
	err := next(ctx, item)
	println("Middleware " + m.Name + " -> before logic")

	return err
}
func main() {
   kp := v2.New[string]()
   middlewareA := Middleware[string]{Name: "A"}
   middlewareB := Middleware[string]{Name: "B"}
   middlewareC := Middleware[string]{Name: "C"}
   kp.AddMiddleware(middlewareA)
   kp.AddMiddleware(middlewareB)
   kp.AddMiddleware(middlewareC)
}
```

1. **Adding Middleware**:
    - Middleware is added to the processor using the `AddMiddleware` method. The order in which you add middleware determines the execution order.
    - For example, if you add middleware in this order: `MiddlewareA`, `MiddlewareB`, and `MiddlewareC`, they will be executed in the same sequence.

2. **Processing the Middleware Stack**:
    - When the `Process` method is called, the middleware functions are executed in the order they were added. The first middleware in the list is executed first.
    - During the execution of each middleware, the control is passed to the next middleware using the `next` function, which recursively processes each middleware until the end of the list is reached.

3. **Control Flow**:
    - Each middleware gets the opportunity to perform operations **before** and **after** the next middleware in the chain. This means that while the middleware chain is called in the order of addition, the response unwinds in reverse order.
    - For example, if `MiddlewareA`, `MiddlewareB`, and `MiddlewareC` were added to the processor, the order of execution would look like this:
      ```
      1. MiddlewareA -> before logic
      2. MiddlewareB -> before logic
      3. MiddlewareC -> before logic
      4. MiddlewareC -> after logic
      5. MiddlewareB -> after logic
      6. MiddlewareA -> after logic
      ```
    - This pattern allows each middleware to wrap around the behavior of the next one in the chain, making it possible to perform actions both before and after the core logic of the middleware stack.

#### Execution Example

Let's say we have three middleware components: `MiddlewareA`, `MiddlewareB`, and `MiddlewareC`. The following sequence of events will occur:

1. **Before Execution Phase**:
    - `MiddlewareA` is called first. It performs its "before" logic.
    - `MiddlewareA` calls the next middleware in the stack, which is `MiddlewareB`.
    - `MiddlewareB` performs its "before" logic and then calls `MiddlewareC`.
    - `MiddlewareC` performs its "before" logic. Since it's the last middleware, it reaches the core logic or the final result.

2. **After Execution Phase**:
    - After `MiddlewareC` completes its process, it returns control back to `MiddlewareB`.
    - `MiddlewareB` now performs its "after" logic and returns control to `MiddlewareA`.
    - Finally, `MiddlewareA` performs its "after" logic.

This pattern allows for a clean and structured way to handle pre-processing and post-processing in a middleware chain.

#### Key Characteristics

- **Forward Execution**: Middleware is executed in the order it is added.
- **Reverse Unwinding**: Once a middleware finishes its operation, control is returned to the previous middleware, allowing it to complete any after-execution logic.
- **Flexible Processing**: Middleware can modify the request, the response, or even handle errors in a centralized way.

This approach allows middleware authors to easily implement logic that needs to occur both **before** and **after** the main processing logic, creating a powerful mechanism for handling cross-cutting concerns like logging, authentication, error handling, and more.

### Usage
There are many middleware that make use of this pattern in KP.
- Consumer middleware retrieves an item from kafka and adds it in. (before logic)
- Deadletter middleware evaluates the result and determines if the message should go to deadletter instead (after logic)
- Backoff middleware delays execution and decreases or increases the interval for the next process (both before and after logic)
- Gracefulshutdown middleware doesn't have before or after logic, it simply stops the execution flow
- Measurement middleware starts a timer in before logic and measures the time taken in after logic.
- Retry middleware is same as deadletter (after logic)
- RetryCount middleware injects the retry count header onto context (before logic)
- Tracing middleware starts a span before and ends the span after (both before and after logic)
- You might have a custom middleware that might make use of any of the above pattern.

#### Summary

- **Order of execution** is **in the order middleware is added**.
- **Order of unwinding** is **in reverse order**.
- This enables middleware to do processing **before** and **after** the next middleware in the chain.

Understanding this order of execution is crucial for correctly implementing middleware logic, as it allows for both sequential processing and reverse unwinding for post-processing.
