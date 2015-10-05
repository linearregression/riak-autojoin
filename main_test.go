package main

import (
	"strings"
	"testing"

	"github.com/hashicorp/consul/api"
)

type testexecuter struct {
	executeCalled bool
	hostname      string
}

func (r *testexecuter) Execute(name string, arg ...string) ([]byte, bool) {
	r.executeCalled = true
	r.hostname = strings.Join(arg, " ")
	return []byte(""), true
}

func TestDoesNotFindNodesToJoin(t *testing.T) {
	executer := new(testexecuter)

	riak := riak{executer: executer}

	entries := []*api.ServiceEntry{}
	joined := riak.join_nodes(entries)

	if joined == true {
		t.Error("Should not be able to join with an empty result set")
	}

}

func TestExecuteIsBeingCalledOnJoinRiak(t *testing.T) {
	executer := new(testexecuter)

	riak := riak{executer: executer}
	riak.join_riak("localhost")
	if !executer.executeCalled {
		t.Error("Execute should have been executed !")
	}
}

func TestJoinRiakWithCorrectHostname(t *testing.T) {
	executer := new(testexecuter)

	riak := riak{executer: executer}
	riak.join_riak("localhost")
	if !strings.HasSuffix(executer.hostname, "localhost") {
		t.Error("Execute should have been executed with localhost, but was ", executer.hostname)
	}
}
