package configstore

import (
	"fmt"
	"github.com/google/uuid"
	"sort"
)

const (
	allConfigs = "config"
	configId   = "config/%s"
	config     = "config/%s/%s"

	allGroups      = "group"
	groupId        = "group/%s"
	groupVer       = "group/%s/%s"
	group          = "group/%s/%s/%s"
	groupWithLabel = "group/%s/%s/%s/%d"
)

func generateConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(config, id, ver), id
}

func constructConfigKey(id string, ver string) string {
	return fmt.Sprintf(config, id, ver)
}

func constructConfigIdKey(id string) string {
	return fmt.Sprintf(configId, id)
}

func generateGroupKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(groupVer, id, ver), id
}

func constructGroupKey(id string, ver string) string {
	return fmt.Sprintf(groupVer, id, ver)
}

func constructGroupIdKey(id string) string {
	return fmt.Sprintf(groupId, id)
}

func constructGroupLabel(id, ver string, index int, config map[string]string) string {
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var kvpairs string
	for k := range keys {

		kvpairs = kvpairs + fmt.Sprintf("%s=%s", keys[k], config[keys[k]]+"&")
	}
	kvpairs = kvpairs[:len(kvpairs)-1]
	return fmt.Sprintf(groupWithLabel, id, ver, kvpairs, index)
}
