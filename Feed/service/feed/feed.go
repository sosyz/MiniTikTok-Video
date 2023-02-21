package feed

import (
	"context"
	"fmt"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/video"
	"github.com/sosyz/mini_tiktok_feed/Feed/model"
	pb_feed "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/feed"
	pb_user "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/user"
)

type FeedService struct {
	pb_feed.UnimplementedFeedServiceServer
	volc        video.Video
	vm          *model.VideoModel
	userService pb_user.UserServiceClient
}

func NewFeedService(
	neo4jConf *conf.Neo4j,
	secret *conf.Secret,
	userService pb_user.UserServiceClient,
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

func (v *FeedService) ListWatchVideos(ctx context.Context, req *pb_feed.ListWatchVideosRequest) (*pb_feed.ListWatchVideosResponse, error) {
	videos, err := v.vm.ListByLastesTime(ctx, req.LastestTime)
	if err != nil {
		return nil, err
	}

	var resp pb_feed.ListWatchVideosResponse
	nextTime := int64(0)
	for _, video := range videos {
		userInfos, err := v.userService.GetFullInfos(ctx, &pb_user.FollowCheckRequests{
			SelfId:    req.UserId,
			TargetIds: []int64{video.Author},
		})
		if err != nil {
			return nil, err
		}

		resp.Videos = append(resp.Videos, &pb_feed.Video{
			Id: video.ID,
			Author: &pb_feed.User{
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

func (v *FeedService) ListPublishVideos(ctx context.Context, req *pb_feed.ListPublishVideosRequest) (*pb_feed.ListPublishVideosResponse, error) {
	videos, err := v.vm.ListByAuthor(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	userInfos, err := v.userService.GetFullInfos(ctx, &pb_user.FollowCheckRequests{
		SelfId:    req.UserId,
		TargetIds: []int64{req.UserId},
	})
	if err != nil {
		return nil, err
	}

	userInfo := &pb_feed.User{
		Id:            userInfos.Infos[0].Id,
		Name:          userInfos.Infos[0].Name,
		FollowCount:   int64(userInfos.Infos[0].FollowCount),
		FollowerCount: int64(userInfos.Infos[0].FollowerCount),
	}

	var resp pb_feed.ListPublishVideosResponse
	for _, video := range videos {

		resp.Videos = append(resp.Videos, &pb_feed.Video{
			Id:       video.ID,
			Author:   userInfo,
			PlayUrl:  video.Name,
			CoverUrl: fmt.Sprintf("%s_cover", video.Name),
			Title:    video.Title,
		})
	}

	return &resp, nil
}

func (f *FeedService) FavoriteVideo(ctx context.Context, req *pb_feed.FavoriteVideoRequest) (*pb_feed.FavoriteVideoResponse, error) {
	ret := &pb_feed.FavoriteVideoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return ret, nil
}

func (f *FeedService) ListFavoriteVideos(ctx context.Context, req *pb_feed.ListFavoriteVideosRequest) (*pb_feed.ListFavoriteVideosResponse, error) {
	ret := &pb_feed.ListFavoriteVideosResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Videos:     []*pb_feed.Video{},
	}
	return ret, nil
}

func (f *FeedService) PublishVideo(ctx context.Context, req *pb_feed.PublishVideoRequest) (*pb_feed.PublishVideoResponse, error) {
	video := &model.VideoInfo{
		Name:   req.PlayUrl,
		Title:  req.Title,
		Author: req.UserId,
	}
	err := f.vm.Create(context.Background(), video)
	if err != nil {
		return nil, err
	}
	userInfo, err := f.userService.GetFullInfos(ctx, &pb_user.FollowCheckRequests{
		SelfId:    req.UserId,
		TargetIds: []int64{req.UserId},
	})
	if err != nil {
		return nil, err
	}
	return &pb_feed.PublishVideoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Video: &pb_feed.Video{
			Id: video.ID,
			Author: &pb_feed.User{
				Id:            req.UserId,
				Name:          userInfo.Infos[0].Name,
				FollowCount:   int64(userInfo.Infos[0].FollowCount),
				FollowerCount: int64(userInfo.Infos[0].FollowerCount),
				IsFollow:      userInfo.Infos[0].IsFollow,
			},
			PlayUrl:       video.Name,
			Title:         video.Title,
			CoverUrl:      fmt.Sprintf("%s_cover", video.Name),
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
		},
	}, nil
}
