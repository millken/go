package gfx

/*
#include "sokol_gfx.h"
#cgo nocallback sg_setup
#cgo nocallback sg_shutdown
#cgo nocallback sg_isvalid
#cgo nocallback sg_begin_pass
#cgo nocallback sg_end_pass
#cgo nocallback sg_commit
#cgo nocallback sg_make_buffer
#cgo nocallback sg_make_shader
#cgo nocallback sg_make_pipeline
#cgo nocallback sg_apply_pipeline
#cgo nocallback sg_apply_bindings
#cgo nocallback sg_draw

#cgo noescape sg_setup
#cgo noescape sg_shutdown
#cgo noescape sg_isvalid
#cgo noescape sg_begin_pass
#cgo noescape sg_end_pass
#cgo noescape sg_commit
#cgo noescape sg_make_buffer
#cgo noescape sg_make_shader
#cgo noescape sg_make_pipeline
#cgo noescape sg_apply_pipeline
#cgo noescape sg_apply_bindings
#cgo noescape sg_draw
*/
import "C"

import (
	"unsafe"
)

func Setup(desc *Desc) {
	C.sg_setup((*C.sg_desc)(unsafe.Pointer(desc)))
}

func Shutdown() {
	C.sg_shutdown()
}

func IsValid() bool {
	return bool(C.sg_isvalid())
}

func BeginPass(pass *Pass) {
	C.sg_begin_pass((*C.sg_pass)(unsafe.Pointer(pass)))
}

func EndPass() {
	C.sg_end_pass()
}

func Commit() {
	C.sg_commit()
}

func MakeBuffer(desc *BufferDesc) Buffer {
	buf := C.sg_make_buffer((*C.sg_buffer_desc)(unsafe.Pointer(desc)))
	return Buffer{Id: uint32(buf.id)}
}

func MakeShader(desc *ShaderDesc) Shader {
	shd := C.sg_make_shader((*C.sg_shader_desc)(unsafe.Pointer(desc)))
	return Shader{Id: uint32(shd.id)}
}

func MakePipeline(desc *PipelineDesc) Pipeline {
	pip := C.sg_make_pipeline((*C.sg_pipeline_desc)(unsafe.Pointer(desc)))
	return Pipeline{Id: uint32(pip.id)}
}

func ApplyPipeline(pipeline Pipeline) {
	C.sg_apply_pipeline(*((*C.sg_pipeline)(unsafe.Pointer(&pipeline))))
}

func ApplyBindings(bindings *Bindings) {
	C.sg_apply_bindings((*C.sg_bindings)(unsafe.Pointer(bindings)))
}

func Draw(baseElement, numElements, numInstances int) {
	C.sg_draw(C.int(baseElement), C.int(numElements), C.int(numInstances))
}
