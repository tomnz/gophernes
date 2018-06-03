package cpu

const maxOps = 1024

type opQueue struct {
	ops        [maxOps]func()
	start, end int
}

func (o *opQueue) empty() bool {
	return o.start == o.end
}

func (o *opQueue) push(fn func()) {
	o.ops[o.end] = fn
	o.end++
	o.end %= maxOps
	if o.empty() {
		panic("overflowed op queue")
	}
}

func (o *opQueue) pop() func() {
	if o.empty() {
		panic("tried to pop from empty op queue")
	}
	fn := o.ops[o.start]
	o.start++
	o.start %= maxOps
	return fn
}
