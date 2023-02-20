package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/video"
	auth "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/auth"
	feed "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/feed"
	"github.com/sosyz/mini_tiktok_feed/Feed/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	FeedService feed.FeedServiceClient
	AuthService auth.AuthServiceClient
	S3Service   *service.S3Service
	VideoHandle video.Video
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
				&Common{
					StatusCode: -1,
					StatusMsg:  "unauthorized",
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
				&Common{
					StatusCode: -1,
					StatusMsg:  "unauthorized",
				},
			)
			return
		}
		c.Set("reqUserID", res.UserId)
		c.Next()
	}
}

func InitService(feedConfig *conf.Server, authConfig *conf.Server, s3Config *conf.S3, volConfig *conf.Secret) error {
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

	S3Service, err = service.NewS3Service(s3Config.Region, s3Config.Endpoint, s3Config.SecretId, s3Config.SecretKey, s3Config.Bucket)
	if err != nil {
		return err
	}

	VideoHandle = video.NewVolcengineService(volConfig.SecretId, volConfig.SecretKey)
	return nil
}

func PlayerList(c *gin.Context) {
	latest_time := c.Query("latest_time")
	userId := c.GetInt64("reqUserID")

	respHandler := func(code int64, msg string, nextTime int64, videos []Video) *FeedPlayerListReply {
		ret := &FeedPlayerListReply{}
		if code != 0 {
			ret.StatusCode = code
			ret.StatusMsg = msg
			return ret
		}
		ret.StatusCode = 0
		ret.StatusMsg = msg
		ret.NextTime = nextTime
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
		UserId:      userId,
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

func PushList(c *gin.Context) {
	userId := c.GetInt64("reqUserID")
	respHandler := func(code int64, msg string, videos []Video) *FeedPushListReply {
		ret := &FeedPushListReply{}
		if code != 0 {
			ret.StatusCode = code
			ret.StatusMsg = msg
			return ret
		}
		ret.StatusCode = 0
		ret.StatusMsg = msg
		ret.VideoList = videos
		return ret
	}

	videos, err := FeedService.ListPublishVideos(context.Background(), &feed.ListPublishVideosRequest{
		UserId: userId,
	})
	if err != nil {
		c.JSON(500, respHandler(
			-1,
			"feed service error",
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
			retVideos,
		),
	)
}

func CreateWorks(c *gin.Context) {
	reqUserID := c.GetInt64("requsterID")
	title := c.PostForm("title")
	file, err := c.FormFile("file")

	respHandler := func(c *gin.Context, headCode int, code int64, msg string) {
		c.JSON(
			headCode,
			&Common{
				StatusCode: code,
				StatusMsg:  msg,
			},
		)
	}

	if err != nil {
		respHandler(c, 400, -1, "invalid file")
		return
	}

	fileHandle, err := file.Open()
	if err != nil {
		respHandler(c, 400, -1, "invalid file")
		return
	}

	defer fileHandle.Close()
	fileBytes, err := ioutil.ReadAll(fileHandle)

	file_path := fmt.Sprintf("%s/%s", md5.New().Sum([]byte(fmt.Sprintf("%d", reqUserID))), uuid.New().String())

	public_url, err := S3Service.SaveFile(file_path, fileBytes)
	if err != nil {
		respHandler(c, 500, -1, "save file error")
		return
	}

	cover, err := VideoHandle.GetVideoCover(public_url)
	if err != nil || len(cover) == 0 {
		respHandler(c, 500, -1, "get video cover error")
		return
	}

	cover_buff, err := io.ReadAll(cover[0])
	if err != nil {
		respHandler(c, 500, -1, "get video cover error")
		return
	}

	cover_path, err := S3Service.SaveFile(fmt.Sprintf("%s_cover", file_path), cover_buff)
	if err != nil {
		respHandler(c, 500, -1, "save cover error")
		return
	}

	res, err := FeedService.PublishVideo(c.Request.Context(), &feed.PublishVideoRequest{
		UserId:   reqUserID,
		Title:    title,
		PlayUrl:  public_url,
		CoverUrl: cover_path,
	})

	if err != nil {
		respHandler(c, 500, -1, "feed service error")
		return
	}

	if res.StatusCode != 0 {
		respHandler(c, 403, -1, res.StatusMsg)
		return
	}
	respHandler(c, 200, 0, "ok")
}

func FavoriteCreate(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		&Common{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
	)
}

func FavoriteList(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		&Common{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
	)
}
