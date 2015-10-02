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

type ConsulResponse struct {
	Checks []struct {
		CheckID     string `json:"CheckID"`
		Name        string `json:"Name"`
		Node        string `json:"Node"`
		Notes       string `json:"Notes"`
		Output      string `json:"Output"`
		ServiceID   string `json:"ServiceID"`
		ServiceName string `json:"ServiceName"`
		Status      string `json:"Status"`
	} `json:"Checks"`
	Node struct {
		Address string `json:"Address"`
		Node    string `json:"Node"`
	} `json:"Node"`
	Service struct {
		Address string   `json:"Address"`
		ID      string   `json:"ID"`
		Port    int      `json:"Port"`
		Service string   `json:"Service"`
		Tags    []string `json:"Tags"`
	} `json:"Service"`
}

type Test []ConsulResponse

var service = flag.String("service", "riak", "The service name to listen for")
var tag = flag.String("tag", "", "The tag name to listen for")

func main() {
	flag.Parse()
	main_loop()
}

func main_loop() {
	// Wait for 3 min 36*5 = 180
	for i := 0; i < 36; i++ {
		cmd := exec.Command("/usr/bin/consul", "watch", "-service="+*service, "-tag="+*tag, "-type=service", "-passingonly=true")

		out, _ := cmd.CombinedOutput()

		if cmd.ProcessState.Success() {
			var resp Test
			if err := json.NewDecoder(bytes.NewReader(out)).Decode(&resp); err != nil {
				log.Fatal(err)
			}

			for _, k := range resp {
				if join_riak(k.Node.Node) {
					os.Exit(0)
				}
			}
		} else {
			log.Println("Consul watch didn't execute successfully.")
			log.Println(string(out))
		}
		time.Sleep(5 * time.Second)
	}
	os.Exit(1)
}

func join_riak(nodename string) bool {
	cmd := exec.Command("sudo", "-H", "-u riak", "riak-admin", "cluster", "join", "riak@"+nodename)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(out))
	}
	fmt.Println(string(out))
	return cmd.ProcessState.Success()
}
