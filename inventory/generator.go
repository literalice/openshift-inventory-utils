package inventory

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/literalice/openshift-inventory-utils/node"
)

// Generate ansible inventory for openshift
func Generate(nodes []*node.Node, dedicatedMasters []*node.Node, dedcatedEtcd []*node.Node, inventoryPath string) (string, error) {
	inventory, rErr := readInventory(inventoryPath)
	if rErr != nil {
		return "", rErr
	}

	setInventoryHosts(inventory, "nodes", nodes)

	var masters []*node.Node
	if len(dedicatedMasters) > 0 {
		masters = dedicatedMasters
	} else {
		masters = nodes
	}
	setInventoryHosts(inventory, "masters", masters)

	var etcd []*node.Node
	if len(dedcatedEtcd) > 0 {
		etcd = dedcatedEtcd
	} else {
		etcd = masters
	}
	setInventoryHosts(inventory, "etcd", etcd)

	data, mErr := yaml.Marshal(inventory)
	if mErr != nil {
		return "", mErr
	}

	return string(data), nil
}

func readInventory(path string) (inventory map[string]interface{}, err error) {
	inventory = make(map[string]interface{})
	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	yaml.Unmarshal(data, &inventory)
	return
}

func setInventoryHosts(inventory map[string]interface{}, nodeType string, nodes []*node.Node) {
	nodeValue := make(map[string]interface{})
	for _, node := range nodes {
		if nodeType == "nodes" && len(node.Vars) > 0 {
			nodeValue[node.Host] = node.Vars
		} else {
			nodeValue[node.Host] = ""
		}
	}

	ose := inventory["OSEv3"].(map[string]interface{})
	children := ose["children"].(map[string]interface{})
	nodeInfo := children[nodeType].(map[string]interface{})
	nodeInfo["hosts"] = nodeValue
}
