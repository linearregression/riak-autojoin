package main

import (
	"encoding/json"
	"log"
	"testing"
)

type TestExecuter struct{}

func (r *TestExecuter) Execute(name string, arg ...string) ([]byte, bool) {
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

func TestFindNodesToJoin(t *testing.T) {
	executer := new(JSONExecuter)
	executer.empty = false
	riak := Riak{executer: executer}
	exitCode := riak.find_nodes()
	if exitCode != true {
		t.Error("Should have found some nodes")
	}
}

func TestDoesNotFindNodesToJoin(t *testing.T) {
	executer := new(JSONExecuter)
	executer.empty = true

	riak := Riak{executer: executer}
	exitCode := riak.find_nodes()
	if exitCode != false {
		t.Error("Should have found some nodes")
	}
}
