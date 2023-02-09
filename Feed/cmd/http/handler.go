package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	auth "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/auth"
	feed "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/feed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	FeedService feed.FeedServiceClient
	AuthService auth.AuthServiceClient
)

func AuthdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			// 例外路径
			if c.Request.URL.Path == "/douyin/feed/" {
				c.Next()
				return
			}
			c.JSON(
				401,
				&FeedPlayerListReply{
					StatusCode: -1,
					StatusMsg:  &[]string{"unauthorized"}[0],
				},
			)
			return
		}
		res, err := AuthService.Auth(c.Request.Context(), &auth.AuthRequest{
			Token: token,
		})
		if err != nil {
			c.JSON(
				401,
				&FeedPlayerListReply{
					StatusCode: -1,
					StatusMsg:  &[]string{"unauthorized"}[0],
				},
			)
			return
		}
		c.Set("requsterID", res.UserId)
		c.Next()
	}
}

func InitService(feedConfig *conf.Server, authConfig *conf.Server) error {
	feedConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", feedConfig.Host, feedConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return err
	}
	FeedService = feed.NewFeedServiceClient(feedConn)

	authConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", authConfig.Host, authConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return err
	}

	AuthService = auth.NewAuthServiceClient(authConn)
	return nil
}

func FeedPlayerList(c *gin.Context) {
	latest_time := c.Query("latest_time")
	userId := c.GetUint("user_id")

	respHandler := func(code int64, msg string, nextTime int64, videos []Video) *FeedPlayerListReply {
		ret := &FeedPlayerListReply{}
		if code != 0 {
			ret.StatusCode = code
			ret.StatusMsg = &msg
			return ret
		}
		ret.StatusCode = 0
		ret.StatusMsg = &msg
		ret.NextTime = &nextTime
		ret.VideoList = videos
		return ret
	}

	latestTime, err := strconv.ParseInt(latest_time, 10, 64)
	if err != nil {
		c.JSON(
			400,
			respHandler(
				-1,
				"invalid latest_time",
				0,
				nil,
			),
		)
		return
	}
	res, err := FeedService.ListWatchVideos(context.Background(), &feed.ListWatchVideosRequest{
		LastestTime: latestTime,
		UserId:      uint32(userId),
	})
	if err != nil {
		c.JSON(500, respHandler(
			-1,
			"feed service error",
			0,
			nil),
		)
		return
	}
	videos := make([]Video, len(res.Videos))
	for i, v := range res.Videos {
		videos[i] = Video{
			Author: User{
				ID:            int64(v.Author.Id),
				Name:          v.Author.Name,
				IsFollow:      v.Author.IsFollow,
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
			},
			CommentCount:  v.CommentCount,
			CoverURL:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			ID:            v.Id,
			IsFavorite:    v.IsFavorite,
			PlayURL:       v.PlayUrl,
			Title:         v.Title,
		}
	}
	c.AsciiJSON(
		200,
		respHandler(
			0,
			"ok",
			res.NextTime,
			videos,
		),
	)
}

func FeedPushList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	respHandler := func(code int64, msg string, nextTime int64, videos []Video) *FeedPushListReply {
		ret := &FeedPushListReply{}
		if code != 0 {
			ret.StatusCode = code
			ret.StatusMsg = &msg
			return ret
		}
		ret.StatusCode = 0
		ret.StatusMsg = &msg
		ret.NextTime = &nextTime
		ret.VideoList = videos
		return ret
	}

	userId, err := strconv.ParseUint(user_id, 10, 32)
	if err != nil {
		c.JSON(
			400,
			respHandler(
				-1,
				"invalid user_id",
				0,
				nil,
			),
		)
		return
	}

	res, err := AuthService.Auth(c.Request.Context(), &auth.AuthRequest{
		Token: token,
	})

	if err != nil {
		c.JSON(
			401,
			respHandler(
				-1,
				"unauthorized",
				0,
				nil,
			),
		)
		return
	}

	if res.UserId != uint32(userId) {
		c.JSON(
			401,
			respHandler(
				-1,
				"unauthorized",
				0,
				nil,
			),
		)
		return
	}

	videos, err := FeedService.ListPublishVideos(context.Background(), &feed.ListPublishVideosRequest{
		UserId: uint32(userId),
	})
	if err != nil {
		c.JSON(500, respHandler(
			-1,
			"feed service error",
			0,
			nil),
		)
		return
	}

	retVideos := make([]Video, len(videos.Videos))
	for i, v := range videos.Videos {
		retVideos[i] = Video{
			Author: User{
				ID:            int64(v.Author.Id),
				Name:          v.Author.Name,
				IsFollow:      v.Author.IsFollow,
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
			},
			CommentCount:  v.CommentCount,
			CoverURL:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			ID:            v.Id,
			IsFavorite:    v.IsFavorite,
			PlayURL:       v.PlayUrl,
			Title:         v.Title,
		}
	}
	c.AsciiJSON(
		200,
		respHandler(
			0,
			"ok",
			0,
			retVideos,
		),
	)
}

func FeedCreate(c *gin.Context) {
	fmt.Println("abc")
}

func FavoriteCreate(c *gin.Context) {

}

func FavoriteList(c *gin.Context) {

}
