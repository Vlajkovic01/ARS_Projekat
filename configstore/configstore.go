package configstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
)

type ConfigStore struct {
	cli *api.Client
}

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}

func (cs *ConfigStore) FindConf(id string, ver string) (*Config, error) {
	kv := cs.cli.KV()
	key := constructConfigKey(id, ver)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("That item does not exist!")
	}

	config := &Config{}
	err = json.Unmarshal(data.Value, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (cs *ConfigStore) DeleteConfig(id, ver string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(constructConfigKey(id, ver), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (cs *ConfigStore) FindConfVersions(id string) ([]*Config, error) {
	kv := cs.cli.KV()

	key := constructConfigIdKey(id)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, err
	}

	var configs []*Config

	for _, pair := range data {
		config := &Config{}
		err := json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}

	return configs, nil
}

func (cs *ConfigStore) CreateConfig(config *Config) (*Config, error) {
	kv := cs.cli.KV()

	sid, rid := generateConfigKey(config.Version)
	config.ID = rid

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (cs *ConfigStore) UpdateConfigVersion(config *Config) (*Config, error) {
	kv := cs.cli.KV()

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	_, err = cs.FindConf(config.ID, config.Version)

	//Does exist
	if err == nil {
		return nil, errors.New("Given config version already exists! ")
	}

	c := &api.KVPair{Key: constructConfigKey(config.ID, config.Version), Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}
	return config, nil

}

func (cs *ConfigStore) CreateGroup(group *Group) (*Group, error) {
	kv := cs.cli.KV()

	sid, rid := generateGroupKey(group.Version)
	group.ID = rid

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	g := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		return nil, err
	}

	err = cs.CreateLabels(group.Configs, group.ID, group.Version)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (cs *ConfigStore) CreateLabels(configs []map[string]string, id, ver string) error {
	kv := cs.cli.KV()
	if keys, _, err := kv.Get(constructGroupKey(id, ver), nil); err != nil || keys == nil {
		return errors.New("Group doesn't exists")
	}

	for _, config := range configs {
		cid := constructGroupLabel(id, ver, uuid.New().String(), config)
		cdata, err := json.Marshal(config)

		log.Default().Printf("adding new config: %q. under key %q", config, cdata)
		if err != nil {
			return err
		}

		c := &api.KVPair{Key: cid, Value: cdata}
		_, err = kv.Put(c, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs *ConfigStore) AddLabelsToGroup(configs []map[string]string, id, ver string) ([]map[string]string, error) {
	kv := cs.cli.KV()
	gr, err := cs.FindGroup(id, ver)
	if err != nil || gr == nil {
		return nil, err
	}

	for _, config := range configs {
		log.Default().Printf("%q", configs)
		gr.Configs = append(gr.Configs, config)
	}

	data, err := json.Marshal(gr)
	if err != nil {
		return nil, err
	}

	sid := constructGroupKey(id, ver)

	g := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		return nil, err
	}

	err = cs.CreateLabels(configs, gr.ID, gr.Version)
	if err != nil {
		return nil, err
	}

	return gr.Configs, nil
}

func (cs *ConfigStore) FindLabels(id, ver, kvpairs string) ([]map[string]string, error) {
	kv := cs.cli.KV()
	labelkey := fmt.Sprintf(group, id, ver, kvpairs) + "/"
	keys, _, err := kv.List(labelkey, nil)
	if err != nil {
		return nil, err
	}

	configs := make([]map[string]string, len(keys))
	for i, k := range keys {
		var config map[string]string
		json.Unmarshal(k.Value, &config)
		log.Default().Printf("%q", config)
		configs[i] = config
	}

	return configs, nil
}

func (cs *ConfigStore) FindGroup(id string, ver string) (*Group, error) {
	kv := cs.cli.KV()
	key := constructGroupKey(id, ver)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("That item does not exist!")
	}

	group := &Group{}
	err = json.Unmarshal(data.Value, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (cs *ConfigStore) UpdateGroupVersion(group *Group) (*Group, error) {
	kv := cs.cli.KV()

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	_, err = cs.FindGroup(group.ID, group.Version)

	//Does exist
	if err == nil {
		return nil, errors.New("Given group version already exists! ")
	}

	c := &api.KVPair{Key: constructGroupKey(group.ID, group.Version), Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	err = cs.CreateLabels(group.Configs, group.ID, group.Version)
	if err != nil {
		return nil, err
	}

	return group, nil

}

func (cs *ConfigStore) DeleteGroup(id, ver string) error {
	kv := cs.cli.KV()

	_, err := kv.DeleteTree(constructGroupKey(id, ver), nil)

	return err
}
