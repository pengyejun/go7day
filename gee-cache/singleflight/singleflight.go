package singleflight

import "sync"

type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	m map[string]*call
}

func New() *Group {
	return &Group{m:make(map[string]*call)}
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {

	if c, ok := g.m[key]; ok {
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c

	c.val, c.err = fn()
	c.wg.Done()

	delete(g.m, key)
	return c.val, c.err
}

