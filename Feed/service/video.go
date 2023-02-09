package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/video"
	"github.com/sosyz/mini_tiktok_feed/Feed/model"
	feed "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/feed"
	user "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/user"
)

type VideoService struct {
	feed.UnimplementedFeedServiceServer
	s3          S3Service
	volc        video.Video
	vm          *model.VideoModel
	userService user.UserServiceClient
}

func NewVideoService(
	s3Conf conf.S3,
	neo4jConf conf.Neo4j,
	secret conf.Secret,
	userService user.UserServiceClient,
	node int64) (*VideoService, error) {
	s3, err := NewS3Service(
		s3Conf.Region,
		s3Conf.Endpoint,
		s3Conf.Secret.Id,
		s3Conf.Secret.Key,
		s3Conf.Bucket,
	)
	if err != nil {
		return nil, err
	}

	n4, err := model.NewVideoModel(
		fmt.Sprintf("bolt://%s:%d", neo4jConf.Host, neo4jConf.Port),
		neo4jConf.User,
		neo4jConf.Password,
		neo4jConf.Realm,
		node,
	)
	if err != nil {
		return nil, err
	}

	return &VideoService{
		s3:          *s3,
		volc:        video.NewVolcengineService(secret.Id, secret.Key),
		vm:          n4,
		userService: userService,
	}, nil
}

func (v *VideoService) CreateVideo(file io.Reader, autherId uint32, title string) (*model.VideoInfo, error) {
	file_path := fmt.Sprintf("%s/%s", md5.New().Sum([]byte(fmt.Sprintf("%d", autherId))), uuid.New().String())
	file_bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	public_url := fmt.Sprintf("https://%s.%s/%s", v.s3.bucket, v.s3.endpoint, file_path)
	err = v.s3.SaveFile(file_path, file_bytes)
	if err != nil {
		return nil, err
	}

	cover, err := v.volc.GetVideoCover(public_url)
	if err != nil {
		return nil, err
	}

	if len(cover) == 0 {
		return nil, fmt.Errorf("no cover")
	}

	cover_buff, err := io.ReadAll(cover[0])
	if err != nil {
		return nil, err
	}

	err = v.s3.SaveFile(fmt.Sprintf("%s/cover", file_path), cover_buff)
	if err != nil {
		return nil, err
	}

	// 视频信息存入库
	video := &model.VideoInfo{
		Name:   file_path,
		Title:  title,
		Size:   int64(len(file_bytes)),
		Author: autherId,
	}
	err = v.vm.Create(context.Background(), video)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (v *VideoService) GetVideo(videoId int64) (*model.VideoInfo, error) {
	return v.vm.Get(context.Background(), videoId)
}

func (v *VideoService) ListWatchVideos(ctx context.Context, req *feed.ListWatchVideosRequest) (*feed.ListWatchVideosResponse, error) {
	videos, err := v.vm.ListByLastesTime(ctx, req.LastestTime)
	if err != nil {
		return nil, err
	}

	var resp feed.ListWatchVideosResponse
	for _, video := range videos {
		userInfo, err := v.userService.GetInfo(ctx, &user.UserInfoRequest{
			TargetId: video.Author,
			SelfId:   req.UserId,
		})
		if err != nil {
			return nil, err
		}

		resp.Videos = append(resp.Videos, &feed.Video{
			Id: video.ID,
			Author: &feed.User{
				Id:            userInfo.UserId,
				Name:          userInfo.Username,
				FollowCount:   int64(userInfo.FollowCount),
				FollowerCount: int64(userInfo.FollowerCount),
			},
			Title:    video.Title,
			PlayUrl:  fmt.Sprintf("https://%s.%s/%s", v.s3.bucket, v.s3.endpoint, video.Name),
			CoverUrl: fmt.Sprintf("https://%s.%s/%s/cover", v.s3.bucket, v.s3.endpoint, video.Name),
		})
	}

	return &resp, nil
}

func (v *VideoService) ListVideosByAuthor(ctx context.Context, req *feed.ListPublishVideosRequest) (*feed.ListPublishVideosResponse, error) {
	videos, err := v.vm.ListByAuthor(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	userRes, err := v.userService.GetInfo(ctx, &user.UserInfoRequest{
		TargetId: req.UserId,
		SelfId:   req.UserId,
	})
	if err != nil {
		return nil, err
	}

	userInfo := &feed.User{
		Id:            userRes.UserId,
		Name:          userRes.Username,
		FollowCount:   int64(userRes.FollowCount),
		FollowerCount: int64(userRes.FollowerCount),
	}

	var resp feed.ListPublishVideosResponse
	for _, video := range videos {

		resp.Videos = append(resp.Videos, &feed.Video{
			Id:       video.ID,
			Author:   userInfo,
			PlayUrl:  fmt.Sprintf("https://%s.%s/%s", v.s3.bucket, v.s3.endpoint, video.Name),
			CoverUrl: fmt.Sprintf("https://%s.%s/%s/cover", v.s3.bucket, v.s3.endpoint, video.Name),
			Title:    video.Title,
		})
	}

	return &resp, nil
}
