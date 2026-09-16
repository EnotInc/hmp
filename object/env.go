package object

type Enviroment struct {
	consts map[string]bool
	store  map[string]Object
	outer  *Enviroment
}

func NewEnviroment() *Enviroment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Enviroment{consts: c, store: s, outer: nil}
}

func NewEnclosedEnviroment(outer *Enviroment) *Enviroment {
	env := NewEnviroment()
	env.outer = outer
	return env
}

func (e *Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Enviroment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Enviroment) Const(name string) {
	e.consts[name] = true
}
func (e *Enviroment) IsConst(name string) bool {
	_, isConst := e.consts[name]
	return isConst
}
