package main

type executer interface {
	Execute(name string, arg ...string) ([]byte, bool)
}

type realexecuter struct{}

type riak struct {
	executer executer
}
