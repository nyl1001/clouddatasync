package clouds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/nyl1001/clouddatasync/clouds/ali"
	"github.com/nyl1001/clouddatasync/clouds/config"
	"github.com/nyl1001/clouddatasync/clouds/wanjies3"
	"github.com/nyl1001/clouddatasync/utils"
)

func TestAliSync(t *testing.T) {
	cli, err := ali.NewClient("<>", "<>", "<>") // 或者 tencent.Download
	if err != nil {
		log.Fatalf("NewClient failed: %v", err)
	}

	// 指定临时目录的前缀，例如"clouddatasync_"，os.MkdirTemp会自动添加随机后缀以保证唯一性
	dirPrefix := "clouddatasync_"
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", dirPrefix)
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	log.Println("tempDir", tempDir)
	err = cli.ListAndDownloadDir(context.Background(), "wj-devops", "drivers/NVIDIA-Linux-x86_64-535.86.05.run", tempDir)
	if err != nil {
		log.Fatalf("ListAndDownloadDir failed: %v", err)
	}

	log.Println("Data synced successfully.")
}

func TestAliSyncByEnv(t *testing.T) {
	userFsMountPoint := os.Getenv("USER_FS_MOUNT_POINT") // 万界用户文件系统挂载点根目录，对应的api参数
	dstPath := os.Getenv("DST_PATH")                     // 数据存储的相对路径，相对于万界文件系统挂载点根目录的相对路径
	platformType := os.Getenv("PLATFORM_TYPE")           // 数据平台类型，如ali-oss、wanjie-fs
	switch platformType {
	case "ali-oss":
		aliEndpointAddr := os.Getenv("ALI_ENDPOINT_ADDR") // 阿里云oss endpoint，如: https://oss-cn-beijing.aliyuncs.com
		aliAk := os.Getenv("ALI_OSS_ACCESS_KEY")          // 阿里云oss账户的access key
		aliSk := os.Getenv("ALI_OSS_SECRET_KEY")          // 阿里云oss账户的secret key
		bucketName := os.Getenv("BUCKET")                 // 阿里云oss bucket
		dataPath := os.Getenv("SRC_DATA_PATH")            // 阿里云oss bucket中的数据相对路径

		cli, err := ali.NewClient(aliEndpointAddr, aliAk, aliSk) // 或者 tencent.Download
		if err != nil {
			log.Fatalf("NewClient failed: %v", err)
		}

		// 指定临时目录的前缀，例如"clouddatasync_"，os.MkdirTemp会自动添加随机后缀以保证唯一性
		dirPrefix := "clouddatasync_"
		// 创建临时目录
		tempDir, err := os.MkdirTemp("", dirPrefix)
		if err != nil {
			log.Fatalf("创建临时目录失败: %v", err)
		}
		log.Println("tempDir", tempDir)
		err = cli.ListAndDownloadDir(context.Background(), bucketName, dataPath, tempDir)
		if err != nil {
			log.Fatalf("ListAndDownloadDir failed: %v", err)
		}

		dstDir := path.Join(userFsMountPoint, dstPath)

		err = utils.CopyDir(tempDir, dstDir)
		if err != nil {
			log.Fatalf("CopyDir failed: %v", err)
		}

	case "wanjie-public-fs":
		publicMountPoint := os.Getenv("PUBLIC_FS_MOUNT_POINT") // 万界公共文件系统挂载点根目录
		dataPath := os.Getenv("SRC_DATA_PATH")                 // 数据在万界公共文件系统根目录中的相对路径
		srcDir := path.Join(publicMountPoint, dataPath)
		dstDir := path.Join(userFsMountPoint, dstPath)
		err := utils.CopyDir(srcDir, dstDir)
		if err != nil {
			log.Fatalf("CopyDir failed: %v", err)
		}
	default:

	}

	log.Println("Data synced successfully.")
}

func TestSyncByConfigFile(t *testing.T) {
	defConfigFilePath := "/Users/nieyinliang/work/go/src/nyl1001/clouddatasync/clouds/config/config.toml"
	cfg, err := config.Init(defConfigFilePath)
	if err != nil {
		log.Fatalf("Init Config failed: %v", err)
	}
	cloudCfg := cfg.Clouds
	cfgStr, _ := json.Marshal(cloudCfg)
	fmt.Println("cloud config:", string(cfgStr))
	switch cloudCfg.Platform {
	case "ali-oss":
		aliEndpointAddr := cloudCfg.ALIOSSConfig.EndpointAddr // 阿里云oss endpoint，如: https://oss-cn-beijing.aliyuncs.com
		aliAk := cloudCfg.ALIOSSConfig.AccessKey              // 阿里云oss账户的access key
		aliSk := cloudCfg.ALIOSSConfig.SecretKey              // 阿里云oss账户的secret key
		bucketName := cloudCfg.ALIOSSConfig.Bucket            // 阿里云oss bucket
		dataPath := cloudCfg.ALIOSSConfig.SrcDataPath         // 阿里云oss bucket中的数据相对路径

		cli, err := ali.NewClient(aliEndpointAddr, aliAk, aliSk) // 或者 tencent.Download
		if err != nil {
			log.Fatalf("NewClient failed: %v", err)
		}

		// 指定临时目录的前缀，例如"clouddatasync_"，os.MkdirTemp会自动添加随机后缀以保证唯一性
		dirPrefix := "clouddatasync_"
		// 创建临时目录
		tempDir, err := os.MkdirTemp("", dirPrefix)
		if err != nil {
			log.Fatalf("创建临时目录失败: %v", err)
		}
		log.Println("tempDir", tempDir)
		err = cli.ListAndDownloadDir(context.Background(), bucketName, dataPath, tempDir)
		if err != nil {
			log.Fatalf("ListAndDownloadDir failed: %v", err)
		}

		dstDir := path.Join(cloudCfg.ALIOSSConfig.UserFSMountPoint, cloudCfg.ALIOSSConfig.DstPath)
		// 判断目录是否存在
		if _, err := os.Stat(dstDir); os.IsNotExist(err) {
			// 目录不存在，创建目录
			fmt.Printf("目录 %s 不存在，正在创建...\n", dstDir)
			err := os.MkdirAll(dstDir, os.ModePerm)
			if err != nil {
				fmt.Printf("创建目录 %s 时出错: %v\n", dstDir, err)
			} else {
				fmt.Printf("目录 %s 创建成功\n", dstDir)
			}
		}
		err = utils.CopyDir(tempDir, dstDir)
		if err != nil {
			log.Fatalf("CopyDir failed: %v", err)
		}

	case "wanjie-public-fs":
		publicMountPoint := cloudCfg.WanJiePublicFS.PublicFSMountPoint // 万界公共文件系统挂载点根目录
		dataPath := cloudCfg.WanJiePublicFS.SrcDataPath                // 数据在万界公共文件系统根目录中的相对路径
		srcDir := path.Join(publicMountPoint, dataPath)
		dstDir := path.Join(cloudCfg.WanJiePublicFS.UserFSMountPoint, cloudCfg.WanJiePublicFS.DstPath)
		err = utils.CopyDir(srcDir, dstDir)
		if err != nil {
			log.Fatalf("CopyDir failed: %v", err)
		}
	case "wanjie-s3":
		wanjieS3EndpointAddr := cloudCfg.WanJieS3.EndpointAddr // 万界s3 endpoint，如: https://oss-cn-beijing.aliyuncs.com
		wanjieS3Region := cloudCfg.WanJieS3.Region             // 万界s3 endpoint，如: https://oss-cn-beijing.aliyuncs.com
		wanjieAk := cloudCfg.WanJieS3.AccessKey                // 万界s3 账户的access key
		wanjieSk := cloudCfg.WanJieS3.SecretKey                // 万界s3 账户的secret key
		bucketName := cloudCfg.WanJieS3.Bucket                 // 万界s3 bucket
		dataPath := cloudCfg.WanJieS3.SrcDataPath              // 万界s3 bucket中的数据相对路径
		cli, err := wanjies3.NewClient(wanjieS3EndpointAddr, wanjieS3Region, wanjieAk, wanjieSk)
		if err != nil {
			log.Fatalf("NewClient failed: %v", err)
		}

		// 指定临时目录的前缀，例如"clouddatasync_"，os.MkdirTemp会自动添加随机后缀以保证唯一性
		dirPrefix := "clouddatasync_"
		// 创建临时目录
		tempDir, err := os.MkdirTemp("", dirPrefix)
		if err != nil {
			log.Fatalf("创建临时目录失败: %v", err)
		}
		log.Println("tempDir", tempDir)
		err = cli.ListAndDownloadDir(context.Background(), bucketName, dataPath, tempDir)
		if err != nil {
			log.Fatalf("ListAndDownloadDir failed: %v", err)
		}

		dstDir := path.Join(cloudCfg.WanJieS3.UserFSMountPoint, cloudCfg.WanJieS3.DstPath)
		// 判断目录是否存在
		if _, err := os.Stat(dstDir); os.IsNotExist(err) {
			// 目录不存在，创建目录
			fmt.Printf("目录 %s 不存在，正在创建...\n", dstDir)
			err := os.MkdirAll(dstDir, os.ModePerm)
			if err != nil {
				fmt.Printf("创建目录 %s 时出错: %v\n", dstDir, err)
			} else {
				fmt.Printf("目录 %s 创建成功\n", dstDir)
			}
		}
		err = utils.CopyDir(tempDir, dstDir)
		if err != nil {
			log.Fatalf("CopyDir failed: %v", err)
		}

	default:

	}

	log.Println("Data synced successfully.")
}
