package main

import (
	"fmt"
	"strings"
)

func IsNotNill(e error) bool {
	return e != nil
}

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {

	if len(e.errors) == 1 {
		return fmt.Sprintf("1 error occured:\n\t* %s\n", e.errors[0])
	}

	points := make([]string, len(e.errors))
	for i, err := range e.errors {
		points[i] = fmt.Sprintf("* %s", err)
	}
	// ;)
	return fmt.Sprintf(
		"%d errors occured:\n\t%s\n",
		len(e.errors), strings.Join(points, "\t"))
}

func Append(err error, errs ...error) *MultiError {
	var n int

	if err != nil {
		n++
	}

	for _, e := range errs {
		if e != nil {
			n++
		}
	}

	result := &MultiError{
		errors: make([]error, 0, n),
	}

	add := func(e error) {
		switch t := e.(type) {
		case *MultiError:
			result.errors = append(result.errors, t.errors...)
		default:
			result.errors = append(result.errors, e)

		}
	}

	if IsNotNill(err) {
		add(err)
	}

	for _, e := range errs {
		if IsNotNill(e) {
			add(e)
		}
	}

	return result
}
