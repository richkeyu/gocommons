package assert

import "fmt"

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: "+msg, v...))
	}
}

func Assert(condition bool, msg string, v ...interface{}) {
	_assert(condition, msg, v)
}
