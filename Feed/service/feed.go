package service

import (
	"context"
	"fmt"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/video"
	"github.com/sosyz/mini_tiktok_feed/Feed/model"
	feed "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/feed"
	user "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/user"
)

type FeedService struct {
	feed.UnimplementedFeedServiceServer
	volc        video.Video
	vm          *model.VideoModel
	userService user.UserServiceClient
}

func NewFeedService(
	neo4jConf *conf.Neo4j,
	secret *conf.Secret,
	userService user.UserServiceClient,
	node int64) (*FeedService, error) {

	n4, err := model.NewVideoModel(
		fmt.Sprintf("bolt+s://%s:%d", neo4jConf.Host, neo4jConf.Port),
		neo4jConf.User,
		neo4jConf.Password,
		neo4jConf.Realm,
		node,
	)
	if err != nil {
		return nil, err
	}

	return &FeedService{
		volc:        video.NewVolcengineService(secret.SecretId, secret.SecretKey),
		vm:          n4,
		userService: userService,
	}, nil
}

func (v *FeedService) GetVideo(videoId int64) (*model.VideoInfo, error) {
	return v.vm.Get(context.Background(), videoId)
}

func (v *FeedService) ListWatchVideos(ctx context.Context, req *feed.ListWatchVideosRequest) (*feed.ListWatchVideosResponse, error) {
	videos, err := v.vm.ListByLastesTime(ctx, req.LastestTime)
	if err != nil {
		return nil, err
	}

	var resp feed.ListWatchVideosResponse
	nextTime := int64(0)
	for _, video := range videos {
		userInfos, err := v.userService.GetFullInfos(ctx, &user.FollowCheckRequests{
			SelfId:    req.UserId,
			TargetIds: []int64{video.Author},
		})
		if err != nil {
			return nil, err
		}

		resp.Videos = append(resp.Videos, &feed.Video{
			Id: video.ID,
			Author: &feed.User{
				Id:            userInfos.Infos[0].Id,
				Name:          userInfos.Infos[0].Name,
				FollowCount:   int64(userInfos.Infos[0].FollowCount),
				FollowerCount: int64(userInfos.Infos[0].FollowerCount),
			},
			Title:    video.Title,
			PlayUrl:  video.Name,
			CoverUrl: fmt.Sprintf("%s/cover", video.Name),
		})

		if video.UploadTime < nextTime {
			nextTime = video.UploadTime
		}
	}
	resp.NextTime = nextTime
	return &resp, nil
}

func (v *FeedService) ListPublishVideos(ctx context.Context, req *feed.ListPublishVideosRequest) (*feed.ListPublishVideosResponse, error) {
	videos, err := v.vm.ListByAuthor(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	userInfos, err := v.userService.GetFullInfos(ctx, &user.FollowCheckRequests{
		SelfId:    req.UserId,
		TargetIds: []int64{req.UserId},
	})
	if err != nil {
		return nil, err
	}

	userInfo := &feed.User{
		Id:            userInfos.Infos[0].Id,
		Name:          userInfos.Infos[0].Name,
		FollowCount:   int64(userInfos.Infos[0].FollowCount),
		FollowerCount: int64(userInfos.Infos[0].FollowerCount),
	}

	var resp feed.ListPublishVideosResponse
	for _, video := range videos {

		resp.Videos = append(resp.Videos, &feed.Video{
			Id:       video.ID,
			Author:   userInfo,
			PlayUrl:  video.Name,
			CoverUrl: fmt.Sprintf("%s_cover", video.Name),
			Title:    video.Title,
		})
	}

	return &resp, nil
}

func (f *FeedService) FavoriteVideo(ctx context.Context, req *feed.FavoriteVideoRequest) (*feed.FavoriteVideoResponse, error) {
	return nil, nil
}

func (f *FeedService) ListFavoriteVideos(ctx context.Context, req *feed.ListFavoriteVideosRequest) (*feed.FavoriteVideoResponse, error) {
	ret := &feed.FavoriteVideoResponse{}
	return ret, nil
}

func (f *FeedService) PublishVideo(ctx context.Context, req *feed.PublishVideoRequest) (*feed.FavoriteVideoResponse, error) {
	video := &model.VideoInfo{
		Name:   req.PlayUrl,
		Title:  req.Title,
		Author: req.UserId,
	}
	err := f.vm.Create(context.Background(), video)
	if err != nil {
		return nil, err
	}
	return &feed.FavoriteVideoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}
