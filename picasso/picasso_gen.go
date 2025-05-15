package picasso

// import (
// 	"fmt"
// 	"runtime"
// 	"time"

// 	"github.com/xlab/c-for-go/generator" // 导入代码生成器的包
// 	"github.com/xlab/c-for-go/parser"
// 	"github.com/xlab/c-for-go/translator" // 导入翻译器的包
// 	"gopkg.in/yaml.v3"
// )

// type (
// 	Gen struct {
// 		PackageName        string
// 		fileName           string
// 		Includes           []string //.h
// 		IncludePaths       []string //.h dir
// 		SourcesPaths       []string //.h in src dir
// 		Defines            map[string]any
// 		ConstCharIsString  bool
// 		ConstUCharIsString bool
// 	}
// )

// func GenYaml(g Gen) {
// 	processConfig := &ProcessConfig{

// 		//生成器配置
// 		Generator: &generator.Config{
// 			PackageName:        g.PackageName, // Go 代码的包名
// 			PackageDescription: "",            // Go 代码的包描述
// 			PackageLicense:     "",            // Go 代码的包许可证
// 			PkgConfigOpts:      []string{},    // PackageName-config 的选项
// 			FlagGroups: []generator.TraitFlagGroup{
// 				{
// 					Name:   "",         // 标志组名称
// 					Traits: []string{}, // 标志组特性
// 					Flags:  []string{}, // 标志组标志
// 				},
// 			},
// 			SysIncludes: []string{}, // 系统头文件
// 			Includes:    []string{}, // 头文件
// 			Options: generator.GenOptions{
// 				SafeStrings:     false, // 是否启用安全字符串
// 				StructAccessors: false, // 是否启用结构体访问器
// 				KeepAlive:       false, // 是否保持存活状态 runtime
// 			},
// 		},

// 		//翻译器配置,主要还是这里的规则
// 		Translator: &translator.Config{
// 			Rules: translator.Rules{
// 				//translator.TargetGlobal: { // 翻译规则目标  cgo_helpers.go
// 				//	translator.RuleSpec{
// 				//		From:      "",                        // 源类型
// 				//		To:        "",                        // 目标类型
// 				//		Action:    translator.ActionAccept,   // 动作
// 				//		Transform: translator.TransformUpper, // 转换
// 				//		Load:      "",                        // 加载
// 				//	},
// 				//},
// 				//translator.TargetConst: { // 翻译规则目标  const.go
// 				//	translator.RuleSpec{
// 				//		From:      "",                         // 源类型
// 				//		To:        "",                         // 目标类型
// 				//		Action:    translator.ActionAccept,    // 动作
// 				//		Transform: translator.TransformExport, // 转换
// 				//		Load:      "",                         // 加载
// 				//	},
// 				//},
// 				//translator.TargetType: { // 翻译规则目标     types.go
// 				//	translator.RuleSpec{
// 				//		From: "", // 源类型
// 				//		To:   "", // 目标类型
// 				//		//Action:    translator.ActionAccept,    // 动作
// 				//		//Transform: translator.TransformExport, // 转换
// 				//		Load: "", // 加载
// 				//	},
// 				//},
// 				translator.TargetFunction: { // 翻译规则目标
// 					translator.RuleSpec{
// 						From:      "",                         // 源类型
// 						To:        "",                         // 目标类型
// 						Action:    translator.ActionAccept,    // 动作
// 						Transform: translator.TransformExport, // 转换
// 						Load:      "",                         // 加载
// 					},
// 				},

// 				//translator.TargetPrivate: { // 翻译规则目标
// 				//	translator.RuleSpec{
// 				//		From: "", // 源类型
// 				//		To:   "", // 目标类型
// 				//		//Action:    translator.ActionAccept,    // 动作
// 				//		//Transform: translator.TransformExport, // 转换
// 				//		Load: "", // 加载
// 				//	},
// 				//},
// 				translator.TargetPublic: { // 翻译规则目标
// 					translator.RuleSpec{
// 						From:      "",                        // 源类型
// 						To:        "",                        // 目标类型
// 						Action:    translator.ActionAccept,   // 动作
// 						Transform: translator.TransformUpper, // 转换
// 						Load:      "",                        // 加载
// 					},
// 				},
// 				//translator.TargetPostGlobal: { // 翻译规则目标
// 				//	translator.RuleSpec{
// 				//		From:      "",                         // 源类型
// 				//		To:        "",                         // 目标类型
// 				//		Action:    translator.ActionAccept,    // 动作
// 				//		Transform: translator.TransformExport, // 转换
// 				//		Load:      "",                         // 加载
// 				//	},
// 				//},
// 			},
// 			ConstRules: translator.ConstRules{
// 				"": "", // 常量规则
// 			},
// 			PtrTips: translator.PtrTips{
// 				"": []translator.TipSpec{
// 					{
// 						Target: "", // 目标指针类型
// 						Tips: translator.Tips{ // 提示信息
// 							"",
// 							"",
// 							"",
// 						},
// 						Self:    "", // 自指针类型
// 						Default: "", // 默认值
// 					},
// 				},
// 			},
// 			TypeTips: translator.TypeTips{
// 				"": []translator.TipSpec{
// 					{
// 						Target: "", // 目标类型
// 						Tips: translator.Tips{
// 							"",
// 							"",
// 							"",
// 						}, // 提示信息
// 						Self:    "", // 自身类型
// 						Default: "", // 默认值
// 					},
// 				},
// 			},
// 			MemTips: []translator.TipSpec{
// 				{
// 					Target: "", // 目标类型
// 					Tips: translator.Tips{
// 						"",
// 						"",
// 						"",
// 					}, // 提示信息
// 					Self:    "", // 自身类型
// 					Default: "", // 默认值
// 				},
// 			},
// 			Typemap: map[translator.CTypeSpec]translator.GoTypeSpec{
// 				translator.CTypeSpec{
// 					Raw:      "",    // 原始类型的字符串表示
// 					Base:     "",    // 基本类型的字符串表示
// 					Const:    false, // 是否为常量
// 					Signed:   false, // 是否为有符号类型
// 					Unsigned: false, // 是否为无符号类型
// 					Short:    false, // 是否为短整型
// 					Long:     false, // 是否为长整型
// 					Complex:  false, // 是否为复数类型
// 					Opaque:   false, // 是否为不透明类型
// 					Pointers: 0,     // 指针数量
// 					InnerArr: "",    // 内部数组的长度
// 					OuterArr: "",    // 外部数组的长度
// 				}: {
// 					Slices:   0,     // 切片数量
// 					Pointers: 0,     // 指针数量
// 					InnerArr: "",    // 内部数组的长度
// 					OuterArr: "",    // 外部数组的长度
// 					Unsigned: false, // 是否为无符号类型
// 					Kind:     0,     // 基本类型的种类
// 					Base:     "",    // 基本类型的字符串表示
// 					Raw:      "",    // 原始类型的字符串表示
// 					Bits:     0,     // 位数
// 				},
// 			},
// 			ConstCharIsString:  &g.ConstCharIsString,  // 是否将常量字符视为字符串
// 			ConstUCharIsString: &g.ConstUCharIsString, // 是否将常量无符号字符视为字符串
// 			LenFields:          map[string]string{},   // 长度字段
// 			IgnoredFiles:       []string{},            // 被忽略的文件列表 todo zycore
// 			LongIs64Bit:        false,                 // long 类型是否为 64 位整数类型
// 		},

// 		//解析器配置
// 		Parser: &parser.Config{
// 			Arch:         runtime.GOARCH, // 架构类型
// 			IncludePaths: g.IncludePaths, // 包含文件的路径 todo
// 			SourcesPaths: g.SourcesPaths, // 源代码文件的路径  todo
// 			IgnoredPaths: []string{},     // 忽略的文件路径
// 			Defines:      g.Defines,      // 宏定义的 map 类型
// 			CCDefs:       false,          //cpp 是否使用 CCDefs  todo
// 			CCIncl:       false,          //cpp 是否使用 CCIncl
// 		},
// 	}
// 	out, err := yaml.Marshal(processConfig)
// 	if !mylog.Error(err) {
// 		return
// 	}
// 	tool.File().WriteTruncate(g.fileName+".yaml", out)

// 	t0 := time.Now()
// 	process, err := NewProcess(g.fileName+".yaml", g.fileName)
// 	if !mylog.Error(err) {
// 		return
// 	}
// 	process.Generate(*noCGO)
// 	if !mylog.Error(process.Flush(*noCGO)) {
// 		return
// 	}
// 	fmt.Printf("done in %v\n", time.Now().Sub(t0))
// }
