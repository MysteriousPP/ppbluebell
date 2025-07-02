package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	bucketName string
	region     string

	imgPath   string = "image"
	cosClient *cos.Client
)

func init() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatalf("加载环境变量失败: %v", err)
	}

	// 获取环境变量
	bucketName = os.Getenv("BUCKETNAME")
	region = os.Getenv("REGION")

	if bucketName == "" || region == "" {
		log.Fatal("BUCKETNAME 或 REGION 环境变量未设置")
	}

	// 初始化 COS 客户端
	var err error
	cosClient = createCOSClient()
	if err != nil {
		log.Fatalf("初始化 COS 客户端失败: %v", err)
	}

	fmt.Printf("COS 客户端初始化成功，存储桶: %s, 区域: %s\n", bucketName, region)
}
func createCOSClient() *cos.Client {
	//client := NewClient(uri *BaseURL, httpClient *http.Client)
	// 将 examplebucket-1250000000 和 COS_REGION 修改为用户真实的信息
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	//u, _ := url.Parse("https://pp-blog-images-1356288435.cos.ap-guangzhou.myqcloud.com")
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucketName, region))
	// 用于 Get Service 查询，默认全地域 service.cos.myqcloud.com
	su, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", region))
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	// 1.永久密钥
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: os.Getenv("SECRETKEY"), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return c
}

func generateImageID(originalFileName string) string {
	// 生成 UUID
	uniqueID := uuid.New().String()

	// 获取文件扩展名
	ext := filepath.Ext(originalFileName)

	// 组合：时间戳 + UUID + 扩展名
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d_%s%s", timestamp, uniqueID, ext)
}

func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile(imgPath)
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, gin.H{"error": err.Error()})
		return
	}

	// 生成唯一文件名
	imageID := generateImageID(file.Filename)

	// 打开文件
	src, err := file.Open()
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, gin.H{"error": err.Error()})
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// 上传到 COS
	// 创建 COS 客户端的函数
	cosPath := fmt.Sprintf("images/%s", imageID) // COS 存储路径

	_, err = cosClient.Object.Put(
		context.Background(),
		cosPath,
		src,
		nil,
	)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to COS"})
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回图片 URL（COS 访问地址）
	imageURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s",
		bucketName, region, cosPath)

	// c.JSON(http.StatusOK, gin.H{
	// 	"imageId": imageID,
	// 	"url":     imageURL,
	// })
	ResponseSuccess(c, gin.H{
		"imageId": imageID,
		"url":     imageURL,
	})
}
