package Config

import (
	"Utils"
	"encoding/xml"
	"os"
	"flag"
	"time"
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
	flag.StringVar(&ConfigFile,"c","config.xml","config file path")
	flag.Parse()
	if ConfigFile==""{
		ConfigFile="./config.xml"
	}
	ParseXml(ConfigFile)
}
/**
解析xml文件
**/
func ParseXml(configFile string){
	file, err := os.Open(configFile)
	if err != nil {
		Utils.LogPanicErr(err)
		return
	}
	xmlObj := xml.NewDecoder(file)
	err = xmlObj.Decode(&vsConfig)
	if err != nil {
		Utils.LogPanicErr(err)
		return
	}
	Utils.LogInfo("parse xml=%v\n",vsConfig)
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