package main

type PointerReceiver struct {
	count int	
}

func (p *PointerReceiver) Increment() {
	p.count++
}
