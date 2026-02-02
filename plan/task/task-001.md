## Task 001: Create the Virtual Machine node structure.


**Description**
Setup the node and connection struct

``` golang
type Connection struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}
```

``` golang
type Node struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
	PID     int32  `json:"pid"`
	Active  bool
	conn    net.Conn
}
```

**Commits**

feat(feature-001): add node registration (1985ec17640cca6d32a3cd9d9730419383370993)