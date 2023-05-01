package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

const ErrorLogColor = "\033[1;31m%s\033[0m"

var logger = log.New(os.Stderr, "[LOG] ", log.Ldate|log.Ltime)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ErrorHandler(err error) {
	logger.Printf(ErrorLogColor, err)
	if err, ok := errors.Cause(err).(stackTracer); ok {
		stacks := err.StackTrace()
		// skip golang runtime error frame(runtime.proc, goexit)
		fmt.Printf("%+v", stacks[0:len(stacks)-2])
	}
}
