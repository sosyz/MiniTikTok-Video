package service

import (
	"os"
	"testing"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	"github.com/spf13/viper"
)

var (
	s3Conf  conf.S3
	cfgFile = "../../config-dev.yaml"
)

func TestUpload(t *testing.T) {
	var cfg struct {
		S3 conf.S3
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
	s3Conf = cfg.S3
	t.Log(s3Conf)
	// Create a new S3 client
	s, err := NewS3Service(
		s3Conf.Region,
		s3Conf.Endpoint,
		s3Conf.SecretId,
		s3Conf.SecretKey,
		s3Conf.Bucket,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Open the file for use
	file, err := os.Open("IMG_1596.MP4")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Get the file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	// t.Logf("file content: %s\n", buffer)
	path, err := s.SaveFile(file.Name(), buffer)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Successfully uploaded %q to %q\n", file.Name(), path)
}
