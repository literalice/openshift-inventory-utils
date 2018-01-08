package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/literalice/openshift-inventory-utils/inventory"
)

func main() {
	masterArg := flag.String("masters", "", "Comma-sepalated domain list for master node")
	etcdArg := flag.String("etcd", "", "Comma-sepalated domain list for etcd node")
	nodesArg := flag.String("nodes", "", "Comma-sepalated domain list for nodes")
	inventoryPath := flag.String("inventory", "", "Inventory file on which the new inventory based")
	flag.Parse()

	nodes := parseNodeArg(*nodesArg)
	masters := parseNodeArg(*masterArg)
	etcd := parseNodeArg(*etcdArg)

	inventory, err := inventory.Generate(nodes, masters, etcd, *inventoryPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(inventory)
}

func parseNodeArg(arg string) []string {
	if arg == "" {
		return nil
	}
	return strings.Split(arg, ",")
}
