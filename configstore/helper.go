package configstore

import (
	"ARS_Projekat/tracer"
	"context"
	"fmt"
	"github.com/google/uuid"
	"sort"
)

const (
	configId = "config/%s"
	config   = "config/%s/%s"

	groupVer       = "group/%s/%s"
	group          = "group/%s/%s/%s"
	groupWithLabel = "group/%s/%s/%s/%s"

	requestId = "request/%s"
)

func generateConfigKey(ctx context.Context, ver string) (string, string) {
	span := tracer.StartSpanFromContext(ctx, "generateConfigKey")
	defer span.Finish()

	id := uuid.New().String()
	return fmt.Sprintf(config, id, ver), id
}

func constructConfigKey(ctx context.Context, id string, ver string) string {
	span := tracer.StartSpanFromContext(ctx, "constructConfigKey")
	defer span.Finish()

	return fmt.Sprintf(config, id, ver)
}

func constructConfigIdKey(ctx context.Context, id string) string {
	span := tracer.StartSpanFromContext(ctx, "constructConfigIdKey")
	defer span.Finish()

	return fmt.Sprintf(configId, id)
}

func generateGroupKey(ctx context.Context, ver string) (string, string) {
	span := tracer.StartSpanFromContext(ctx, "generateGroupKey")
	defer span.Finish()

	id := uuid.New().String()
	return fmt.Sprintf(groupVer, id, ver), id
}

func constructGroupKey(ctx context.Context, id string, ver string) string {
	span := tracer.StartSpanFromContext(ctx, "constructGroupKey")
	defer span.Finish()

	return fmt.Sprintf(groupVer, id, ver)
}
func constructGroupLabel(ctx context.Context, id, ver, index string, config map[string]string) string {
	span := tracer.StartSpanFromContext(ctx, "constructGroupLabel")
	defer span.Finish()

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

func generateRequestId(ctx context.Context) string {
	span := tracer.StartSpanFromContext(ctx, "generateRequestId")
	defer span.Finish()

	rid := uuid.New().String()
	return rid
}
