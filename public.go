package respond

// Public allows types to control how they are publically viewed.
// The Public method will be called before the data is written.
type Public interface {
	Public() interface{}
}

const publicRecursiveLimit = 100

func public(o interface{}) interface{} {
	n := 0
	for {
		if p, ok := o.(Public); ok {
			o = p.Public()
		} else {
			break
		}
		n++
		if n >= publicRecursiveLimit {
			panic("respond: Public recursion limit reached")
		}
	}
	return o
}
