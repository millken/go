package c2go

import (
	"strings"
)

// 基础类型映射
var typeMap = map[string]string{
	// 整型
	"int":                "C.int",
	"unsigned int":       "C.uint",
	"short":              "C.short",
	"unsigned short":     "C.ushort",
	"long":               "C.long",
	"unsigned long":      "C.ulong",
	"long long":          "C.longlong",
	"unsigned long long": "C.ulonglong",

	// 字符型
	"char":          "C.char",
	"unsigned char": "C.uchar",
	"signed char":   "C.schar",

	// 浮点型
	"float":       "C.float",
	"double":      "C.double",
	"long double": "C.longdouble",

	// 空类型
	"void": "",

	// 标准库类型
	"size_t":    "C.size_t",
	"ssize_t":   "C.ssize_t",
	"ptrdiff_t": "C.ptrdiff_t",
	"uintptr_t": "C.uintptr_t",
	"intptr_t":  "C.intptr_t",
	"time_t":    "C.time_t",

	// 布尔类型 (C99)
	"bool":  "C.bool",
	"_Bool": "C.bool",
}

// 指针类型映射
var pointerTypeMap = map[string]string{
	"void*":          "unsafe.Pointer",
	"char*":          "*C.char",
	"const char*":    "*C.char",
	"unsigned char*": "*C.uchar",
	"int*":           "*C.int",
	"float*":         "*C.float",
	"double*":        "*C.double",
	"void**":         "**unsafe.Pointer",
	"char**":         "**C.char",
}

// Go原生类型映射
var goNativeTypeMap = map[string]string{
	"char":               "byte",
	"unsigned char":      "byte",
	"short":              "int16",
	"unsigned short":     "uint16",
	"int":                "int32",
	"unsigned int":       "uint32",
	"long":               "int64",  // 在64位系统上
	"unsigned long":      "uint64", // 在64位系统上
	"long long":          "int64",
	"unsigned long long": "uint64",
	"float":              "float32",
	"double":             "float64",
	"size_t":             "uintptr",
	"bool":               "bool",
}

func MapCTypeToGo(cType string) string {
	cType = strings.TrimSpace(cType)

	// 直接映射
	if mapped, ok := typeMap[cType]; ok {
		return mapped
	}

	// 指针类型映射
	if mapped, ok := pointerTypeMap[cType]; ok {
		return mapped
	}

	// 处理指针类型
	if strings.HasSuffix(cType, "*") {
		baseType := strings.TrimSpace(cType[:len(cType)-1])
		if mappedBase, ok := typeMap[baseType]; ok {
			if mappedBase == "" { // void*
				return "unsafe.Pointer"
			}
			return "*" + mappedBase
		}
		return "*C." + sanitizeType(baseType)
	}

	// 处理const修饰符
	if strings.HasPrefix(cType, "const ") {
		baseType := strings.TrimSpace(cType[6:])
		return MapCTypeToGo(baseType)
	}

	// 处理数组类型
	if strings.Contains(cType, "[") {
		return MapArrayType(cType)
	}

	// 处理函数指针类型
	if strings.Contains(cType, ")(") {
		return MapFunctionPointerType(cType)
	}

	// 处理自定义类型
	return "C." + sanitizeType(cType)
}

// MapCTypeToGoNative 将C类型映射到Go原生类型
func MapCTypeToGoNative(cType string) string {
	cType = strings.TrimSpace(cType)

	// 直接映射
	if mapped, ok := goNativeTypeMap[cType]; ok {
		return mapped
	}

	// 处理指针类型
	if strings.HasSuffix(cType, "*") {
		baseType := strings.TrimSpace(cType[:len(cType)-1])
		if mappedBase, ok := goNativeTypeMap[baseType]; ok {
			return "*" + mappedBase
		}
		return "*" + baseType
	}

	// 处理const修饰符
	if strings.HasPrefix(cType, "const ") {
		baseType := strings.TrimSpace(cType[6:])
		return MapCTypeToGoNative(baseType)
	}

	return cType
}

// MapArrayType 处理数组类型
func MapArrayType(arrayType string) string {
	// 简化处理，返回指针类型
	baseType := strings.Split(arrayType, "[")[0]
	return "*" + MapCTypeToGo(baseType)
}

// MapFunctionPointerType 处理函数指针类型
func MapFunctionPointerType(funcType string) string {
	// 简化处理，返回unsafe.Pointer
	return "unsafe.Pointer"
}

// sanitizeType 清理类型名称，移除非法字符
func sanitizeType(typeName string) string {
	// 替换空格为下划线
	typeName = strings.ReplaceAll(typeName, " ", "_")
	// 移除其他非法字符
	var result strings.Builder
	for _, r := range typeName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// IsPointerType 检查是否为指针类型
func IsPointerType(cType string) bool {
	return strings.HasSuffix(strings.TrimSpace(cType), "*")
}

// GetBaseType 获取基础类型（去除指针）
func GetBaseType(cType string) string {
	cType = strings.TrimSpace(cType)
	if IsPointerType(cType) {
		return strings.TrimSpace(cType[:len(cType)-1])
	}
	return cType
}

// IsNumericType 检查是否为数值类型
func IsNumericType(cType string) bool {
	baseType := GetBaseType(cType)
	numericTypes := []string{
		"int", "short", "long", "char", "float", "double",
		"unsigned int", "unsigned short", "unsigned long", "unsigned char",
		"long long", "unsigned long long", "size_t", "ssize_t",
	}

	for _, nt := range numericTypes {
		if strings.Contains(baseType, nt) {
			return true
		}
	}
	return false
}
