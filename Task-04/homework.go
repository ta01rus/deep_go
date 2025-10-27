package main

type OrderedMap struct {
	size int
	root *node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = newNode(key, value, nil)
		m.size = 1
		return

	}
	n := m.root.add(key, value)
	if n != nil {
		m.size++
	}
}

func (m *OrderedMap) Erase(key int) {
	n := m.root.del(key)
	if n != nil {
		m.size--
	}
}

func (m *OrderedMap) Contains(key int) bool {
	f := m.root.find(key)
	if f != nil {
		return true
	}
	return false // need to implement
}

func (m *OrderedMap) Size() int {

	return m.size // need to implement
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	for n := range m.root.Next(m.size) {
		action(n.key, n.value)
	}
}
