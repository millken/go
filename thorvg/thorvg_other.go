//go:build !windows

package thorvg

import (
	"slices"
	"unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo darwin LDFLAGS: -L${SRCDIR} -lc++
#cgo darwin,arm64 LDFLAGS: -lthorvg_darwin_arm64 -L/opt/homebrew/opt/libomp/lib -lomp
#cgo linux LDFLAGS: -L${SRCDIR} -lskia_linux -lfontconfig -lfreetype -lGL -ldl -lm -lstdc++

#include <stdlib.h>
#include <string.h>
#include "thorvg_capi.h"
*/
import "C"

// EngineVersion 获取 TVG 版本号
func EngineVersion() (major, minor, micro uint32, version string, res Result) {
	var cMajor, cMinor, cMicro C.uint32_t
	var cVersion *C.char
	r := C.tvg_engine_version(&cMajor, &cMinor, &cMicro, (**C.char)(unsafe.Pointer(&cVersion)))
	major = uint32(cMajor)
	minor = uint32(cMinor)
	micro = uint32(cMicro)
	if cVersion != nil {
		version = C.GoString(cVersion)
	}
	res = Result(r)
	return
}

// EngineInit 初始化 ThorVG 引擎
func EngineInit(threads uint) Result {
	return Result(C.tvg_engine_init(C.uint(threads)))
}

// EngineTerm 终止 ThorVG 引擎
func EngineTerm() Result {
	return Result(C.tvg_engine_term())
}

type Canvas struct {
	ptr *C.Tvg_Canvas
}

func NewSwcanvas() *Canvas {
	c := C.tvg_swcanvas_create()
	if c == nil {
		return nil
	}
	return &Canvas{ptr: c}
}

// SwcanvasSetTarget 设置软件渲染 Canvas 的目标缓冲区
func (c *Canvas) SwcanvasSetTarget(buffer []uint32, stride, w, h uint32, cs Colorspace) Result {
	if c == nil || c.ptr == nil || len(buffer) == 0 {
		return ResultInvalidArgument // TVG_RESULT_INVALID_ARGUMENT
	}
	return Result(C.tvg_swcanvas_set_target(
		c.ptr,
		(*C.uint32_t)(unsafe.Pointer(&buffer[0])),
		C.uint32_t(stride),
		C.uint32_t(w),
		C.uint32_t(h),
		C.Tvg_Colorspace(cs),
	))
}

// Push Push 画布
func (c *Canvas) Push(paint *Paint) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_canvas_push(c.ptr, paint.ptr))
}

// Remove 移除画布
func (c *Canvas) Remove(paint *Paint) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_canvas_remove(c.ptr, paint.ptr))
}

// Update 更新画布
func (c *Canvas) Update() Result {
	if c == nil || c.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_canvas_update(c.ptr))
}

// UpdatePaint 更新画布
func (c *Canvas) UpdatePaint(paint *Paint) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_canvas_update_paint(c.ptr, paint.ptr))
}

// Draw 绘制画布
func (c *Canvas) Draw(clear bool) Result {
	if c == nil || c.ptr == nil {
		return ResultInvalidArgument
	}

	return Result(C.tvg_canvas_draw(c.ptr, C.bool(clear)))
}

// Sync 同步画布
func (c *Canvas) Sync() Result {
	if c == nil || c.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_canvas_sync(c.ptr))
}

// TVG_API Tvg_Result tvg_canvas_set_viewport(Tvg_Canvas* canvas, int32_t x, int32_t y, int32_t w, int32_t h);
// SetViewport 设置画布视口
func (c *Canvas) SetViewport(x, y, w, h int32) Result {
	if c == nil || c.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_canvas_set_viewport(c.ptr, C.int32_t(x), C.int32_t(y), C.int32_t(w), C.int32_t(h)))
}

// TVG_API Tvg_Result tvg_paint_scale(Tvg_Paint* paint, float factor);
// Scale 缩放画布
func (c *Canvas) Scale(paint *Paint, factor float32) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_scale(paint.ptr, C.float(factor)))
}

// TVG_API Tvg_Result tvg_paint_rotate(Tvg_Paint* paint, float degree);
// Rotate 旋转画布
func (c *Canvas) Rotate(paint *Paint, degree float32) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_rotate(paint.ptr, C.float(degree)))
}

// TVG_API Tvg_Result tvg_paint_translate(Tvg_Paint* paint, float x, float y);
// Translate 平移画布
func (c *Canvas) Translate(paint *Paint, x, y float32) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_translate(paint.ptr, C.float(x), C.float(y)))
}

// TVG_API Tvg_Result tvg_paint_set_transform(Tvg_Paint* paint, const Tvg_Matrix* m);
// SetTransform 设置画布变换矩阵
func (c *Canvas) SetTransform(paint *Paint, m Matrix) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_set_transform(paint.ptr, m.CPtr()))
}

// TVG_API Tvg_Result tvg_paint_get_transform(Tvg_Paint* paint, Tvg_Matrix* m);
// GetTransform 获取画布变换矩阵
func (c *Canvas) GetTransform(paint *Paint, m *Matrix) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil || m == nil {
		return ResultInvalidArgument
	}
	var cm C.Tvg_Matrix
	res := C.tvg_paint_get_transform(paint.ptr, &cm)
	m.E11, m.E12, m.E13 = float32(cm.e11), float32(cm.e12), float32(cm.e13)
	m.E21, m.E22, m.E23 = float32(cm.e21), float32(cm.e22), float32(cm.e23)
	m.E31, m.E32, m.E33 = float32(cm.e31), float32(cm.e32), float32(cm.e33)
	return Result(res)
}

// TVG_API Tvg_Result tvg_paint_set_opacity(Tvg_Paint* paint, uint8_t opacity);
// SetOpacity 设置画布不透明度
func (c *Canvas) SetOpacity(paint *Paint, opacity uint8) Result {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_set_opacity(paint.ptr, C.uint8_t(opacity)))
}

// TVG_API Tvg_Result tvg_paint_get_opacity(const Tvg_Paint* paint, uint8_t* opacity);
// GetOpacity 获取画布不透明度
func (c *Canvas) GetOpacity(paint *Paint) (uint8, Result) {
	if c == nil || c.ptr == nil || paint == nil || paint.ptr == nil {
		return 0, ResultInvalidArgument
	}
	var copacity C.uint8_t
	res := C.tvg_paint_get_opacity(paint.ptr, &copacity)
	return uint8(copacity), Result(res)
}

type Matrix struct {
	E11, E12, E13 float32
	E21, E22, E23 float32
	E31, E32, E33 float32
}

// 转换函数（可选）
func (m Matrix) CPtr() *C.Tvg_Matrix {
	return &C.Tvg_Matrix{
		e11: C.float(m.E11),
		e12: C.float(m.E12),
		e13: C.float(m.E13),
		e21: C.float(m.E21),
		e22: C.float(m.E22),
		e23: C.float(m.E23),
		e31: C.float(m.E31),
		e32: C.float(m.E32),
		e33: C.float(m.E33),
	}
}

type Paint struct {
	ptr *C.Tvg_Paint
}

// TVG_API Tvg_Paint* tvg_paint_duplicate(Tvg_Paint* paint);
func (p *Paint) Duplicate() *Paint {
	if p == nil || p.ptr == nil {
		return nil
	}
	c := C.tvg_paint_duplicate(p.ptr)
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

// TVG_API Tvg_Result tvg_paint_get_aabb(const Tvg_Paint* paint, float* x, float* y, float* w, float* h);
func (p *Paint) GetAABB() (x, y, w, h float32, res Result) {
	if p == nil || p.ptr == nil {
		return 0, 0, 0, 0, ResultInvalidArgument
	}
	var cx, cy, cw, ch C.float
	res = Result(C.tvg_paint_get_aabb(p.ptr, &cx, &cy, &cw, &ch))
	x = float32(cx)
	y = float32(cy)
	w = float32(cw)
	h = float32(ch)
	return
}

// TVG_API Tvg_Result tvg_paint_get_obb(const Tvg_Paint* paint, Tvg_Point* pt4);
func (p *Paint) GetOBB() (pt4 [4]Point, res Result) {
	if p == nil || p.ptr == nil {
		return pt4, ResultInvalidArgument
	}
	var cpt4 [4]C.Tvg_Point
	res = Result(C.tvg_paint_get_obb(p.ptr, &cpt4[0]))
	for i := 0; i < 4; i++ {
		pt4[i].X = float32(cpt4[i].x)
		pt4[i].Y = float32(cpt4[i].y)
	}
	return
}

// TVG_API Tvg_Result tvg_paint_set_mask_method(Tvg_Paint* paint, Tvg_Paint* target, Tvg_Mask_Method method);
// SetMaskMethod 设置遮罩方法
func (p *Paint) SetMaskMethod(target *Paint, method MaskMethod) Result {
	if p == nil || p.ptr == nil || target == nil || target.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_set_mask_method(p.ptr, target.ptr, C.Tvg_Mask_Method(method)))
}

// TVG_API Tvg_Result tvg_paint_get_mask_method(const Tvg_Paint* paint, const Tvg_Paint** target, Tvg_Mask_Method* method);
// GetMaskMethod 获取遮罩方法
func (p *Paint) GetMaskMethod() (target *Paint, method MaskMethod, res Result) {
	if p == nil || p.ptr == nil {
		return nil, 0, ResultInvalidArgument
	}
	var ctarget *C.Tvg_Paint
	var cmethod C.Tvg_Mask_Method
	res = Result(C.tvg_paint_get_mask_method(p.ptr, &ctarget, &cmethod))
	if ctarget != nil {
		target = &Paint{ptr: ctarget}
	}
	method = MaskMethod(cmethod)
	return
}

// TVG_API Tvg_Result tvg_paint_clip(Tvg_Paint* paint, Tvg_Paint* clipper);
// Clip 设置裁剪器
func (p *Paint) Clip(clipper *Paint) Result {
	if p == nil || p.ptr == nil || clipper == nil || clipper.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_clip(p.ptr, clipper.ptr))
}

// TVG_API const Tvg_Paint* tvg_paint_get_parent(const Tvg_Paint* paint);
// GetParent 获取父画布
func (p *Paint) GetParent() *Paint {
	if p == nil || p.ptr == nil {
		return nil
	}
	c := C.tvg_paint_get_parent(p.ptr)
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

// TVG_API Tvg_Result tvg_paint_get_type(const Tvg_Paint* paint, Tvg_Type* type);
func (p *Paint) GetType() (Type, Result) {
	if p == nil || p.ptr == nil {
		return TypeUndef, ResultInvalidArgument
	}
	var ctype C.Tvg_Type
	res := Result(C.tvg_paint_get_type(p.ptr, &ctype))
	return Type(ctype), res
}

// TVG_API Tvg_Result tvg_paint_set_blend_method(Tvg_Paint* paint, Tvg_Blend_Method method);
// SetBlendMethod 设置混合方法
func (p *Paint) SetBlendMethod(method BlendMethod) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_paint_set_blend_method(p.ptr, C.Tvg_Blend_Method(method)))
}

/************************************************************************/
/* Shape API                                                            */
/************************************************************************/
//TVG_API Tvg_Paint* tvg_shape_new(void);
func ShapeNew() *Paint {
	c := C.tvg_shape_new()
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

// TVG_API Tvg_Result tvg_shape_reset(Tvg_Paint* paint);
// ShapeReset 重置形状
func (p *Paint) ShapeReset() Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_reset(p.ptr))
}

// TVG_API Tvg_Result tvg_shape_move_to(Tvg_Paint* paint, float x, float y);
// ShapeMoveTo 移动到指定坐标
func (p *Paint) ShapeMoveTo(x, y float32) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_move_to(p.ptr, C.float(x), C.float(y)))
}

// TVG_API Tvg_Result tvg_shape_line_to(Tvg_Paint* paint, float x, float y);
// ShapeLineTo 画线到指定坐标
func (p *Paint) ShapeLineTo(x, y float32) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_line_to(p.ptr, C.float(x), C.float(y)))
}

// TVG_API Tvg_Result tvg_shape_cubic_to(Tvg_Paint* paint, float cx1, float cy1, float cx2, float cy2, float x, float y);
// ShapeCubicTo 画三次贝塞尔曲线
func (p *Paint) ShapeCubicTo(cx1, cy1, cx2, cy2, x, y float32) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_cubic_to(p.ptr, C.float(cx1), C.float(cy1), C.float(cx2), C.float(cy2), C.float(x), C.float(y)))
}

// TVG_API Tvg_Result tvg_shape_close(Tvg_Paint* paint);
// ShapeClose 关闭形状
func (p *Paint) ShapeClose() Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_close(p.ptr))
}

// TVG_API Tvg_Result tvg_shape_append_rect(Tvg_Paint* paint, float x, float y, float w, float h, float rx, float ry, bool cw);
// ShapeAppendRect 添加矩形
func (p *Paint) ShapeAppendRect(x, y, w, h, rx, ry float32, cw bool) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_append_rect(p.ptr, C.float(x), C.float(y), C.float(w), C.float(h), C.float(rx), C.float(ry), C.bool(cw)))
}

// TVG_API Tvg_Result tvg_shape_append_circle(Tvg_Paint* paint, float cx, float cy, float rx, float ry, bool cw);
/*
* @param [in]绘制一个tvg_paint指向形状对象的指针。
* @param [in] cx椭圆中心的水平坐标。
* @param [in] cy椭圆中心的垂直坐标。
* @param [in] rx椭圆的X轴半径。
* @param [in] ry椭圆的Y轴半径。
* @param [in] cw指定路径方向：@c true for calketwise，@c false for Complclockwise。
 */
// ShapeAppendCircle 添加椭圆
func (p *Paint) ShapeAppendCircle(cx, cy, rx, ry float32, cw bool) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_append_circle(p.ptr, C.float(cx), C.float(cy), C.float(rx), C.float(ry), C.bool(cw)))
}

// TVG_API Tvg_Result tvg_shape_append_path(Tvg_Paint* paint, const Tvg_Path_Command* cmds, uint32_t cmdCnt, const Tvg_Point* pts, uint32_t ptsCnt);
// ShapeAppendPath 添加路径
func (p *Paint) ShapeAppendPath(cmds PathCommand, cmdCnt uint32, pts []Point, ptsCnt uint32) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	if len(pts) < int(ptsCnt) {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_append_path(p.ptr, (*C.Tvg_Path_Command)(unsafe.Pointer(&cmds)), C.uint32_t(cmdCnt), (*C.Tvg_Point)(unsafe.Pointer(&pts[0])), C.uint32_t(ptsCnt)))
}

// TVG_API Tvg_Result tvg_shape_get_path(const Tvg_Paint* paint, const Tvg_Path_Command** cmds, uint32_t* cmdsCnt, const Tvg_Point** pts, uint32_t* ptsCnt);
// ShapeGetPath 获取路径
func (p *Paint) ShapeGetPath() (cmds []PathCommand, pts []Point, res Result) {
	if p == nil || p.ptr == nil {
		return nil, nil, ResultInvalidArgument
	}
	var ccmds *C.Tvg_Path_Command
	var ccmdsCnt C.uint32_t
	var cpts *C.Tvg_Point
	var cptsCnt C.uint32_t
	res = Result(C.tvg_shape_get_path(p.ptr, &ccmds, &ccmdsCnt, &cpts, &cptsCnt))
	if res != ResultSuccess {
		return nil, nil, res
	}
	// 将 C 数组映射为 Go slice
	if ccmdsCnt > 0 {
		cmds = unsafe.Slice((*PathCommand)(unsafe.Pointer(ccmds)), int(ccmdsCnt))
		cmds = slices.Clone(cmds)
	}
	if cptsCnt > 0 {
		cptsSlice := unsafe.Slice((*C.Tvg_Point)(unsafe.Pointer(cpts)), int(cptsCnt))
		pts = make([]Point, cptsCnt)
		for i := 0; i < int(cptsCnt); i++ {
			pts[i].X = float32(cptsSlice[i].x)
			pts[i].Y = float32(cptsSlice[i].y)
		}
	}
	return cmds, pts, res
}

// TVG_API Tvg_Result tvg_shape_set_stroke_width(Tvg_Paint* paint, float width);
// ShapeSetStrokeWidth 设置描边宽度
func (p *Paint) ShapeSetStrokeWidth(width float32) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_stroke_width(p.ptr, C.float(width)))
}

// TVG_API Tvg_Result tvg_shape_get_stroke_width(const Tvg_Paint* paint, float* width);
// ShapeGetStrokeWidth 获取描边宽度
func (p *Paint) ShapeGetStrokeWidth() (float32, Result) {
	if p == nil || p.ptr == nil {
		return 0, ResultInvalidArgument
	}
	var cwidth C.float
	res := Result(C.tvg_shape_get_stroke_width(p.ptr, &cwidth))
	return float32(cwidth), res
}

// TVG_API Tvg_Result tvg_shape_set_stroke_color(Tvg_Paint* paint, uint8_t r, uint8_t g, uint8_t b, uint8_t a);
// ShapeSetStrokeColor 设置描边颜色
func (p *Paint) ShapeSetStrokeColor(r, g, b, a uint8) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_stroke_color(p.ptr, C.uint8_t(r), C.uint8_t(g), C.uint8_t(b), C.uint8_t(a)))
}

// TVG_API Tvg_Result tvg_shape_get_stroke_color(const Tvg_Paint* paint, uint8_t* r, uint8_t* g, uint8_t* b, uint8_t* a);
// ShapeGetStrokeColor 获取描边颜色
func (p *Paint) ShapeGetStrokeColor() (r, g, b, a uint8, res Result) {
	if p == nil || p.ptr == nil {
		return 0, 0, 0, 0, ResultInvalidArgument
	}
	var cr, cg, cb, ca C.uint8_t
	res = Result(C.tvg_shape_get_stroke_color(p.ptr, &cr, &cg, &cb, &ca))
	r = uint8(cr)
	g = uint8(cg)
	b = uint8(cb)
	a = uint8(ca)
	return
}

// TVG_API Tvg_Result tvg_shape_set_stroke_gradient(Tvg_Paint* paint, Tvg_Gradient* grad);
// ShapeSetStrokeGradient 设置描边渐变
func (p *Paint) ShapeSetStrokeGradient(grad *Gradient) Result {
	if p == nil || p.ptr == nil || grad == nil || grad.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_stroke_gradient(p.ptr, grad.ptr))
}

// TVG_API Tvg_Result tvg_shape_get_stroke_gradient(const Tvg_Paint* paint, Tvg_Gradient** grad);
func (p *Paint) ShapeGetStrokeGradient() (*Gradient, Result) {
	if p == nil || p.ptr == nil {
		return nil, ResultInvalidArgument
	}
	var cgrad *C.Tvg_Gradient
	res := Result(C.tvg_shape_get_stroke_gradient(p.ptr, &cgrad))
	if cgrad == nil {
		return nil, res
	}
	return &Gradient{ptr: cgrad}, res
}

// TVG_API Tvg_Result tvg_shape_set_stroke_dash(Tvg_Paint* paint, const float* dashPattern, uint32_t cnt, float offset);
// ShapeSetStrokeDash 设置描边虚线
func (p *Paint) ShapeSetStrokeDash(dashPattern []float32, cnt uint32, offset float32) Result {
	if p == nil || p.ptr == nil || len(dashPattern) < int(cnt) {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_stroke_dash(p.ptr, (*C.float)(unsafe.Pointer(&dashPattern[0])), C.uint32_t(cnt), C.float(offset)))
}

// TVG_API Tvg_Result tvg_shape_set_fill_color(Tvg_Paint* paint, uint8_t r, uint8_t g, uint8_t b, uint8_t a);
// ShapeSetFillColor 设置填充颜色
func (p *Paint) ShapeSetFillColor(r, g, b, a uint8) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_fill_color(p.ptr, C.uint8_t(r), C.uint8_t(g), C.uint8_t(b), C.uint8_t(a)))
}

// TVG_API Tvg_Result tvg_shape_get_fill_color(const Tvg_Paint* paint, uint8_t* r, uint8_t* g, uint8_t* b, uint8_t* a);
// ShapeGetFillColor 获取填充颜色
func (p *Paint) ShapeGetFillColor() (r, g, b, a uint8, res Result) {
	if p == nil || p.ptr == nil {
		return 0, 0, 0, 0, ResultInvalidArgument
	}
	var cr, cg, cb, ca C.uint8_t
	res = Result(C.tvg_shape_get_fill_color(p.ptr, &cr, &cg, &cb, &ca))
	r = uint8(cr)
	g = uint8(cg)
	b = uint8(cb)
	a = uint8(ca)
	return
}

// TVG_API Tvg_Result tvg_shape_set_fill_rule(Tvg_Paint* paint, Tvg_Fill_Rule rule);
// ShapeSetFillRule 设置填充规则
func (p *Paint) ShapeSetFillRule(rule FillRule) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_fill_rule(p.ptr, C.Tvg_Fill_Rule(rule)))
}

// TVG_API Tvg_Result tvg_shape_get_fill_rule(const Tvg_Paint* paint, Tvg_Fill_Rule* rule);
// ShapeGetFillRule 获取填充规则
func (p *Paint) ShapeGetFillRule() (*FillRule, Result) {
	if p == nil || p.ptr == nil {
		return nil, ResultInvalidArgument
	}
	var crule C.Tvg_Fill_Rule
	res := Result(C.tvg_shape_get_fill_rule(p.ptr, &crule))
	if res != ResultSuccess {
		return nil, res
	}
	rule := FillRule(crule)
	return &rule, res
}

// TVG_API Tvg_Result tvg_shape_set_paint_order(Tvg_Paint* paint, bool strokeFirst);
// ShapeSetPaintOrder 设置绘制顺序
func (p *Paint) ShapeSetPaintOrder(strokeFirst bool) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_paint_order(p.ptr, C.bool(strokeFirst)))
}

// TVG_API Tvg_Result tvg_shape_set_gradient(Tvg_Paint* paint, Tvg_Gradient* grad);
// ShapeSetGradient 设置渐变
func (p *Paint) ShapeSetGradient(grad *Gradient) Result {
	if p == nil || p.ptr == nil || grad == nil || grad.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_shape_set_gradient(p.ptr, grad.ptr))
}

// TVG_API Tvg_Result tvg_shape_get_gradient(const Tvg_Paint* paint, Tvg_Gradient** grad);
// ShapeGetGradient 获取渐变
func (p *Paint) ShapeGetGradient() (*Gradient, Result) {
	if p == nil || p.ptr == nil {
		return nil, ResultInvalidArgument
	}
	var cgrad *C.Tvg_Gradient
	res := Result(C.tvg_shape_get_gradient(p.ptr, &cgrad))
	if cgrad == nil {
		return nil, res
	}
	return &Gradient{ptr: cgrad}, res
}

/************************************************************************/
/* Gradient API                                                         */
/************************************************************************/

// TVG_API Tvg_Gradient* tvg_linear_gradient_new(void);
// LinearGradientNew 创建线性渐变
func LinearGradientNew() *Gradient {
	c := C.tvg_linear_gradient_new()
	if c == nil {
		return nil
	}
	return &Gradient{ptr: c}
}

// TVG_API Tvg_Gradient* tvg_radial_gradient_new(void);
// RadialGradientNew 创建径向渐变
func RadialGradientNew() *Gradient {
	c := C.tvg_radial_gradient_new()
	if c == nil {
		return nil
	}
	return &Gradient{ptr: c}
}

// TVG_API Tvg_Result tvg_linear_gradient_set(Tvg_Gradient* grad, float x1, float y1, float x2, float y2);
// LinearGradientSet 设置线性渐变
func (g *Gradient) LinearGradientSet(x1, y1, x2, y2 float32) Result {
	if g == nil || g.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_linear_gradient_set(g.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2)))
}

// TVG_API Tvg_Result tvg_linear_gradient_get(Tvg_Gradient* grad, float* x1, float* y1, float* x2, float* y2);
// LinearGradientGet 获取线性渐变
func (g *Gradient) LinearGradientGet() (x1, y1, x2, y2 float32, res Result) {
	if g == nil || g.ptr == nil {
		return 0, 0, 0, 0, ResultInvalidArgument
	}
	var cx1, cy1, cx2, cy2 C.float
	res = Result(C.tvg_linear_gradient_get(g.ptr, &cx1, &cy1, &cx2, &cy2))
	x1 = float32(cx1)
	y1 = float32(cy1)
	x2 = float32(cx2)
	y2 = float32(cy2)
	return
}

// TVG_API Tvg_Result tvg_radial_gradient_set(Tvg_Gradient* grad, float cx, float cy, float r, float fx, float fy, float fr);
// RadialGradientSet 设置径向渐变
func (g *Gradient) RadialGradientSet(cx, cy, r, fx, fy, fr float32) Result {
	if g == nil || g.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_radial_gradient_set(g.ptr, C.float(cx), C.float(cy), C.float(r), C.float(fx), C.float(fy), C.float(fr)))
}

// TVG_API Tvg_Result tvg_radial_gradient_get(Tvg_Gradient* grad, float* cx, float* cy, float* r, float* fx, float* fy, float* fr);
// RadialGradientGet 获取径向渐变
func (g *Gradient) RadialGradientGet() (cx, cy, r, fx, fy, fr float32, res Result) {
	if g == nil || g.ptr == nil {
		return 0, 0, 0, 0, 0, 0, ResultInvalidArgument
	}
	var ccx, ccy, cr, cfx, cfy, cfr C.float
	res = Result(C.tvg_radial_gradient_get(g.ptr, &ccx, &ccy, &cr, &cfx, &cfy, &cfr))
	cx = float32(ccx)
	cy = float32(ccy)
	r = float32(cr)
	fx = float32(cfx)
	fy = float32(cfy)
	fr = float32(cfr)
	return
}

// ColorStop 表示渐变中的一个颜色点
type ColorStop struct {
	Offset float32 // 颜色在渐变中的相对位置
	R      uint8   // 红色通道 [0,255]
	G      uint8   // 绿色通道 [0,255]
	B      uint8   // 蓝色通道 [0,255]
	A      uint8   // alpha通道 [0,255]
}

// Cptr 返回对应的 *C.Tvg_Color_Stop 指针
func (c *ColorStop) Cptr() *C.Tvg_Color_Stop {
	return (*C.Tvg_Color_Stop)(unsafe.Pointer(c))
}

// TVG_API Tvg_Result tvg_gradient_set_color_stops(Tvg_Gradient* grad, const Tvg_Color_Stop* color_stop, uint32_t cnt);
// SetColorStops 设置渐变颜色点
func (g *Gradient) SetColorStops(colorStops []ColorStop, cnt uint32) Result {
	if g == nil || g.ptr == nil || len(colorStops) < int(cnt) {
		return ResultInvalidArgument
	}
	return Result(C.tvg_gradient_set_color_stops(g.ptr, (*C.Tvg_Color_Stop)(unsafe.Pointer(&colorStops[0])), C.uint32_t(cnt)))
}

// TVG_API Tvg_Result tvg_gradient_get_color_stops(const Tvg_Gradient* grad, const Tvg_Color_Stop** color_stop, uint32_t* cnt);
// GetColorStops 获取渐变颜色点
func (g *Gradient) GetColorStops() (colorStops []ColorStop, res Result) {
	if g == nil || g.ptr == nil {
		return nil, ResultInvalidArgument
	}
	var ccolorStops *C.Tvg_Color_Stop
	var ccnt C.uint32_t
	res = Result(C.tvg_gradient_get_color_stops(g.ptr, &ccolorStops, &ccnt))
	if res != ResultSuccess {
		return nil, res
	}
	if ccnt > 0 {
		colorStops = unsafe.Slice((*ColorStop)(unsafe.Pointer(ccolorStops)), int(ccnt))
		colorStops = slices.Clone(colorStops)
	}
	return colorStops, res
}

// TVG_API Tvg_Result tvg_gradient_set_spread(Tvg_Gradient* grad, const Tvg_Stroke_Fill spread);
// SetSpread 设置渐变扩展模式
func (g *Gradient) SetSpread(spread StrokeFill) Result {
	if g == nil || g.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_gradient_set_spread(g.ptr, C.Tvg_Stroke_Fill(spread)))
}

// TVG_API Tvg_Result tvg_gradient_get_spread(const Tvg_Gradient* grad, Tvg_Stroke_Fill* spread);
// GetSpread 获取渐变扩展模式
func (g *Gradient) GetSpread() (*StrokeFill, Result) {
	if g == nil || g.ptr == nil {
		return nil, ResultInvalidArgument
	}
	var cspread C.Tvg_Stroke_Fill
	res := Result(C.tvg_gradient_get_spread(g.ptr, &cspread))
	if res != ResultSuccess {
		return nil, res
	}
	spread := StrokeFill(cspread)
	return &spread, res
}

// TVG_API Tvg_Result tvg_gradient_set_transform(Tvg_Gradient* grad, const Tvg_Matrix* m);
// SetTransform 设置渐变变换矩阵
func (g *Gradient) SetTransform(m Matrix) Result {
	if g == nil || g.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_gradient_set_transform(g.ptr, m.CPtr()))
}

// TVG_API Tvg_Result tvg_gradient_get_transform(const Tvg_Gradient* grad, Tvg_Matrix* m);
// GetTransform 获取渐变变换矩阵
func (g *Gradient) GetTransform() (Matrix, Result) {
	if g == nil || g.ptr == nil {
		return Matrix{}, ResultInvalidArgument
	}
	var cm C.Tvg_Matrix
	res := Result(C.tvg_gradient_get_transform(g.ptr, &cm))
	if res != ResultSuccess {
		return Matrix{}, res
	}
	return Matrix{
		E11: float32(cm.e11),
		E12: float32(cm.e12),
		E13: float32(cm.e13),
		E21: float32(cm.e21),
		E22: float32(cm.e22),
		E23: float32(cm.e23),
		E31: float32(cm.e31),
		E32: float32(cm.e32),
		E33: float32(cm.e33),
	}, res
}

// TVG_API Tvg_Result tvg_gradient_get_type(const Tvg_Gradient* grad, Tvg_Type* type);
// GetType 获取渐变类型
func (g *Gradient) GetType() (Type, Result) {
	if g == nil || g.ptr == nil {
		return TypeUndef, ResultInvalidArgument
	}
	var ctype C.Tvg_Type
	res := Result(C.tvg_gradient_get_type(g.ptr, &ctype))
	return Type(ctype), res
}

// TVG_API Tvg_Gradient* tvg_gradient_duplicate(Tvg_Gradient* grad);
// Duplicate 复制渐变
func (g *Gradient) Duplicate() *Gradient {
	if g == nil || g.ptr == nil {
		return nil
	}
	c := C.tvg_gradient_duplicate(g.ptr)
	if c == nil {
		return nil
	}
	return &Gradient{ptr: c}
}

// TVG_API Tvg_Result tvg_gradient_del(Tvg_Gradient* grad);
// Del 删除渐变
func (g *Gradient) Del() Result {
	if g == nil || g.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_gradient_del(g.ptr))
}

/************************************************************************/
/* Picture API                                                          */
/************************************************************************/

// TVG_API Tvg_Paint* tvg_picture_new(void);
func NewPicture() *Paint {
	c := C.tvg_picture_new()
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

// TVG_API Tvg_Result tvg_picture_load(Tvg_Paint* paint, const char* path);
// PictureLoad 加载图片
func (p *Paint) LoadPicture(path string) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	return Result(C.tvg_picture_load(p.ptr, cPath))
}

// TVG_API Tvg_Result tvg_picture_load_raw(Tvg_Paint* paint, uint32_t *data, uint32_t w, uint32_t h, Tvg_Colorspace cs, bool copy);
// LoadPictureRaw 加载原始图片数据
func (p *Paint) LoadPictureRaw(data []byte, w, h uint32, cs Colorspace, copy bool) Result {
	if p.ptr == nil {
		return ResultInvalidArgument
	}
	var cData *C.uint32_t
	if len(data) > 0 {
		cData = (*C.uint32_t)(unsafe.Pointer(&data[0]))
	} else {
		cData = nil
	}
	return Result(C.tvg_picture_load_raw(p.ptr, cData, C.uint32_t(w), C.uint32_t(h), C.Tvg_Colorspace(cs), C.bool(copy)))
}

// TVG_API Tvg_Result tvg_picture_load_data(Tvg_Paint* paint, const char *data, uint32_t size, const char *mimetype, const char* rpath, bool copy);
// LoadPictureData 加载图片数据
func (p *Paint) LoadPictureData(data []byte, size uint32, mimetype, rpath string, copy bool) Result {
	if p == nil || p.ptr == nil || len(data) < int(size) {
		return ResultInvalidArgument
	}
	cData := C.CString(string(data))
	defer C.free(unsafe.Pointer(cData))
	cMimetype := C.CString(mimetype)
	defer C.free(unsafe.Pointer(cMimetype))
	cRpath := C.CString(rpath)
	defer C.free(unsafe.Pointer(cRpath))
	return Result(C.tvg_picture_load_data(p.ptr, cData, C.uint32_t(size), cMimetype, cRpath, C.bool(copy)))
}

// TVG_API Tvg_Result tvg_picture_set_size(Tvg_Paint* paint, float w, float h);
// SetPictureSize 设置图片大小
func (p *Paint) SetPictureSize(w, h float32) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_picture_set_size(p.ptr, C.float(w), C.float(h)))
}

// TVG_API Tvg_Result tvg_picture_get_size(const Tvg_Paint* paint, float* w, float* h);
// GetPictureSize 获取图片大小
func (p *Paint) GetPictureSize() (w, h float32, res Result) {
	if p == nil || p.ptr == nil {
		return 0, 0, ResultInvalidArgument
	}
	var cw, ch C.float
	res = Result(C.tvg_picture_get_size(p.ptr, &cw, &ch))
	w = float32(cw)
	h = float32(ch)
	return
}

// TVG_API const Tvg_Paint* tvg_picture_get_paint(Tvg_Paint* paint, uint32_t id);
// GetPicturePaint 获取图片画布
func (p *Paint) GetPicturePaint(id uint32) *Paint {
	if p == nil || p.ptr == nil {
		return nil
	}
	c := C.tvg_picture_get_paint(p.ptr, C.uint32_t(id))
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

/************************************************************************/
/* Scene API                                                            */
/************************************************************************/
//TVG_API Tvg_Paint* tvg_scene_new(void);
func SceneNew() *Paint {
	c := C.tvg_scene_new()
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

// TVG_API Tvg_Result tvg_scene_push(Tvg_Paint* scene, Tvg_Paint* paint);
// ScenePush 将画布添加到场景中
func (s *Paint) ScenePush(paint *Paint) Result {
	if s == nil || s.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_scene_push(s.ptr, paint.ptr))
}

// TVG_API Tvg_Result tvg_scene_push_at(Tvg_Paint* scene, Tvg_Paint* target, Tvg_Paint* at);
// ScenePushAt 将画布添加到场景中的指定位置
func (s *Paint) ScenePushAt(target, at *Paint) Result {
	if s == nil || s.ptr == nil || target == nil || target.ptr == nil || at == nil || at.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_scene_push_at(s.ptr, target.ptr, at.ptr))
}

// TVG_API Tvg_Result tvg_scene_remove(Tvg_Paint* scene, Tvg_Paint* paint);
// SceneRemove 从场景中移除画布
func (s *Paint) SceneRemove(paint *Paint) Result {
	if s == nil || s.ptr == nil || paint == nil || paint.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_scene_remove(s.ptr, paint.ptr))
}

/************************************************************************/
/* Text API                                                            */
/************************************************************************/

// TVG_API Tvg_Paint* tvg_text_new(void);
// NewText 创建文本对象
func NewText() *Paint {
	c := C.tvg_text_new()
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

// TVG_API Tvg_Result tvg_text_set_font(Tvg_Paint* paint, const char* name, float size, const char* style);
// TextSetFont 设置字体
func (p *Paint) TextSetFont(name string, size float32, style string) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if name == "" {
		cName = nil
	}
	cStyle := C.CString(style)
	defer C.free(unsafe.Pointer(cStyle))
	return Result(C.tvg_text_set_font(p.ptr, cName, C.float(size), cStyle))
}

// TVG_API Tvg_Result tvg_text_set_text(Tvg_Paint* paint, const char* text);
// TextSetText 设置文本内容
func (p *Paint) TextSetText(text string) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return Result(C.tvg_text_set_text(p.ptr, cText))
}

// TVG_API Tvg_Result tvg_text_set_fill_color(Tvg_Paint* paint, uint8_t r, uint8_t g, uint8_t b);
// TextSetFillColor 设置填充颜色
func (p *Paint) TextSetFillColor(r, g, b uint8) Result {
	if p == nil || p.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_text_set_fill_color(p.ptr, C.uint8_t(r), C.uint8_t(g), C.uint8_t(b)))
}

// TVG_API Tvg_Result tvg_text_set_gradient(Tvg_Paint* paint, Tvg_Gradient* gradient);
// TextSetGradient 设置渐变
func (p *Paint) TextSetGradient(gradient *Gradient) Result {
	if p == nil || p.ptr == nil || gradient == nil || gradient.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_text_set_gradient(p.ptr, gradient.ptr))
}

// TVG_API Tvg_Result tvg_font_load(const char* path);
// LoadFont 加载字体
func LoadFont(path string) Result {
	if path == "" {
		return ResultInvalidArgument
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	return Result(C.tvg_font_load(cPath))
}

// TVG_API Tvg_Result tvg_font_load_data(const char* name, const char* data, uint32_t size, const char *mimetype, bool copy);
// LoadFontData 加载字体数据
func LoadFontData(name, data string, mimetype string, copy bool) Result {
	if name == "" || data == "" {
		return ResultInvalidArgument
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cData := C.CString(string(data))
	defer C.free(unsafe.Pointer(cData))
	cMimetype := C.CString(mimetype)
	defer C.free(unsafe.Pointer(cMimetype))
	return Result(C.tvg_font_load_data(cName, cData, C.uint32_t(len(data)), cMimetype, C.bool(copy)))
}

// TVG_API Tvg_Result tvg_font_unload(const char* path);
// UnloadFont 卸载字体
func UnloadFont(path string) Result {
	if path == "" {
		return ResultInvalidArgument
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	return Result(C.tvg_font_unload(cPath))
}

/************************************************************************/
/* Saver API                                                            */
/************************************************************************/
type Saver struct {
	ptr *C.Tvg_Saver
}

// TVG_API Tvg_Saver* tvg_saver_new(void);
func NewSaver() *Saver {
	c := C.tvg_saver_new()
	if c == nil {
		return nil
	}
	return &Saver{ptr: c}
}

// TVG_API Tvg_Result tvg_saver_save(Tvg_Saver* saver, Tvg_Paint* paint, const char* path, uint32_t quality);
// TODO：原本是上面的定义
// TVG_API Tvg_Result tvg_saver_save(Tvg_Saver* saver, Tvg_Animation* animation, const char* filename, uint32_t quality, uint32_t fps);
// Save 保存画布到文件
func (s *Saver) SaveAnimation(animation *Animation, filename string, quality, fps uint32) Result {
	if s == nil || s.ptr == nil || animation == nil || animation.ptr == nil {
		return ResultInvalidArgument
	}
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return Result(C.tvg_saver_save(s.ptr, animation.ptr, cFilename, C.uint32_t(quality), C.uint32_t(fps)))
}

// TVG_API Tvg_Result tvg_saver_sync(Tvg_Saver* saver);
// Sync 同步保存
func (s *Saver) Sync() Result {
	if s == nil || s.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_saver_sync(s.ptr))
}

// TVG_API Tvg_Result tvg_saver_del(Tvg_Saver* saver);
// Del 删除 Saver
func (s *Saver) Del() Result {
	if s == nil || s.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_saver_del(s.ptr))
}

/************************************************************************/
/* Animation API                                                        */
/************************************************************************/
type Animation struct {
	ptr *C.Tvg_Animation
}

// TVG_API Tvg_Animation* tvg_animation_new(void);
// NewAnimation 创建动画
func NewAnimation() *Animation {
	c := C.tvg_animation_new()
	if c == nil {
		return nil
	}
	return &Animation{ptr: c}
}

func (a *Animation) Unsafe() unsafe.Pointer {
	if a == nil || a.ptr == nil {
		return nil
	}
	return unsafe.Pointer(a.ptr)
}

// TVG_API Tvg_Result tvg_animation_set_frame(Tvg_Animation* animation, float no);
// SetFrame 设置动画帧
func (a *Animation) SetFrame(no float32) Result {
	if a == nil || a.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_animation_set_frame(a.ptr, C.float(no)))
}

// TVG_API Tvg_Paint* tvg_animation_get_picture(Tvg_Animation* animation);
// GetPicture 获取动画的图片
func (a *Animation) GetPicture() *Paint {
	if a == nil || a.ptr == nil {
		return nil
	}
	c := C.tvg_animation_get_picture(a.ptr)
	if c == nil {
		return nil
	}
	return &Paint{ptr: c}
}

// TVG_API Tvg_Result tvg_animation_get_frame(Tvg_Animation* animation, float* no);
// GetFrame 获取动画帧
func (a *Animation) GetFrame() (float32, Result) {
	if a == nil || a.ptr == nil {
		return 0, ResultInvalidArgument
	}
	var cno C.float
	res := Result(C.tvg_animation_get_frame(a.ptr, &cno))
	return float32(cno), res
}

// TVG_API Tvg_Result tvg_animation_get_total_frame(Tvg_Animation* animation, float* cnt);
// GetTotalFrame 获取动画总帧数
func (a *Animation) GetTotalFrame() (float32, Result) {
	if a == nil || a.ptr == nil {
		return 0, ResultInvalidArgument
	}
	var ccnt C.float
	res := Result(C.tvg_animation_get_total_frame(a.ptr, &ccnt))
	return float32(ccnt), res
}

// TVG_API Tvg_Result tvg_animation_get_duration(Tvg_Animation* animation, float* duration);
// GetDuration 获取动画持续时间
func (a *Animation) GetDuration() (float32, Result) {
	if a == nil || a.ptr == nil {
		return 0, ResultInvalidArgument
	}
	var cduration C.float
	res := Result(C.tvg_animation_get_duration(a.ptr, &cduration))
	return float32(cduration), res
}

// TVG_API Tvg_Result tvg_animation_set_segment(Tvg_Animation* animation, float begin, float end);
// SetSegment 设置动画段
func (a *Animation) SetSegment(begin, end float32) Result {
	if a == nil || a.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_animation_set_segment(a.ptr, C.float(begin), C.float(end)))
}

// TVG_API Tvg_Result tvg_animation_get_segment(Tvg_Animation* animation, float* begin, float* end);
// GetSegment 获取动画段
func (a *Animation) GetSegment() (begin, end float32, res Result) {
	if a == nil || a.ptr == nil {
		return 0, 0, ResultInvalidArgument
	}
	var cbeg, cend C.float
	res = Result(C.tvg_animation_get_segment(a.ptr, &cbeg, &cend))
	begin = float32(cbeg)
	end = float32(cend)
	return
}

// TVG_API Tvg_Result tvg_animation_del(Tvg_Animation* animation);
// Del 删除动画
func (a *Animation) Del() Result {
	if a == nil || a.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_animation_del(a.ptr))
}

/************************************************************************/
/* Accessor API                                                         */
/************************************************************************/
type Accessor struct {
	ptr *C.Tvg_Accessor
}

// TVG_API Tvg_Accessor* tvg_accessor_new();
func NewAccessor() *Accessor {
	c := C.tvg_accessor_new()
	if c == nil {
		return nil
	}
	return &Accessor{ptr: c}
}

// TVG_API Tvg_Result tvg_accessor_del(Tvg_Accessor* accessor);
func (a *Accessor) Del() Result {
	if a == nil || a.ptr == nil {
		return ResultInvalidArgument
	}
	return Result(C.tvg_accessor_del(a.ptr))
}

// TVG_API Tvg_Result tvg_accessor_set(Tvg_Accessor* accessor, Tvg_Paint* paint, bool (*func)(Tvg_Paint* paint, void* data), void* data);
// func (a *Accessor) Set(paint *Paint, f func(*Paint, unsafe.Pointer) bool, data unsafe.Pointer) Result {
// 	if a == nil || a.ptr == nil || paint == nil || paint.ptr == nil {
// 		return ResultInvalidArgument
// 	}
// 	return Result(C.tvg_accessor_set(a.ptr, paint.ptr, C.Tvg_Accessor_Func(unsafe.Pointer(f)), data))
// }

// TVG_API uint32_t tvg_accessor_generate_id(const char* name);
func (a *Accessor) GenerateID(name string) uint32 {
	if a == nil || a.ptr == nil {
		return 0
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return uint32(C.tvg_accessor_generate_id(cName))
}

// Destroy 销毁软件渲染 Canvas
func (c *Canvas) Destroy() Result {
	if c == nil || c.ptr == nil {
		return ResultMemoryCorruption
	}
	return Result(C.tvg_canvas_destroy(c.ptr))
}

// Point 表示二维空间中的点
type Point struct {
	X float32
	Y float32
}

// Cptr 返回对应的 *C.Tvg_Point 指针
func (p *Point) Cptr() *C.Tvg_Point {
	return (*C.Tvg_Point)(unsafe.Pointer(p))
}

// Tvg_Gradient
type Gradient struct {
	ptr *C.Tvg_Gradient
}
