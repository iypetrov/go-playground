![diagram](diagram.png)

allows you to process a large number of jobs concurrently, making the most of your system’s resources

core idea is to create a fixed number of worker goroutines that continuously pull jobs from a shared channel

this approach helps limit the number of concurrently running goroutines, preventing resource exhaustion while still allowing for parallel processing

great at:
- batch processing
- handling multiple api requests simultaneously
- load balancing, as tasks are distributed evenly among the available workers
