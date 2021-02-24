package workers

type job struct {
	channel     chan interface{}
	subscribers int
}

func (j *job) BroadcastResult(result interface{}) {
	for j.subscribers > 0 {
		j.channel <- result
		j.subscribers--
	}
	close(j.channel)
}
