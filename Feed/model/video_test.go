package model

import (
	"context"
	"fmt"
	"testing"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	"github.com/spf13/viper"
)

var (
	neo4jConf conf.Neo4j
	cfgFile         = "../../config-dev.yaml"
	id        int64 = 7108107930963968
)

func TestCreate(t *testing.T) {
	t.Log("test create")
	var cfg struct {
		Neo4j conf.Neo4j
	}
	v := viper.New()
	v.SetConfigFile(cfgFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	t.Log(v.AllSettings())
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	neo4jConf = cfg.Neo4j
	t.Log(neo4jConf)
	ctx := context.Background()
	vm, err := NewVideoModel(
		fmt.Sprintf("bolt://%s:%d", neo4jConf.Host, neo4jConf.Port),
		neo4jConf.User,
		neo4jConf.Password,
		neo4jConf.Realm,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := vm.Close(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	video := &VideoInfo{
		Name:   "test",
		Title:  "test",
		Size:   100,
		Author: 1,
		Status: 0,
	}
	err = vm.Create(ctx, video)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("create video %v", video)
}

func TestGet(t *testing.T) {
	t.Log("test get")
	var cfg struct {
		Neo4j conf.Neo4j
	}
	v := viper.New()
	v.SetConfigFile(cfgFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	neo4jConf = cfg.Neo4j
	t.Log(neo4jConf)
	ctx := context.Background()
	vm, err := NewVideoModel(
		fmt.Sprintf("bolt://%s:%d", neo4jConf.Host, neo4jConf.Port),
		neo4jConf.User,
		neo4jConf.Password,
		neo4jConf.Realm,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := vm.Close(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	video, err := vm.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(video)
}

func TestListByLastesTime(t *testing.T) {
	t.Log("test list")
	var cfg struct {
		Neo4j conf.Neo4j
	}
	v := viper.New()
	v.SetConfigFile(cfgFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	neo4jConf = cfg.Neo4j
	t.Log(neo4jConf)
	ctx := context.Background()
	vm, err := NewVideoModel(
		fmt.Sprintf("bolt://%s:%d", neo4jConf.Host, neo4jConf.Port),
		neo4jConf.User,
		neo4jConf.Password,
		neo4jConf.Realm,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := vm.Close(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()
	videos, err := vm.ListByLastesTime(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(videos)
}

func TestUpdate(t *testing.T) {
	t.Log("test update")
	var cfg struct {
		Neo4j conf.Neo4j
	}
	v := viper.New()
	v.SetConfigFile(cfgFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	neo4jConf = cfg.Neo4j
	t.Log(neo4jConf)
	ctx := context.Background()
	vm, err := NewVideoModel(
		fmt.Sprintf("bolt://%s:%d", neo4jConf.Host, neo4jConf.Port),
		neo4jConf.User,
		neo4jConf.Password,
		neo4jConf.Realm,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := vm.Close(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	video, err := vm.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(video)
	video.Title = "test update"
	err = vm.Update(ctx, video)
	if err != nil {
		t.Fatal(err)
	}
	video, err = vm.Get(ctx, video.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(video)
}

func TestDelete(t *testing.T) {
	var cfg struct {
		Neo4j conf.Neo4j
	}
	t.Log("test delete")
	v := viper.New()
	v.SetConfigFile(cfgFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	neo4jConf = cfg.Neo4j
	t.Log(neo4jConf)
	ctx := context.Background()
	vm, err := NewVideoModel(
		fmt.Sprintf("bolt://%s:%d", neo4jConf.Host, neo4jConf.Port),
		neo4jConf.User,
		neo4jConf.Password,
		neo4jConf.Realm,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := vm.Close(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	video, err := vm.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(video)
	err = vm.Delete(ctx, video)
	if err != nil {
		t.Logf("err: %v", err)
		t.Fatal(err)
	}
	_, err = vm.Get(ctx, id)
	if err != nil {
		t.Logf("err: %v", err)
		t.Log("delete success")
		return
	}
	t.Fatal("delete failed")
}
