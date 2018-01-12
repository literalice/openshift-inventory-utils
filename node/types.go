package node

// Node of OpenShift
type Node struct {
	Host string
	Vars map[string]interface{}
}
