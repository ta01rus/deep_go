package main

import "fmt"

/*
	Я изменил сигнатуру так как, я не понимаю как в конечном счете обрабатывать конструктор
	которые не являются функциями. Структуры или пойнторы только усложнят код Resolve придется применять reflect.

	А если конструктор функция с несколькими параметрами то как их выбирать из интерфейса, или к примеру функция возвращает кортеж.

	Считаю что все это можно решить через замыкание при вызове RegisterType.

*/

type ConstructorFN func() interface{}

type Container struct {
	constructors map[string]ConstructorFN
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]ConstructorFN, 4),
	}
}

func (c *Container) RegisterType(name string, fn ConstructorFN) {
	c.constructors[name] = fn
}

func (c *Container) Resolve(name string) (interface{}, error) {

	fn, ok := c.constructors[name]
	if !ok {
		return nil, fmt.Errorf("constructor for %s not find", name)
	}

	return fn(), nil
}
