package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/consul/api"
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
var process_name = flag.String("process_name", "riak", "The process_name of the node we are connecting too")

var host = flag.String("host", "localhost", "The hostname for the consul")
var port = flag.String("port", "8500", "Port of the consul server")

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
	client, err := api.NewClient(&api.Config{
		Address: *host + ":" + *port,
	})
	if err != nil {
		log.Println("Unable to contact consul cluster")
	}
	health := client.Health()
	serviceEntries, _, err := health.Service(*service, *tag, true, &api.QueryOptions{})
	if err != nil {
		log.Println("Unable to talk to consul", err)
	}
	for _, v := range serviceEntries {
		log.Printf("Found node '%s' trying to join it\n", v.Node.Node)
		if r.join_riak(v.Node.Node) {
			return true
		}
	}
	return false
}

func (r *Riak) join_riak(nodehostname string) bool {
	out, success := r.executer.Execute("sudo", "-H", "-u riak", "riak-admin", "cluster", "join", *process_name+"@"+nodehostname)

	if !success {
		fmt.Println(string(out))
	}
	return success
}
