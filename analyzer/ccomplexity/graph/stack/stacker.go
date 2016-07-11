//Package graph.stack provides an implementation of the
//Stack abstract-datastructure, basically providing
//the standard graph.stack operating primitives.
//
//Author: Mark Summerfield, found in the book:
//"Programming in Go, Creating Applications for the 21st Century"
//
//Author: Christian Bergum Bergersen, added comments to satisfy
//Golang code convention.
package stack

import "errors"

//Stack is an array data-structure with the properties of a LIFO queue.
type Stack []interface{}

//Len returns the size(length) of the graph.stack.
func (stack Stack) Len() int {
	return len(stack)
}

//Cap returns the capacity of the graph.stack.
func (stack Stack) Cap() int {
	return cap(stack)
}

//IsEmpty returns true if graph.stack is empty, false otherwise.
func (stack Stack) IsEmpty() bool {
	return len(stack) == 0
}

//Push puts x on the top of the graph.stack.
func (stack *Stack) Push(x interface{}) {
	*stack = append(*stack, x)
}

//Pop removes and returns the first(top) element of the
//graph.stack, or returns an error message if the graph.stack is empty.
func (stack *Stack) Pop() (interface{}, error) {
	theStack := *stack
	if len(theStack) == 0 {
		return nil, errors.New("can't Pop() an empty stack")
	}
	x := theStack[len(theStack)-1]
	*stack = theStack[:len(theStack)-1]
	return x, nil
}

//Top returns the first first(top) element in the
//graph.stack, or error message if the graph.stack is empty.
func (stack Stack) Top() (interface{}, error) {
	if len(stack) == 0 {
		return nil, errors.New("cant't Top() an empty stack")
	}
	return stack[len(stack)-1], nil
}
