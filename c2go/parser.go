// Package c2go provides functionality to parse C header files and generate Go wrapper code.
package c2go

/*
#cgo CFLAGS: -I/usr/lib/llvm-14/include
#cgo LDFLAGS: -L/usr/lib/llvm-14/lib -lclang
#include <clang-c/Index.h>
#include <stdlib.h>
#include <stdint.h>

// Global visitor context pointer
static void* g_visitor_context = NULL;
static void* g_visitor_results = NULL;

// C callback function that calls into Go
extern enum CXChildVisitResult goVisitChild(CXCursor cursor, CXCursor parent, uintptr_t handle);

static enum CXChildVisitResult visitCallback(CXCursor cursor, CXCursor parent, CXClientData client_data) {
    return goVisitChild(cursor, parent, (uintptr_t)client_data);
}

static void visitCursorChildren(CXCursor cursor, uintptr_t handle) {
    clang_visitChildren(cursor, visitCallback, (CXClientData)handle);
}

// Helper functions for Go to access cursor information
static enum CXCursorKind getCursorKind(CXCursor cursor) {
    return clang_getCursorKind(cursor);
}

static CXString getCursorSpelling(CXCursor cursor) {
    return clang_getCursorSpelling(cursor);
}

static CXType getCursorType(CXCursor cursor) {
    return clang_getCursorType(cursor);
}

static CXString getTypeSpelling(CXType type) {
    return clang_getTypeSpelling(type);
}

static CXType getCursorResultType(CXCursor cursor) {
    return clang_getCursorResultType(cursor);
}

static unsigned getCursorNumArguments(CXCursor cursor) {
    return clang_Cursor_getNumArguments(cursor);
}

static CXCursor getCursorArgument(CXCursor cursor, unsigned i) {
    return clang_Cursor_getArgument(cursor, i);
}

static long long getEnumConstantDeclValue(CXCursor cursor) {
    return clang_getEnumConstantDeclValue(cursor);
}
*/
import "C"

import (
	"fmt"
	"os"
	"runtime/cgo"
	"unsafe"
)

type Function struct {
	Name       string
	ReturnType string
	Params     []Param
}

type Param struct {
	Name string
	Type string
}

type Enum struct {
	Name   string
	Values []EnumValue
}

type EnumValue struct {
	Name  string
	Value int64
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

type Typedef struct {
	Name    string
	Aliased string
}

type ParsedHeader struct {
	Functions []Function
	Enums     []Enum
	Structs   []Struct
	Typedefs  []Typedef
}

// ParseContext 解析上下文，用于在访问者函数间传递数据
type ParseContext struct {
	Results *ParsedHeader
}

// ParseHeader 解析 C 头文件
func ParseHeader(filePath string) (*ParsedHeader, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// 创建索引
	cFilePath := C.CString(filePath)
	defer C.free(unsafe.Pointer(cFilePath))

	index := C.clang_createIndex(0, 0)
	defer C.clang_disposeIndex(index)

	// 解析翻译单元
	tu := C.clang_parseTranslationUnit(
		index,
		cFilePath,
		nil, 0,
		nil, 0,
		C.CXTranslationUnit_None,
	)
	if tu == nil {
		return nil, fmt.Errorf("failed to parse translation unit for file: %s", filePath)
	}
	defer C.clang_disposeTranslationUnit(tu)

	// 创建解析上下文
	context := &ParseContext{
		Results: &ParsedHeader{},
	}

	// 获取根游标并开始解析
	rootCursor := C.clang_getTranslationUnitCursor(tu)

	handle := cgo.NewHandle(context)
	defer handle.Delete()

	C.visitCursorChildren(rootCursor, C.uintptr_t(handle))

	return context.Results, nil
}

// collectDeclarations 手动收集顶级声明
func collectDeclarations(cursor C.CXCursor, context *ParseContext) {
	cursorKind := C.clang_getCursorKind(cursor)

	switch cursorKind {
	case C.CXCursor_FunctionDecl:
		fn := parseFunctionFromCursor(cursor)
		if fn != nil {
			context.Results.Functions = append(context.Results.Functions, *fn)
		}

	case C.CXCursor_EnumDecl:
		enum := parseEnumFromCursor(cursor)
		if enum != nil {
			context.Results.Enums = append(context.Results.Enums, *enum)
		}

	case C.CXCursor_StructDecl:
		struct_ := parseStructFromCursor(cursor)
		if struct_ != nil {
			context.Results.Structs = append(context.Results.Structs, *struct_)
		}

	case C.CXCursor_TypedefDecl:
		typedef := parseTypedefFromCursor(cursor)
		if typedef != nil {
			context.Results.Typedefs = append(context.Results.Typedefs, *typedef)
		}
	}

	// 子节点遍历由 clang_visitChildren 驱动，这里无需额外处理
}

// parseFunctionFromCursor 从游标解析函数
func parseFunctionFromCursor(cursor C.CXCursor) *Function {
	name := getCursorSpelling(cursor)
	returnType := getTypeSpelling(C.getCursorResultType(cursor))

	if name == "" {
		return nil
	}

	fn := Function{
		Name:       name,
		ReturnType: returnType,
	}

	// 解析参数
	paramCount := C.getCursorNumArguments(cursor)
	for i := C.uint(0); i < paramCount; i++ {
		paramCursor := C.getCursorArgument(cursor, i)
		if C.clang_Cursor_isNull(paramCursor) != 0 {
			continue
		}

		paramName := getCursorSpelling(paramCursor)
		paramType := getTypeSpelling(C.getCursorType(paramCursor))

		if paramName == "" {
			paramName = fmt.Sprintf("arg%d", i+1)
		}

		fn.Params = append(fn.Params, Param{
			Name: paramName,
			Type: paramType,
		})
	}

	return &fn
}

// parseEnumFromCursor 从游标解析枚举
func parseEnumFromCursor(cursor C.CXCursor) *Enum {
	name := getCursorSpelling(cursor)
	if name == "" {
		return nil
	}

	enum := Enum{Name: name}

	// 由于缺少遍历子节点的函数，暂时返回基本枚举结构
	// 在实际使用中，可以扩展这部分来解析枚举值

	return &enum
}

// parseStructFromCursor 从游标解析结构体
func parseStructFromCursor(cursor C.CXCursor) *Struct {
	name := getCursorSpelling(cursor)
	if name == "" {
		return nil
	}

	struct_ := Struct{Name: name}

	// 由于缺少遍历子节点的函数，暂时返回基本结构体结构
	// 在实际使用中，可以扩展这部分来解析字段

	return &struct_
}

// parseTypedefFromCursor 从游标解析类型定义
func parseTypedefFromCursor(cursor C.CXCursor) *Typedef {
	name := getCursorSpelling(cursor)
	underlyingType := getTypeSpelling(C.getCursorType(cursor))

	if name == "" {
		return nil
	}

	return &Typedef{
		Name:    name,
		Aliased: underlyingType,
	}
}

//export goVisitChild
func goVisitChild(cursor C.CXCursor, parent C.CXCursor, handle C.uintptr_t) C.enum_CXChildVisitResult {
	ctxHandle := cgo.Handle(handle)
	ctx, ok := ctxHandle.Value().(*ParseContext)
	if !ok {
		return C.CXChildVisit_Break
	}

	collectDeclarations(cursor, ctx)
	return C.CXChildVisit_Continue
}

// Helper functions for accessing Clang data
func getCursorSpelling(cursor C.CXCursor) string {
	cstr := C.getCursorSpelling(cursor)
	defer C.clang_disposeString(cstr)

	return C.GoString(C.clang_getCString(cstr))
}

func getTypeSpelling(t C.CXType) string {
	cstr := C.getTypeSpelling(t)
	defer C.clang_disposeString(cstr)

	return C.GoString(C.clang_getCString(cstr))
}
