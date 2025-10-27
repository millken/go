package c2go

import (
	"bytes"
	"fmt"
	"text/template"
)

const wrapperTemplate = `
package {{.PackageName}}

/*
#include "{{.HeaderFile}}"
*/
import "C"
import (
	"unsafe"
)

{{range .Functions}}
// {{.GoName}} 包装原生函数 {{.CName}}
func {{.GoName}}({{.ParamList}}) {{.ReturnType}} {
	{{if .HasReturn}}return {{end}}C.{{.CName}}({{.CallArgs}})
}
{{end}}
`

type FunctionWrapper struct {
	CName      string
	GoName     string
	ReturnType string
	ParamList  string
	CallArgs   string
	HasReturn  bool
}

type PackageTemplateData struct {
	PackageName string
	HeaderFile  string
	Functions   []FunctionWrapper
}

func GenerateWrapper(pkgName, headerFile string, funcs []Function) (string, error) {
	var wrappers []FunctionWrapper

	for _, fn := range funcs {
		wrapper := FunctionWrapper{
			CName:      fn.Name,
			GoName:     toCamelCase(fn.Name),
			ReturnType: MapCTypeToGo(fn.ReturnType),
			HasReturn:  fn.ReturnType != "void",
		}

		// 生成参数列表和调用参数
		var paramList bytes.Buffer
		var callArgs bytes.Buffer

		for i, param := range fn.Params {
			goType := MapCTypeToGo(param.Type)
			goName := safeVarName(param.Name, i)

			if i > 0 {
				paramList.WriteString(", ")
				callArgs.WriteString(", ")
			}

			paramList.WriteString(goName + " " + goType)

			// 特殊处理字符串转换
			if param.Type == "const char*" || param.Type == "char*" {
				callArgs.WriteString("C.CString(" + goName + ")")
			} else {
				callArgs.WriteString(goName)
			}
		}

		wrapper.ParamList = paramList.String()
		wrapper.CallArgs = callArgs.String()
		wrappers = append(wrappers, wrapper)
	}

	// 应用模板
	tmpl, err := template.New("wrapper").Parse(wrapperTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := PackageTemplateData{
		PackageName: pkgName,
		HeaderFile:  headerFile,
		Functions:   wrappers,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func toCamelCase(s string) string {
	var result bytes.Buffer
	upperNext := true

	for _, r := range s {
		if r == '_' {
			upperNext = true
			continue
		}

		if upperNext {
			result.WriteRune(r - 32) // 转换为大写
			upperNext = false
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func safeVarName(name string, index int) string {
	if name == "" {
		return fmt.Sprintf("arg%d", index)
	}

	// 替换无效字符
	var result bytes.Buffer
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
	}

	return result.String()
}
