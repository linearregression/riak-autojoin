package main

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
	Nodes   Node `json:"Node"`
	Service struct {
		Address string   `json:"Address"`
		ID      string   `json:"ID"`
		Port    int      `json:"Port"`
		Service string   `json:"Service"`
		Tags    []string `json:"Tags"`
	} `json:"Service"`
}

type Node struct {
	Address string `json:"Address"`
	Node    string `json:"Node"`
}

type Test []ConsulResponse

type Executer interface {
	Execute(name string, arg ...string) ([]byte, bool)
}

type RealExecuter struct{}

type Riak struct {
	executer Executer
}
