#### Concurrency in Go: Goroutines, Channels, and Select

This guide explains the three core building blocks of concurrency in Go: goroutines, channels, and the select statement. Together they let you write programs that do many things at once in a clear and safe way.

#### Concurrency vs. Parallelism

Concurrency is about structuring a program as several independent tasks that can make progress without waiting on each other. Parallelism is about actually running those tasks at the same moment on multiple CPU cores. Go gives you the tools to express concurrency, and its runtime decides how to spread the work across available cores. A concurrent program can run on a single core, and it can also run in parallel when more cores are available.

Go's approach is summed up in a well-known saying: do not communicate by sharing memory; instead, share memory by communicating. Rather than having many threads read and write the same variables protected by locks, Go encourages you to pass data between tasks through channels.

### Goroutines

A goroutine is a lightweight unit of execution managed by the Go runtime, not by the operating system. You start one by placing the go keyword in front of a function call. The function then runs independently while the calling code carries on immediately.

Goroutines are very cheap. They start with only a few kilobytes of stack space, which grows and shrinks as needed, so a program can comfortably run thousands or even millions of them. The Go scheduler multiplexes these goroutines onto a small number of operating system threads, which is why creating them costs so little compared to traditional threads.

There are a few important things to remember:

When the main function returns, the whole program exits, even if other goroutines are still running. Goroutines do not keep a program alive.
The order in which goroutines run is not guaranteed. Two goroutines started one after the other may run in either order, or interleave.
Pausing the main function with a sleep to wait for goroutines is unreliable. The proper tools for waiting are a WaitGroup from the sync package or channels.

A WaitGroup is essentially a counter. You increase it before launching each goroutine, each goroutine decreases it when it finishes, and the waiting code blocks until the counter returns to zero. It is the standard way to say "wait until all of these tasks are done."

#### Channels

A channel is a typed conduit that lets goroutines send values to one another. If a channel carries integers, only integers can pass through it. One goroutine sends a value into the channel, and another receives it from the other end. Because the channel handles the handoff, you avoid many of the problems that come from sharing variables directly.

Unbuffered channels

An unbuffered channel has no storage at all. A send operation blocks until another goroutine is ready to receive, and a receive blocks until another goroutine is ready to send. This means the sender and receiver meet at the same moment, so unbuffered channels act as both a communication mechanism and a synchronization point. When a receive completes, you know for certain that the sender has reached its send.

Buffered channels

A buffered channel has a fixed capacity. Sending only blocks when the buffer is full, and receiving only blocks when the buffer is empty. This lets the sender and receiver run at slightly different speeds without constantly waiting for each other. Buffers are useful for smoothing out bursts of work, but they also make the timing between goroutines less strict, so they should be chosen deliberately rather than added just to make a deadlock disappear.

Closing a channel

A sender can close a channel to signal that no more values will ever be sent. Receivers can then keep reading until the channel is closed, which is commonly done with a range loop that ends automatically once the channel is closed and drained. A receiver can also check whether a channel is closed by using the two-value form of receive, which returns the value along with a boolean indicating whether it came from an open channel.

The rules around closing are important:

Only the sender should close a channel, never the receiver.
Sending on a closed channel causes a panic.
Closing a channel that is already closed also causes a panic.
Receiving from a closed channel returns immediately with the zero value of the type.
You do not need to close every channel. Closing is only necessary when a receiver needs to know that no more data is coming.
Directional channels

A function can declare that it will only send on a channel or only receive from one. This is expressed in the parameter type. Doing so documents your intent, and the compiler will reject code that tries to use the channel in the wrong direction. It is a small feature that prevents a surprising number of bugs.

#### Select

The select statement lets a goroutine wait on several channel operations at the same time. It looks like a switch statement, but each case is a send or a receive on a channel. Select blocks until one of its cases is able to proceed, then runs that case.

Some key behaviors:

If several cases are ready at the same moment, Go picks one of them at random. This prevents any single channel from being starved.
If a default case is included, select never blocks. When no other case is ready, the default case runs immediately.
A select with no cases blocks forever.
Common select patterns

Timeouts. By including a case that receives from a timer channel, such as the one returned by time.After, you can give up waiting for a slow operation after a set amount of time.

Non-blocking operations. Adding a default case turns a send or receive into something that either happens right now or is skipped. This is useful for checking whether a message is available without pausing the goroutine.

Cancellation. A goroutine can listen on a dedicated "done" channel alongside its normal work. When the channel is closed, every goroutine listening on it is notified at once, because receiving from a closed channel succeeds immediately. This is how long-running goroutines are told to stop cleanly. In real programs, the standard library's context package builds on this same idea and is generally preferred for cancellation and timeouts.

Multiplexing. When several producers send on different channels, select lets one goroutine collect results from all of them as they arrive, without caring which one finishes first.

#### How They Work Together

These three features are designed to be combined. A typical pattern is the worker pool: several goroutines act as workers, all reading jobs from one shared channel and writing results to another. The main goroutine feeds jobs into the first channel and then closes it, which tells the workers there is no more work. A WaitGroup tracks when all workers have finished so the results channel can be closed safely. Finally, the collecting code uses select to read results while also watching a timeout or cancellation signal, so the program never hangs forever.

This structure scales naturally. Adding more workers requires no change to the logic, and the channels take care of distributing work and coordinating the handoff.

Common Pitfalls

Deadlock. This happens when every goroutine is blocked waiting on something that will never happen, such as sending on an unbuffered channel with no receiver. Go's runtime detects the simplest cases and crashes with a clear error message, but deadlocks involving only some goroutines may go unnoticed while part of the program silently stalls.

Goroutine leaks. A goroutine that is blocked forever on a channel is never cleaned up and keeps consuming memory. Every goroutine you start should have a clear way to finish, such as a closed channel or a cancellation signal.

Data races. If two goroutines access the same variable at the same time and at least one of them writes to it, the result is unpredictable. The fix is to share data through channels, protect it with a mutex, or use atomic operations. Go includes a built-in race detector, enabled with the -race flag when running or testing, which is well worth using while you learn.

Loop variable capture. In versions of Go before 1.22, a loop variable was shared across all iterations, so goroutines started inside a loop could all see the final value. Go 1.22 changed this so each iteration gets its own copy, but if you work with older code or older versions, pass the value into the goroutine as an argument.

Closing channels incorrectly. Closing from the receiver side, or closing the same channel twice, causes panics. Keep closing the responsibility of a single, clearly identified sender.

Nil channels. Sending to or receiving from a nil channel blocks forever. This is occasionally used deliberately inside select to disable a case, but when it happens by accident it is a hard bug to track down.