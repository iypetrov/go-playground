package main

type ValueReceiver struct {
	count int	
}

func (v ValueReceiver) Increment() {
	v.count++
}
