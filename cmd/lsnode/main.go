package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/literalice/openshift-inventory-utils/node"
)

// Cmd is used for handling cmd interface
func main() {
	clusterArg := flag.String("cluster", "", "Cluster name used in the tag: kubernetes.io/cluster/xxx")
	roleArg := flag.String("role", "", "master / etcd / node")
	roleTagArg := flag.String("role-tag", "Role", "Tag name for specifying node types")
	flag.Parse()

	nodes, err := node.List(*clusterArg, *roleArg, *roleTagArg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(strings.Join(nodes, ","))
}
