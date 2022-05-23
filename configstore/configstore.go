package configstore

import (
	"encoding/json"
	"errors"
	"fmt"
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

	return map[string]string{"Deleted config": id + ver}, nil
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

	log.Default().Println(sid, kv)

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	g := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		return nil, err
	}

	cids := make([]string, 0)
	for i, config := range group.Configs {
		cid := constructGroupLabel(rid, group.Version, i, config)
		cids = append(cids, cid)
		log.Default().Println(cid)
		cdata, err := json.Marshal(config)
		if err != nil {
			return nil, err
		}

		c := &api.KVPair{Key: cid, Value: cdata}
		_, err = kv.Put(c, nil)
		if err != nil {
			return nil, err
		}
	}

	for _, k := range cids {
		val, _, err := kv.List(k, nil)
		if err != nil {
			return nil, err
		}
		log.Default().Println(val)
	}

	log.Default().Println(cs.FindConfVersions(sid))

	return group, nil
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
	return group, nil

}

func (cs *ConfigStore) FindGroupVersions(id string) ([]*Group, error) {
	kv := cs.cli.KV()

	key := constructGroupIdKey(id)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, err
	}

	var groups []*Group

	for _, pair := range data {
		group := &Group{}
		err := json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}
	return groups, nil

}

func (cs *ConfigStore) DeleteConfigGroup(id, ver string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(constructGroupKey(id, ver), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted group": id + ver}, nil
}
