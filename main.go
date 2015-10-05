package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func (r *RealExecuter) Execute(name string, arg ...string) ([]byte, bool) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Getting output from command %s failed with '%s'\n", name, err)
	}
	return out, cmd.ProcessState.Success()
}

var service = flag.String("service", "riak", "The service name to listen for")
var tag = flag.String("tag", "", "The tag name to listen for")
var consul_path = flag.String("consul", "/usr/sbin/consul", "Path to the consul binary")
var riak_user = flag.String("riak-user", "riak", "The user name of the node we are connecting too")

var timeout = flag.Int("timeout", 5, "Timeout in seconds")
var timeout_iterations = flag.Int("timeout-iterations", 36, "Number of iterations to do the timeout")

func main() {
	riak := Riak{executer: new(RealExecuter)}
	flag.Parse()
	os.Exit(riak.main_loop())
}

func (r *Riak) main_loop() int {
	// Wait for 3 min 36*5 = 180

	for i := 0; i < *timeout_iterations; i++ {
		if r.find_nodes() {
			return 0
		}
		time.Sleep(time.Duration(*timeout) * time.Second)
	}
	return 1
}

func (r *Riak) find_nodes() bool {
	out, success := r.executer.Execute(*consul_path, "watch", "-service="+*service, "-tag="+*tag, "-type=service", "-passingonly=true")

	if success {
		var resp Test
		if err := json.NewDecoder(bytes.NewReader(out)).Decode(&resp); err != nil {
			log.Fatal("Unable to parse json from consul: ", err)
		}

		for _, k := range resp {
			if r.join_riak(k.Nodes.Node) {
				return true
			}
		}
	} else {
		log.Println("Consul watch didn't execute successfully.")
		log.Println(string(out))
	}
	return false
}

func (r *Riak) join_riak(nodehostname string) bool {
	out, success := r.executer.Execute("sudo", "-H", "-u riak", "riak-admin", "cluster", "join", *riak_user+"@"+nodehostname)

	if !success {
		fmt.Println(string(out))
	}
	return success
}
