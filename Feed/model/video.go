package model

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/snowflake"
)

type VideoInfo struct {
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
	Author uint32
	// 状态
	Status uint8
}

var ()

type VideoModel struct {
	db neo4j.DriverWithContext
	sf *snowflake.Worker
}

func NewVideoModel(uri, user, password, realm string, node int64) (*VideoModel, error) {
	db, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(
			user,
			password,
			realm,
		),
	)
	if err != nil {
		return nil, err
	}

	sf, err := snowflake.NewWorker(node)
	if err != nil {
		return nil, err
	}
	return &VideoModel{db, sf}, nil
}

func (v *VideoModel) Close(ctx context.Context) error {
	if v.db == nil {
		return nil
	}
	err := v.db.Close(ctx)
	if err != nil {
		return err
	}
	v.db = nil
	return nil
}

func videoToMap(v *VideoInfo) map[string]interface{} {
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

func (v *VideoModel) Create(ctx context.Context, video *VideoInfo) error {
	session := v.db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		id := v.sf.GetId()
		video.ID = id
		video.UploadTime = time.Now().Unix()
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
			videoToMap(video),
		)
		return nil, err
	})

	return err
}

func (v *VideoModel) Update(ctx context.Context, video *VideoInfo) error {
	session := v.db.NewSession(ctx, neo4j.SessionConfig{})
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
			videoToMap(video),
		)
		return nil, err
	})

	return err
}

func (v *VideoModel) Delete(ctx context.Context, video *VideoInfo) error {
	session := v.db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			`MATCH (v:Video {id: $id}) 
			DETACH DELETE v`,
			videoToMap(video),
		)
		return nil, err
	})

	return err
}

func (v *VideoModel) Get(ctx context.Context, id int64) (*VideoInfo, error) {
	session := v.db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	video := &VideoInfo{
		ID: id,
	}
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx,
			`MATCH (v:Video {id: $id}) 
			RETURN v`,
			videoToMap(video),
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
		return nil, err
	}

	node := result.(neo4j.Node)

	video.ID = node.Props["id"].(int64)
	video.Name = node.Props["name"].(string)
	video.Title = node.Props["title"].(string)
	video.Size = node.Props["size"].(int64)
	video.UploadTime = node.Props["upload_time"].(int64)
	video.Author = node.Props["author"].(uint32)
	video.Status = (uint8)(node.Props["status"].(int64))

	return video, nil
}

func (v *VideoModel) ListByLastesTime(ctx context.Context, lastest_time int64) ([]*VideoInfo, error) {
	if lastest_time == 0 {
		lastest_time = time.Now().Unix()
	}
	session := v.db.NewSession(ctx, neo4j.SessionConfig{})
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

		var videos []*VideoInfo
		for records.Next(ctx) {
			record := records.Record()
			node := record.Values[0].(neo4j.Node)
			v := &VideoInfo{
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

	return result.([]*VideoInfo), nil
}

func (v *VideoModel) ListByAuthor(ctx context.Context, author uint32) ([]*VideoInfo, error) {
	session := v.db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx,
			`MATCH (v:Video) 
			WHERE v.author = $author
			RETURN v`,
			map[string]interface{}{
				"author": author,
			},
		)

		if err != nil {
			return nil, err
		}

		var videos []*VideoInfo
		for records.Next(ctx) {
			record := records.Record()
			node := record.Values[0].(neo4j.Node)
			v := &VideoInfo{
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

	return result.([]*VideoInfo), nil
}
