package inventory

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Generate ansible inventory for openshift
func Generate(nodes []string, dedicatedMasters []string, dedcatedEtcd []string, inventoryPath string) (string, error) {
	inventory, rErr := readInventory(inventoryPath)
	if rErr != nil {
		return "", rErr
	}

	setInventoryHosts(inventory, "nodes", nodes)

	var masters []string
	if len(dedicatedMasters) > 0 {
		masters = dedicatedMasters
	} else {
		masters = nodes
	}
	setInventoryHosts(inventory, "masters", masters)

	var etcd []string
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

func setInventoryHosts(inventory map[string]interface{}, nodeType string, hosts []string) {
	hostValue := make(map[string]interface{})
	for _, host := range hosts {
		hostValue[host] = true
	}

	ose := inventory["OSEv3"].(map[string]interface{})
	children := ose["children"].(map[string]interface{})
	nodeInfo := children[nodeType].(map[string]interface{})
	nodeInfo["hosts"] = hostValue
}
