package gorpool

type Task struct {
	f    func(...interface{}) bool // do the task in this func
	args []interface{}               // the arguments for task func
}

func CreateTask(f func(...interface{}) bool, args ...interface{}) Task {
	return Task{
		f:    f,
		args: args,
	}
}

func (t *Task) Do() {
	t.f(t.args...)
}
