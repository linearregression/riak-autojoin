package main

import (
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/consul/api"
)

type TestExecuter struct {
	executeCalled bool
	hostname      string
}

func (r *TestExecuter) Execute(name string, arg ...string) ([]byte, bool) {
	r.executeCalled = true
	r.hostname = strings.Join(arg, " ")
	return []byte(""), true
}

type JSONExecuter struct {
	empty bool
}

func (r *JSONExecuter) Execute(name string, arg ...string) ([]byte, bool) {
	if !r.empty {
		data := make([]ConsulResponse, 0)
		data = append(data, ConsulResponse{Nodes: Node{Node: ""}})
		byte, err := json.Marshal(data)
		if err != nil {
			log.Println("->", err)
		}
		return byte, true
	} else {
		return []byte(""), false
	}
}

func TestDoesNotFindNodesToJoin(t *testing.T) {
	executer := new(JSONExecuter)
	executer.empty = true

	riak := Riak{executer: executer}

	entries := []*api.ServiceEntry{}
	joined := riak.join_nodes(entries)

	if joined == true {
		t.Error("Should not be able to join with an empty result set")
	}

}

func TestExecuteIsBeingCalledOnJoinRiak(t *testing.T) {
	executer := new(TestExecuter)

	riak := Riak{executer: executer}
	riak.join_riak("localhost")
	if !executer.executeCalled {
		t.Error("Execute should have been executed !")
	}
}

func TestJoinRiakWithCorrectHostname(t *testing.T) {
	executer := new(TestExecuter)

	riak := Riak{executer: executer}
	riak.join_riak("localhost")
	if !strings.HasSuffix(executer.hostname, "localhost") {
		t.Error("Execute should have been executed with localhost, but was ", executer.hostname)
	}
}
