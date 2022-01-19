package webserver

import (
	"fmt"
	"time"
)

type Filter func(next Context)

type FilterBuilder func(next Filter) Filter

var _ FilterBuilder = TimeFilterBuilder

func TimeFilterBuilder(next Filter) Filter {
	return func(c Context) {
		start := time.Now().Nanosecond()
		next(c)
		end := time.Now().Nanosecond()
		fmt.Printf("webFilter takes %d nanoseconds \n", end-start)
	}
}
