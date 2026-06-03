package config

var ExtConfig Extend

// Extend 扩展配置
//  extend:
//    demo:
//      name: demo-name
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	AMap    AMap
	Storage Storage
}

type AMap struct {
	Key string
}

type Storage struct {
	Driver          string
	Endpoint        string
	PublicEndpoint  string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	UseSSL          bool
}
