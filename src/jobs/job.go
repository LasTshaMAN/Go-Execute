package jobs

import (
	"fmt"
	"reflect"
)

type Job struct {
	function interface{}
	args     []interface{}
}

func NewJob(function interface{}, args ...interface{}) *Job {
	if err := validateFunctionArgs(function, args); err != nil {
		panic(fmt.Errorf("function signature doesn't match passed arguments: %s", err))
	}
	return &Job{
		function: function,
		args:     args,
	}
}

func (job *Job) Execute() {
	var argsAsValues []reflect.Value
	for _, arg := range job.args {
		argsAsValues = append(argsAsValues, reflect.ValueOf(arg))
	}
	reflect.ValueOf(job.function).Call(argsAsValues)
}

func validateFunctionArgs(function interface{}, args []interface{}) error {
	funcType := reflect.ValueOf(function).Type()

	if funcType.NumIn() > len(args) {
		return fmt.Errorf("too few arguments were passed")
	}
	if funcType.NumIn() < len(args) {
		return fmt.Errorf("too many arguments were passed")
	}

	for i := 0; i < funcType.NumIn(); i++ {
		funcArgType := funcType.In(i)
		argType := reflect.ValueOf(args[i]).Type()
		if !argType.AssignableTo(funcArgType) {
			return fmt.Errorf("expected '%s', got '%s'", funcArgType.String(), argType.String())
		}
	}

	return nil
}
