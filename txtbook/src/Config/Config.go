package config

import (
	"github.com/qiniu/log"
	"encoding/xml"
	"os"
	"flag"
	"time"
	"os/exec"
	"path"
)


type redisConfig struct {
	Addr, Port string
	PoolSize   int
}
type loopTime struct {
	Mail, Sms, Crawlkangle time.Duration
}

type smtpConfig struct {
	Host, Username, Password, From string
}

type beanstalkConfig struct {
	Server, MailQueue, MobileQueue string
}

type urlConfig struct {
	SmsApi, KangleIp string
}

type VsConfig struct {
	XMLName    xml.Name `xml:"Config"`
	Redis      redisConfig
	Smtp       smtpConfig
	LoopTime  loopTime
	Beanstalk   beanstalkConfig	
	Url         urlConfig
}

var vsConfig VsConfig
var ConfigFile string
var DataFile string
func init() {
	flag.StringVar(&ConfigFile,"c",defaultConfigFile() + "/config.xml","config file path")
	flag.Parse()
	if ConfigFile==""{
		ConfigFile = defaultConfigFile() + "config.xml"
	}
	log.Infof("config file: %s", ConfigFile)
	ParseXml(ConfigFile)
}

func defaultConfigFile() string {
	file, _ := exec.LookPath(os.Args[0])
	dir,_ := path.Split(file)
	os.Chdir(dir)
	wd, _ := os.Getwd()
	return wd
}

/**
解析xml文件
**/
func ParseXml(configFile string){
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	xmlObj := xml.NewDecoder(file)
	err = xmlObj.Decode(&vsConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("parse xml=%v\n",vsConfig)
}
/**
得到redis的配置
**/
func GetRedisConfig() redisConfig {
	return vsConfig.Redis
}

/**
 * get loop time
 */
func GetLoopTime() loopTime {
	return vsConfig.LoopTime
}

/*
 * smtp config
 */
func GetSmtp() smtpConfig {
	return vsConfig.Smtp
}

/*
 * beantalk queue
 */
func  GetBeanstalk() beanstalkConfig {
	return vsConfig.Beanstalk
}

/*
 * url 
 */
func  GetUrl() urlConfig {
	return vsConfig.Url
}