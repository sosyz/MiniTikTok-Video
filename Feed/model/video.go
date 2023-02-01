package model

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/snowflake"
)

type Video struct {
	// 视频唯一标识
	ID int64
	// 文件名
	Name string
	// 视频标题
	Title string
	// 视频大小
	Size int64
	// 上传时间
	UploadTime int64
	// 作者
	Author int64
	// 状态
	Status uint8
}

var (
	db neo4j.DriverWithContext
	sf *snowflake.Worker
)

func InitVideo(uri, user, password, realm string, node int64) error {
	if db != nil {
		return nil
	}
	var err error
	db, err = neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(
			user,
			password,
			realm,
		),
	)
	if err != nil {
		return err
	}

	sf, err = snowflake.NewWorker(node)
	if err != nil {
		return err
	}
	return nil
}

func CloseVideo(ctx context.Context) error {
	if db == nil {
		return nil
	}
	err := db.Close(ctx)
	if err != nil {
		return err
	}
	db = nil
	return nil
}

func videoToMap(v *Video) map[string]interface{} {
	return map[string]interface{}{
		"id":          v.ID,
		"name":        v.Name,
		"title":       v.Title,
		"size":        v.Size,
		"upload_time": v.UploadTime,
		"author":      v.Author,
		"status":      v.Status,
	}
}

func (v *Video) Create(ctx context.Context) error {
	session := db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		id := sf.GetId()
		v.ID = id
		v.UploadTime = time.Now().Unix()
		_, err := tx.Run(ctx,
			`CREATE (v:Video {
				id: $id,
				name: $name,
				title: $title,
				size: $size,
				upload_time: $upload_time,
				author: $author,
				status: $status
			})`,
			videoToMap(v),
		)
		return nil, err
	})

	return err
}

func (v *Video) Update(ctx context.Context) error {
	session := db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			`MATCH (v:Video {id: $id}) 
			SET
				v.name = $name,
				v.title = $title, 
				v.size = $size, 
				v.upload_time = $upload_time, 
				v.author = $author, 
				v.status = $status `,
			videoToMap(v),
		)
		return nil, err
	})

	return err
}

func (v *Video) Delete(ctx context.Context) error {
	session := db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			`MATCH (v:Video {id: $id}) 
			DETACH DELETE v`,
			videoToMap(v),
		)
		return nil, err
	})

	return err
}

func (v *Video) Get(ctx context.Context) error {
	session := db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx,
			`MATCH (v:Video {id: $id}) 
			RETURN v`,
			videoToMap(v),
		)

		if err != nil {
			return nil, err
		}

		record, err := records.Single(ctx)
		if err != nil {
			return nil, err
		}

		return record.Values[0].(neo4j.Node), nil
	})

	if err != nil {
		return err
	}

	node := result.(neo4j.Node)
	v.ID = node.Props["id"].(int64)
	v.Name = node.Props["name"].(string)
	v.Title = node.Props["title"].(string)
	v.Size = node.Props["size"].(int64)
	v.UploadTime = node.Props["upload_time"].(int64)
	v.Author = node.Props["author"].(int64)
	v.Status = (uint8)(node.Props["status"].(int64))

	return nil
}

func List(ctx context.Context, lastest_time int64) ([]Video, error) {
	if lastest_time == 0 {
		lastest_time = time.Now().Unix()
	}
	session := db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx,
			`MATCH (v:Video) 
			WHERE v.upload_time < $lastest_time
			RETURN v`,
			map[string]interface{}{
				"lastest_time": lastest_time,
			},
		)

		if err != nil {
			return nil, err
		}

		var videos []Video
		for records.Next(ctx) {
			record := records.Record()
			node := record.Values[0].(neo4j.Node)
			v := Video{
				ID:         node.Props["id"].(int64),
				Name:       node.Props["name"].(string),
				Title:      node.Props["title"].(string),
				Size:       node.Props["size"].(int64),
				UploadTime: node.Props["upload_time"].(int64),
			}

			videos = append(videos, v)
		}

		return videos, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]Video), nil
}
