package wanjies3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	endpoint        string
	region          string
	accessKey       string
	accessKeySecret string
	s3Cli           *s3.S3
}

func NewClient(endpoint, region, accessKey, accessKeySecret string) (*Client, error) {
	if region == "" {
		region = "us-east-1"
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint), // Ceph RGW 端点
		Credentials: credentials.NewStaticCredentials(accessKey, accessKeySecret, ""),
	})
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	cli := &Client{
		endpoint:        endpoint,
		region:          region,
		accessKey:       accessKey,
		accessKeySecret: accessKeySecret,
		s3Cli:           svc,
	}
	return cli, nil
}

func (cli *Client) Download(_ context.Context, bucketName, objectKey, localDir string) error {
	// 创建本地文件路径
	localFilePath := filepath.Join(localDir, objectKey)
	err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %v", filepath.Dir(localFilePath), err)
	}

	// 下载文件
	output, err := cli.s3Cli.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to get object %s: %v", objectKey, err)
	}
	defer output.Body.Close()

	// 创建本地文件
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", localFilePath, err)
	}
	defer localFile.Close()

	// 将 S3 对象内容写入本地文件
	_, err = io.Copy(localFile, output.Body)
	if err != nil {
		return fmt.Errorf("failed to copy object %s to file %s: %v", objectKey, localFilePath, err)
	}

	return nil

}

func (cli *Client) ListAndDownloadDir(ctx context.Context, bucketName, prefix, localBasePath string) error {
	if strings.HasSuffix(prefix, "/") {
		// 如果 key 以斜杠结尾，表示下载目录
		return cli.DownloadDirectory(ctx, bucketName, prefix, localBasePath)
	} else {
		// 否则，表示下载单个文件
		return cli.Download(ctx, bucketName, prefix, localBasePath)
	}
}

func (cli *Client) DownloadDirectory(ctx context.Context, bucket, prefix, localDir string) error {
	err := cli.s3Cli.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, object := range page.Contents {
			err := cli.Download(ctx, bucket, *object.Key, localDir)
			if err != nil {
				log.Printf("Failed to download file %s: %v", *object.Key, err)
				return false
			}
		}
		return !lastPage
	})
	return err
}
