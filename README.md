# Auto joining Riak cluster through Consul

This is a small to for auto joining Riak nodes to a cluster, based on values from Consul.

```shell
Usage of ./riak-autojoin:
  -host string
    	The hostname for the consul (default "localhost")
  -port string
    	Port of the consul server (default "8500")
  -process_name string
    	The process_name of the node we are connecting too (default "riak")
  -service string
    	The service name to listen for (default "riak")
  -tag string
    	The tag name to listen for
  -timeout int
    	Timeout in seconds (default 5)
  -timeout-iterations int
    	Number of iterations to do the timeout (default 36)
```

The idea is that the tool will listen for events on the Consul cluster regarding the service / tag name combination. And
if there are events, return the node name, and try to join that cluster.
As soon as a successfull joining of the Riak cluster has been made, the program will exit.

Prebuilt linux binaries can be found in the version branches.
