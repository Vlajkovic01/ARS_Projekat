package configstore

import (
	"ARS_Projekat/tracer"
	"context"
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

func (cs *ConfigStore) CreateConfig(ctx context.Context, config *Config) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateConfig")
	defer span.Finish()

	kv := cs.cli.KV()

	childCtx := tracer.ContextWithSpan(ctx, span)

	sid, rid := generateConfigKey(childCtx, config.Version)
	config.ID = rid

	data, err := json.Marshal(config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return config, nil
}

func (cs *ConfigStore) FindConfig(ctx context.Context, id string, ver string) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "FindConfig")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()
	key := constructConfigKey(childCtx, id, ver)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		tracer.LogError(span, err)
		return nil, errors.New("That item does not exist!")
	}

	config := &Config{}
	err = json.Unmarshal(data.Value, config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return config, nil
}

func (cs *ConfigStore) FindConfVersions(ctx context.Context, id string) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "FindConfigVersions")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()

	key := constructConfigIdKey(childCtx, id)
	data, _, err := kv.List(key, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	var configs []*Config

	for _, pair := range data {
		config := &Config{}
		err := json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}

		configs = append(configs, config)
	}

	return configs, nil
}

func (cs *ConfigStore) UpdateConfigVersion(ctx context.Context, config *Config) (*Config, error) {
	span := tracer.StartSpanFromContext(context.Background(), "UpdateConfigVersion")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()

	data, err := json.Marshal(config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	_, err = cs.FindConfig(childCtx, config.ID, config.Version)

	//Does exist
	if err == nil {
		tracer.LogError(span, err)
		return nil, errors.New("Given config version already exists! ")
	}

	c := &api.KVPair{Key: constructConfigKey(childCtx, config.ID, config.Version), Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}
	return config, nil

}

func (cs *ConfigStore) DeleteConfig(ctx context.Context, id, ver string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteConfig")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()
	_, err := kv.Delete(constructConfigKey(childCtx, id, ver), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return map[string]string{"Deleted config": id + ver}, nil
}

func (cs *ConfigStore) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateGroup")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()

	sid, rid := generateGroupKey(childCtx, group.Version)
	group.ID = rid

	data, err := json.Marshal(group)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	g := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	err = cs.CreateLabels(childCtx, group.Configs, group.ID, group.Version)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return group, nil
}

func (cs *ConfigStore) FindGroup(ctx context.Context, id string, ver string) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "FindGroup")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()
	key := constructGroupKey(childCtx, id, ver)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		tracer.LogError(span, err)
		return nil, errors.New("That item does not exist!")
	}

	group := &Group{}
	err = json.Unmarshal(data.Value, group)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return group, nil
}

func (cs *ConfigStore) UpdateGroupVersion(ctx context.Context, group *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateGroupVersion")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	_, err = cs.FindGroup(childCtx, group.ID, group.Version)

	//Does exist
	if err == nil {
		tracer.LogError(span, err)
		return nil, errors.New("Given group version already exists! ")
	}

	c := &api.KVPair{Key: constructGroupKey(childCtx, group.ID, group.Version), Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	err = cs.CreateLabels(childCtx, group.Configs, group.ID, group.Version)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return group, nil

}

func (cs *ConfigStore) DeleteGroup(ctx context.Context, id, ver string) error {
	span := tracer.StartSpanFromContext(ctx, "DeleteGroup")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()

	_, err := kv.DeleteTree(constructGroupKey(childCtx, id, ver), nil)

	return err
}

func (cs *ConfigStore) CreateLabels(ctx context.Context, configs []map[string]string, id, ver string) error {
	span := tracer.StartSpanFromContext(ctx, "CreateLabels")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()
	if keys, _, err := kv.Get(constructGroupKey(childCtx, id, ver), nil); err != nil || keys == nil {
		tracer.LogError(span, err)
		return errors.New("Group doesn't exists")
	}

	for _, config := range configs {
		cid := constructGroupLabel(childCtx, id, ver, uuid.New().String(), config)
		cdata, err := json.Marshal(config)

		log.Default().Printf("adding new config: %q. under key %q", config, cdata)
		if err != nil {
			tracer.LogError(span, err)
			return err
		}

		c := &api.KVPair{Key: cid, Value: cdata}
		_, err = kv.Put(c, nil)
		if err != nil {
			tracer.LogError(span, err)
			return err
		}
	}
	return nil
}

func (cs *ConfigStore) FindLabels(ctx context.Context, id, ver, kvpairs string) ([]map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "FindLabels")
	defer span.Finish()

	kv := cs.cli.KV()
	labelkey := fmt.Sprintf(group, id, ver, kvpairs) + "/"
	keys, _, err := kv.List(labelkey, nil)
	if err != nil {
		tracer.LogError(span, err)
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

func (cs *ConfigStore) AddLabelsToGroup(ctx context.Context, configs []map[string]string, id, ver string) ([]map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "AddLabelsToGroup")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()
	gr, err := cs.FindGroup(childCtx, id, ver)
	if err != nil || gr == nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for _, config := range configs {
		log.Default().Printf("%q", configs)
		gr.Configs = append(gr.Configs, config)
	}

	data, err := json.Marshal(gr)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	sid := constructGroupKey(childCtx, id, ver)

	g := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	err = cs.CreateLabels(childCtx, configs, gr.ID, gr.Version)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return gr.Configs, nil
}

func (cs *ConfigStore) SaveRequestId(ctx context.Context) string {
	span := tracer.StartSpanFromContext(ctx, "SaveRequestId")
	defer span.Finish()

	childCtx := tracer.ContextWithSpan(ctx, span)

	kv := cs.cli.KV()

	reqId := generateRequestId(childCtx)

	i := &api.KVPair{Key: reqId, Value: nil}

	_, err := kv.Put(i, nil)

	if err != nil {
		tracer.LogError(span, err)
		return "error"
	}

	return reqId
}

func (cs *ConfigStore) FindRequestId(ctx context.Context, requestId string) bool {
	span := tracer.StartSpanFromContext(ctx, "FindRequestId")
	defer span.Finish()

	kv := cs.cli.KV()

	key, _, err := kv.Get(requestId, nil)

	fmt.Println(key)

	if err != nil || key == nil {
		tracer.LogError(span, err)
		return false
	}

	return true
}
