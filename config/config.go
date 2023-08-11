package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"os"
	"share.ac.cn/common/rsa"
)

// 系统配置，对应yml
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便

// Conf 全局配置变量
var Conf = new(config)

type config struct {
	System    *SystemConfig    `mapstructure:"system" json:"system"`
	File      *FileConfig      `mapstructure:"file" json:"file"`
	Logs      *LogsConfig      `mapstructure:"logs" json:"logs"`
	Redis     *RedisConfig     `mapstructure:"redis" json:"redis"`
	Jwt       *JwtConfig       `mapstructure:"jwt" json:"jwt"`
	RateLimit *RateLimitConfig `mapstructure:"rate-limit" json:"rateLimit"`
	QinNiu    *QiniuConfig     `mapstructure:"qiniu" json:"qiniu"`
}

func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/")
	// 读取配置信息
	err = viper.ReadInConfig()

	// 热更新配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 将读取的配置信息保存至全局变量Conf
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
		}
		// 读取rsa key
		Conf.System.RSAPublicBytes = rsa.RSAReadKeyFromFile(Conf.System.RSAPublicKey)
		Conf.System.RSAPrivateBytes = rsa.RSAReadKeyFromFile(Conf.System.RSAPrivateKey)
	})

	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}
	// 读取rsa key
	Conf.System.RSAPublicBytes = rsa.RSAReadKeyFromFile(Conf.System.RSAPublicKey)
	Conf.System.RSAPrivateBytes = rsa.RSAReadKeyFromFile(Conf.System.RSAPrivateKey)
}

type SystemConfig struct {
	Mode            string `mapstructure:"mode" json:"mode"`
	Port            int    `mapstructure:"port" json:"port"`
	RSAPublicKey    string `mapstructure:"rsa-public-key" json:"rsaPublicKey"`
	RSAPrivateKey   string `mapstructure:"rsa-private-key" json:"rsaPrivateKey"`
	RSAPublicBytes  []byte `mapstructure:"-" json:"-"`
	RSAPrivateBytes []byte `mapstructure:"-" json:"-"`
	WsPort          int    `mapstructure:"ws-port" json:"ws-port"`
	Host            string `mapstructure:"host" json:"host"`
	HttpBaseWeb     string `mapstructure:"http-base-web" json:"http-base-web"`
	UploadBaseUrl   string `mapstructure:"upload-base-url" json:"upload-base-url"`
}
type JwtConfig struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}
type FileConfig struct {
	AllowMaxSize        int    `mapstructure:"allow-max-size" json:"allow-max-size"`
	DownloadTokenExpire int    `mapstructure:"download-token-expire" json:"download-token-expire"`
	ReNameFile          bool   `mapstructure:"rename-file" json:"rename-file"`
	DownloadUseToken    bool   `mapstructure:"download-use-token" json:"download-use-token"`
	EnableDistinctFile  bool   `mapstructure:"enable-distinct-file" json:"enable-distinct-file"`
	FileSumArithmetic   string `mapstructure:"file-sum-arithmetic" json:"file-sum-arithmetic"`
	AllowExt            string `mapstructure:"allow-extensions" json:"allow-extensions"`
	StoreBasePath       string `mapstructure:"store-base-path" json:"store-base-path"`
}
type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}
type RedisConfig struct {
	Address  string `mapstructure:"address" json:"address"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	Db       int    `mapstructure:"db" json:"db"`
	Size     int    `mapstructure:"size" json:"size"`
	ConnMax  int    `mapstructure:"conn-max" json:"conn-max"`
}

type RateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
}
type QiniuConfig struct {
	Zone        string `mapstructure:"zone" json:"zone"`
	Bucket      string `mapstructure:"bucket" json:"bucket"`
	DownUrl     string `mapstructure:"down-url" json:"down-url"`
	CallbackUrl string `mapstructure:"callback-url" json:"callback-url"`
	AccessKey   string `mapstructure:"access-key" json:"access-key"`
	SecretKey   string `mapstructure:"secret-key" json:"secret-key"`
}
