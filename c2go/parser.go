// Package c2go provides functionality to parse C header files and generate Go wrapper code.
package c2go

/*
#cgo CFLAGS: -I/usr/lib/llvm-14/include
#cgo LDFLAGS: -L/usr/lib/llvm-14/lib -lclang
#include <clang-c/Index.h>
#include <stdlib.h>
#include <stdint.h>

// C callback function that calls into Go
extern enum CXChildVisitResult goVisitChild(CXCursor cursor, CXCursor parent, uintptr_t handle);

// 使用全局句柄避免把小整数句柄当作指针传入 Go 栈
static uintptr_t g_visit_handle = 0;

static enum CXChildVisitResult visitCallback(CXCursor cursor, CXCursor parent, CXClientData client_data) {
    (void)client_data; // 未使用
    return goVisitChild(cursor, parent, g_visit_handle);
}

static void visitCursorChildren(CXCursor cursor, uintptr_t handle) {
    uintptr_t prev = g_visit_handle;
    g_visit_handle = handle;
    // 传 NULL，避免在 Go 栈上出现无效指针；并在返回后恢复句柄，保证可重入
    clang_visitChildren(cursor, visitCallback, NULL);
    g_visit_handle = prev;
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

static int isCursorFromMainFile(CXCursor cursor) {
    CXSourceLocation loc = clang_getCursorLocation(cursor);
    return clang_Location_isFromMainFile(loc);
}

// 新增: typedef/类型相关辅助
static CXType getTypedefDeclUnderlyingType(CXCursor cursor) {
    return clang_getTypedefDeclUnderlyingType(cursor);
}
static enum CXTypeKind getTypeKind(CXType t) {
    return t.kind;
}
static CXCursor getTypeDeclaration(CXType t) {
    return clang_getTypeDeclaration(t);
}
static CXType getNamedType(CXType t) {
    return clang_Type_getNamedType(t);
}
static CXType getCanonicalType(CXType t) {
    return clang_getCanonicalType(t);
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

type visitFunc func(C.CXCursor, C.CXCursor) C.enum_CXChildVisitResult

type visitHandle struct {
	fn visitFunc
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

	C.visitCursorChildren(rootCursor, C.uintptr_t(uintptr(handle)))

	return context.Results, nil
}

// collectDeclarations 手动收集顶级声明
func collectDeclarations(cursor C.CXCursor, context *ParseContext) {
	if C.isCursorFromMainFile(cursor) == 0 {
		return
	}

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
		// 先尝试 typedef 到枚举：若是，则仅加入 Enums
		if enum := parseEnumFromTypedef(cursor); enum != nil {
			context.Results.Enums = append(context.Results.Enums, *enum)
			break
		}
		// 再尝试 typedef 到结构体：若是，则仅加入 Structs
		if st := parseStructFromTypedef(cursor); st != nil {
			context.Results.Structs = append(context.Results.Structs, *st)
			break
		}
		// 普通 typedef
		if typedef := parseTypedefFromCursor(cursor); typedef != nil {
			context.Results.Typedefs = append(context.Results.Typedefs, *typedef)
		}
	}

	// 子节点的深入解析在各自的 parse 函数中完成
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

	visitChildren(cursor, func(child, _ C.CXCursor) C.enum_CXChildVisitResult {
		if C.clang_getCursorKind(child) != C.CXCursor_EnumConstantDecl {
			return C.CXChildVisit_Continue
		}

		valueName := getCursorSpelling(child)
		if valueName == "" {
			return C.CXChildVisit_Continue
		}

		enum.Values = append(enum.Values, EnumValue{
			Name:  valueName,
			Value: int64(C.getEnumConstantDeclValue(child)),
		})

		return C.CXChildVisit_Continue
	})

	return &enum
}

// parseStructFromCursor 从游标解析结构体
func parseStructFromCursor(cursor C.CXCursor) *Struct {
	name := getCursorSpelling(cursor)
	if name == "" {
		return nil
	}

	struct_ := Struct{Name: name}

	visitChildren(cursor, func(child, _ C.CXCursor) C.enum_CXChildVisitResult {
		kind := C.clang_getCursorKind(child)

		switch kind {
		case C.CXCursor_FieldDecl:
			fieldName := getCursorSpelling(child)
			if fieldName == "" {
				fieldName = fmt.Sprintf("field%d", len(struct_.Fields)+1)
			}
			fieldType := getTypeSpelling(C.getCursorType(child))

			struct_.Fields = append(struct_.Fields, Field{
				Name: fieldName,
				Type: fieldType,
			})
		case C.CXCursor_StructDecl, C.CXCursor_UnionDecl:
			return C.CXChildVisit_Recurse
		}

		return C.CXChildVisit_Continue
	})

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

// parseEnumFromTypedef: 当 typedef 的底层类型为枚举时，解析其常量并用 typedef 名作为枚举名
func parseEnumFromTypedef(cursor C.CXCursor) *Enum {
	ut := C.getTypedefDeclUnderlyingType(cursor)

	typedefName := getCursorSpelling(cursor)
	if typedefName == "" {
		return nil
	}

	// 依次尝试从 underlying/named/canonical 类型上找到 EnumDecl
	enumCursor := C.getTypeDeclaration(ut)
	if C.clang_Cursor_isNull(enumCursor) != 0 || C.clang_getCursorKind(enumCursor) != C.CXCursor_EnumDecl {
		t2 := C.getNamedType(ut)
		tmp := C.getTypeDeclaration(t2)
		if C.clang_Cursor_isNull(tmp) == 0 && C.clang_getCursorKind(tmp) == C.CXCursor_EnumDecl {
			enumCursor = tmp
		}
	}
	if C.clang_Cursor_isNull(enumCursor) != 0 || C.clang_getCursorKind(enumCursor) != C.CXCursor_EnumDecl {
		t3 := C.getCanonicalType(ut)
		tmp := C.getTypeDeclaration(t3)
		if C.clang_Cursor_isNull(tmp) == 0 && C.clang_getCursorKind(tmp) == C.CXCursor_EnumDecl {
			enumCursor = tmp
		}
	}
	if C.clang_Cursor_isNull(enumCursor) != 0 || C.clang_getCursorKind(enumCursor) != C.CXCursor_EnumDecl {
		// 仍未找到有效的枚举声明
		return nil
	}

	e := Enum{Name: typedefName}
	visitChildren(enumCursor, func(child, _ C.CXCursor) C.enum_CXChildVisitResult {
		if C.clang_getCursorKind(child) != C.CXCursor_EnumConstantDecl {
			return C.CXChildVisit_Continue
		}
		valueName := getCursorSpelling(child)
		if valueName == "" {
			return C.CXChildVisit_Continue
		}
		e.Values = append(e.Values, EnumValue{
			Name:  valueName,
			Value: int64(C.getEnumConstantDeclValue(child)),
		})
		return C.CXChildVisit_Continue
	})

	return &e
}

// parseStructFromTypedef: 当 typedef 的底层类型为结构体时，解析其字段并用 typedef 名作为结构体名
func parseStructFromTypedef(cursor C.CXCursor) *Struct {
	typedefName := getCursorSpelling(cursor)
	if typedefName == "" {
		return nil
	}
	ut := C.getTypedefDeclUnderlyingType(cursor)

	// 依次尝试 underlying/named/canonical 找到 StructDecl
	rec := C.getTypeDeclaration(ut)
	if C.clang_Cursor_isNull(rec) != 0 || C.clang_getCursorKind(rec) != C.CXCursor_StructDecl {
		t2 := C.getNamedType(ut)
		tmp := C.getTypeDeclaration(t2)
		if C.clang_Cursor_isNull(tmp) == 0 && C.clang_getCursorKind(tmp) == C.CXCursor_StructDecl {
			rec = tmp
		}
	}
	if C.clang_Cursor_isNull(rec) != 0 || C.clang_getCursorKind(rec) != C.CXCursor_StructDecl {
		t3 := C.getCanonicalType(ut)
		tmp := C.getTypeDeclaration(t3)
		if C.clang_Cursor_isNull(tmp) == 0 && C.clang_getCursorKind(tmp) == C.CXCursor_StructDecl {
			rec = tmp
		}
	}
	if C.clang_Cursor_isNull(rec) != 0 || C.clang_getCursorKind(rec) != C.CXCursor_StructDecl {
		return nil
	}

	st := Struct{Name: typedefName}
	visitChildren(rec, func(child, _ C.CXCursor) C.enum_CXChildVisitResult {
		switch C.clang_getCursorKind(child) {
		case C.CXCursor_FieldDecl:
			fieldName := getCursorSpelling(child)
			if fieldName == "" {
				fieldName = fmt.Sprintf("field%d", len(st.Fields)+1)
			}
			fieldType := getTypeSpelling(C.getCursorType(child))
			st.Fields = append(st.Fields, Field{
				Name: fieldName,
				Type: fieldType,
			})
		case C.CXCursor_StructDecl, C.CXCursor_UnionDecl:
			return C.CXChildVisit_Recurse
		}
		return C.CXChildVisit_Continue
	})

	return &st
}

func visitChildren(cursor C.CXCursor, fn visitFunc) {
	if fn == nil {
		return
	}

	state := &visitHandle{fn: fn}
	handle := cgo.NewHandle(state)
	defer handle.Delete()

	C.visitCursorChildren(cursor, C.uintptr_t(uintptr(handle)))
}

//export goVisitChild
func goVisitChild(cursor C.CXCursor, parent C.CXCursor, handle C.uintptr_t) C.enum_CXChildVisitResult {
	if uintptr(handle) == 0 {
		return C.CXChildVisit_Break
	}

	ctxHandle := cgo.Handle(uintptr(handle))
	switch v := ctxHandle.Value().(type) {
	case *ParseContext:
		collectDeclarations(cursor, v)
		return C.CXChildVisit_Continue
	case *visitHandle:
		if v.fn == nil {
			return C.CXChildVisit_Continue
		}
		return v.fn(cursor, parent)
	default:
		return C.CXChildVisit_Break
	}
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
