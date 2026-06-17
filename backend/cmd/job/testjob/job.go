package testjob

import "log"

type job struct {
	msg string
}

func NewJob() *job {
	return &job{
		msg: "hello world",
	}
}

func (s *job) Run() {
	log.Println(s.msg)
}
