package config

import (
	_ "embed"
	"encoding/json"
	"os"

	defaults "github.com/mcuadros/go-defaults"
	toml "github.com/pelletier/go-toml"
)

var (
	cfg *Config
)

// Config 配置
type Config struct {
	Clouds CloudsConfig `toml:"clouds"`
}

// CloudsConfig 代表云端配置
type CloudsConfig struct {
	Platform          string               `toml:"platform"`
	UserS3FsEndpoint  string               `toml:"user_s3fs_endpoint"`
	UserS3FsBucket    string               `toml:"user_s3fs_bucket"`
	UserS3FsAccessKey string               `toml:"user_s3fs_access_key"`
	UserS3FsSecretKey string               `toml:"user_s3fs_secret_key"`
	UserFsMountPoint  string               `toml:"user_fs_mount_point"`
	UserFsSrcAddr     string               `toml:"user_fs_src_addr"`
	UserFSMountPoint  string               `toml:"user_fs_mount_point"`
	DstPath           string               `toml:"dst_path"`
	SrcDataPath       string               `toml:"src_data_path"`
	ALIOSSConfig      ALIOSSConfig         `toml:"ali-oss"`
	WanJiePublicFS    WanJiePublicFSConfig `toml:"wanjie-public-fs"`
	WanJieS3          WanJieS3Config       `toml:"wanjie-s3"`
}

// ALIOSSConfig 代表阿里云OSS配置
type ALIOSSConfig struct {
	EndpointAddr string `toml:"endpoint_addr"`
	AccessKey    string `toml:"access_key"`
	SecretKey    string `toml:"secret_key"`
	Bucket       string `toml:"bucket"`
}

// WanJiePublicFSConfig 代表万界公共文件系统配置
type WanJiePublicFSConfig struct {
	PublicFSSrcAddr    string `toml:"public_fs_src_addr"`
	PublicFSMountPoint string `toml:"public_fs_mount_point"`
}

// WanJieS3Config 代表万界s3配置
type WanJieS3Config struct {
	EndpointAddr string `toml:"endpoint_addr"`
	Region       string `toml:"region"`
	AccessKey    string `toml:"access_key"`
	SecretKey    string `toml:"secret_key"`
	Bucket       string `toml:"bucket"`
}

func (c *CloudsConfig) String() string {
	bs, _ := json.MarshalIndent(c, "", "  ")
	return string(bs)
}

func loadConfigFromBytes(cfgBytes []byte) (*Config, error) {
	cfg = new(Config)
	defaults.SetDefaults(cfg)
	err := toml.Unmarshal(cfgBytes, cfg)
	if err != nil {
		return nil, err
	}

	err = checkConfig(cfg)
	return cfg, err
}

func Init(p string) (*Config, error) {
	cfgBytes, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	if cfg, err = loadConfigFromBytes(cfgBytes); err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetCfg() *Config {
	return cfg
}

func checkConfig(cfg *Config) error {
	// check log config
	return nil
}
