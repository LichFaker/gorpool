package gorpool

type Queue []interface{}

func (q *Queue) Push(obj interface{}) {
	*q = append(*q, obj)
}

func (q *Queue) Pop() interface{} {
	if q.IsEmpty() {
		return nil
	}
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

func (q *Queue) GetTop() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return (*q)[0]
}

func (q *Queue) IsEmpty() bool {
	return len(*q) <= 0
}
