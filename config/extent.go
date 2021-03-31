package config

var ExtConfig Extand

// 扩展配置
//  extend:
//    demo:
//      name: demo-name
// 使用方法： config.ExtConfig......即可！！
type Extand struct {
	// Demo Demo   // 这里配置对应配置文件的结构即可
}

type Demo struct {
	Name string
}