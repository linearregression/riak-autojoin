package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/consul/api"
)

const (
	SUCCESS = 0
	FAILURE = 1
)

func (r *realexecuter) Execute(name string, arg ...string) ([]byte, bool) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Getting output from command %s failed with '%s'\n", name, err)
	}
	return out, cmd.ProcessState.Success()
}

var service = flag.String("service", "riak", "The service name to listen for")
var tag = flag.String("tag", "", "The tag name to listen for")

var process_name = flag.String("process_name", "riak", "The process_name of the node we are connecting too")

var host = flag.String("host", "localhost", "The hostname for the consul")
var port = flag.String("port", "8500", "Port of the consul server")

var timeout = flag.Int("timeout", 5, "Timeout in seconds")
var timeout_iterations = flag.Int("timeout-iterations", 36, "Number of iterations to do the timeout")

func main() {
	riak := riak{executer: new(realexecuter)}
	flag.Parse()
	os.Exit(riak.main_loop())
}

func (r *riak) main_loop() int {
	// Wait for 3 min 36*5 = 180

	for i := 0; i < *timeout_iterations; i++ {
		if r.join_nodes(r.discover_services()) {
			return SUCCESS
		}
		time.Sleep(time.Duration(*timeout) * time.Second)
	}
	return FAILURE
}

func (r *riak) discover_services() []*api.ServiceEntry {
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
	return serviceEntries
}

func (r *riak) join_nodes(serviceEntries []*api.ServiceEntry) bool {
	for _, v := range serviceEntries {
		log.Printf("Found node '%s' trying to join it\n", v.Node.Node)
		if r.join_riak(v.Node.Node) {
			return true
		}
	}
	return false
}

func (r *riak) join_riak(nodehostname string) bool {
	out, success := r.executer.Execute("sudo", "-H", "-u", "riak", "riak-admin", "cluster", "join", *process_name+"@"+nodehostname)

	if !success {
		log.Println(string(out))
	}

	out, success = r.executer.Execute("sudo", "-H", "-u", "riak", "riak-admin", "cluster", "plan")
	if !success {
		log.Println("Was not able to produce the riak plan due to: ", string(out))
	}
	log.Println(string(out))

	log.Println("Comitting the plan !")
	out, success = r.executer.Execute("sudo", "-H", "-u", "riak", "riak-admin", "cluster", "commit")
	if !success {
		log.Println("Was not able to commit the riak plan due to: ", string(out))
	}
	log.Println(string(out))
	return success
}
