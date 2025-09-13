// Package plutovg provides Go bindings for the PlutoVG 2D vector graphics library.
// PlutoVG is a lightweight, single-header library for 2D vector graphics with anti-aliasing.
package plutovg

/*
#cgo CFLAGS: -I./include -DPLUTOVG_BUILD_STATIC
 #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/_libs/linux/amd64 -lplutovg -lm
#include "plutovg.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// Mathematical constants used in PlutoVG.
const (
	PI      = 3.14159265358979323846 // Pi constant
	TWO_PI  = 6.28318530717958647693 // 2 * Pi
	HALF_PI = 1.57079632679489661923 // Pi / 2
	SQRT2   = 1.41421356237309504880 // Square root of 2
	KAPPA   = 0.55228474983079339840 // Kappa constant for Bezier approximation of circles
)

const (
	PATH_COMMAND_MOVE_TO  PathCommand = 0
	PATH_COMMAND_LINE_TO  PathCommand = 1
	PATH_COMMAND_CUBIC_TO PathCommand = 2
	PATH_COMMAND_CLOSE    PathCommand = 3
)

// Basic structure types that correspond to PlutoVG C types.
type (
	Point        C.plutovg_point_t         // 2D point with x, y coordinates
	Rect         C.plutovg_rect_t          // Rectangle with x, y, width, height
	Matrix       C.plutovg_matrix_t        // 2D transformation matrix
	Color        C.plutovg_color_t         // RGBA color representation
	GradientStop C.plutovg_gradient_stop_t // Gradient stop with offset and color
	PathElement  C.plutovg_path_element_t  // Path element with command and points
	PathIterator C.plutovg_path_iterator_t // Iterator for path traversal
	TextIterator C.plutovg_text_iterator_t // Iterator for text processing
)

// Enumeration types that correspond to PlutoVG C enums.
type (
	PathCommand  = C.plutovg_path_command_t  // Path command type (move, line, curve, etc.)
	TextEncoding = C.plutovg_text_encoding_t // Text encoding type (UTF-8, UTF-16, etc.)
	TextureType  = C.plutovg_texture_type_t  // Texture type (plain, tiled)
	SpreadMethod = C.plutovg_spread_method_t // Gradient spread method (pad, reflect, repeat)
	FillRule     = C.plutovg_fill_rule_t     // Fill rule (non-zero, even-odd)
	Operator     = C.plutovg_operator_t      // Compositing operator
	LineCap      = C.plutovg_line_cap_t      // Line cap style (butt, round, square)
	LineJoin     = C.plutovg_line_join_t     // Line join style (miter, round, bevel)
	CodePoint    = C.plutovg_codepoint_t     // Unicode code point
)

// Opaque handle types for PlutoVG objects.
type (
	Path     struct{ ptr *C.plutovg_path_t }      // Vector path object
	FontFace struct{ ptr *C.plutovg_font_face_t } // Font face object
	Surface  struct{ ptr *C.plutovg_surface_t }   // Drawing surface (bitmap)
	Paint    struct{ ptr *C.plutovg_paint_t }     // Paint object (color, gradient, texture)
	Canvas   struct{ ptr *C.plutovg_canvas_t }    // Drawing context
)

var (
	IDENTITY_MATRIX = MakeMatrix(1, 0, 0, 1, 0, 0) // Identity matrix for transformations
)

// ================== Version ==================

// Version returns the version number of PlutoVG library as an integer.
func Version() int {
	return int(C.plutovg_version())
}

// VersionString returns the version string of PlutoVG library.
func VersionString() string {
	return C.GoString(C.plutovg_version_string())
}

// ================== Utils ==================

// Deg2Rad converts degrees to radians.
func Deg2Rad(degrees float32) float32 {
	return degrees * (PI / 180.0)
}

// Rad2Deg converts radians to degrees.
func Rad2Deg(radians float32) float32 {
	return radians * (180.0 / PI)
}

// MakePoint creates a new point with the specified x and y coordinates.
func MakePoint(x, y float32) *Point {
	point := C.plutovg_point_t{C.float(x), C.float(y)}
	return (*Point)(unsafe.Pointer(&point))
}

// MakeRect creates a new rectangle with the specified coordinates and dimensions.
func MakeRect(x, y, w, h float32) *Rect {
	rect := C.plutovg_rect_t{C.float(x), C.float(y), C.float(w), C.float(h)}
	return (*Rect)(unsafe.Pointer(&rect))
}

// MakeMatrix creates a new matrix with the specified values.
func MakeMatrix(a, b, c, d, e, f float32) *Matrix {
	matrix := C.plutovg_matrix_t{C.float(a), C.float(b), C.float(c), C.float(d), C.float(e), C.float(f)}
	return (*Matrix)(unsafe.Pointer(&matrix))
}

// MakeScale creates a new scaling matrix with the specified scale factors.
func MakeScale(sx, sy float32) *Matrix {
	matrix := C.plutovg_matrix_t{C.float(sx), 0, 0, C.float(sy), 0, 0}
	return (*Matrix)(unsafe.Pointer(&matrix))
}

// MakeTranslate creates a new translation matrix with the specified translation values.
func MakeTranslate(tx, ty float32) *Matrix {
	matrix := C.plutovg_matrix_t{1, 0, 0, 1, C.float(tx), C.float(ty)}
	return (*Matrix)(unsafe.Pointer(&matrix))
}

// MakeColor creates a new color with the specified RGBA values.
func MakeColor(r, g, b, a float32) *Color {
	color := C.plutovg_color_t{C.float(r), C.float(g), C.float(b), C.float(a)}
	return (*Color)(unsafe.Pointer(&color))
}

func NewColor() *Color {
	color := C.plutovg_color_t{0, 0, 0, 0} // Default color is black with full opacity
	return (*Color)(unsafe.Pointer(&color))
}

// ================== Matrix ==================

// Init initializes the matrix with the given values.
func (m *Matrix) Init(a, b, c, d, e, f float32) {
	C.plutovg_matrix_init((*C.plutovg_matrix_t)(m), C.float(a), C.float(b), C.float(c), C.float(d), C.float(e), C.float(f))
}

// InitIdentity initializes the matrix as an identity matrix.
func (m *Matrix) InitIdentity() {
	C.plutovg_matrix_init_identity((*C.plutovg_matrix_t)(m))
}

// InitTranslate initializes the matrix as a translation matrix.
func (m *Matrix) InitTranslate(tx, ty float32) {
	C.plutovg_matrix_init_translate((*C.plutovg_matrix_t)(m), C.float(tx), C.float(ty))
}

// InitScale initializes the matrix as a scaling matrix.
func (m *Matrix) InitScale(sx, sy float32) {
	C.plutovg_matrix_init_scale((*C.plutovg_matrix_t)(m), C.float(sx), C.float(sy))
}

// InitRotate initializes the matrix as a rotation matrix.
func (m *Matrix) InitRotate(angle float32) {
	C.plutovg_matrix_init_rotate((*C.plutovg_matrix_t)(m), C.float(angle))
}

// InitShear initializes the matrix as a shear transformation matrix.
func (m *Matrix) InitShear(shx, shy float32) {
	C.plutovg_matrix_init_shear((*C.plutovg_matrix_t)(m), C.float(shx), C.float(shy))
}

// Translate applies a translation transformation to the matrix.
func (m *Matrix) Translate(tx, ty float32) {
	C.plutovg_matrix_translate((*C.plutovg_matrix_t)(m), C.float(tx), C.float(ty))
}

// Scale applies a scaling transformation to the matrix.
func (m *Matrix) Scale(sx, sy float32) {
	C.plutovg_matrix_scale((*C.plutovg_matrix_t)(m), C.float(sx), C.float(sy))
}

// Rotate applies a rotation transformation to the matrix.
func (m *Matrix) Rotate(angle float32) {
	C.plutovg_matrix_rotate((*C.plutovg_matrix_t)(m), C.float(angle))
}

// Shear applies a shear transformation to the matrix.
func (m *Matrix) Shear(shx, shy float32) {
	C.plutovg_matrix_shear((*C.plutovg_matrix_t)(m), C.float(shx), C.float(shy))
}

// Multiply multiplies two matrices and stores the result in this matrix.
func (m *Matrix) Multiply(left, right *Matrix) {
	C.plutovg_matrix_multiply((*C.plutovg_matrix_t)(m), (*C.plutovg_matrix_t)(left), (*C.plutovg_matrix_t)(right))
}

// Invert inverts the matrix and stores the result in the inverse matrix.
func (m *Matrix) Invert(inverse *Matrix) bool {
	return bool(C.plutovg_matrix_invert((*C.plutovg_matrix_t)(m), (*C.plutovg_matrix_t)(inverse)))
}

// Map transforms the given point coordinates using this matrix.
func (m *Matrix) Map(x, y float32) (float32, float32) {
	var xx, yy C.float
	C.plutovg_matrix_map((*C.plutovg_matrix_t)(m), C.float(x), C.float(y), &xx, &yy)
	return float32(xx), float32(yy)
}

// MapPoint transforms a point using this matrix.
func (m *Matrix) MapPoint(src *Point, dst *Point) {
	C.plutovg_matrix_map_point((*C.plutovg_matrix_t)(m), (*C.plutovg_point_t)(src), (*C.plutovg_point_t)(dst))
}

// MapPoints transforms multiple points using this matrix.
func (m *Matrix) MapPoints(src, dst *Point, count int) {
	C.plutovg_matrix_map_points((*C.plutovg_matrix_t)(m), (*C.plutovg_point_t)(src), (*C.plutovg_point_t)(dst), C.int(count))
}

// MapRect transforms a rectangle using this matrix.
func (m *Matrix) MapRect(src *Rect, dst *Rect) {
	C.plutovg_matrix_map_rect((*C.plutovg_matrix_t)(m), (*C.plutovg_rect_t)(src), (*C.plutovg_rect_t)(dst))
}

// Parse parses a matrix from a string representation.
func (m *Matrix) Parse(data string, length int) bool {
	cstr := C.CString(data)
	defer C.free(unsafe.Pointer(cstr))
	return bool(C.plutovg_matrix_parse((*C.plutovg_matrix_t)(m), cstr, C.int(length)))
}

// ================== Path ==================

// PathCreate creates a new path object.
func PathCreate() *Path {
	p := C.plutovg_path_create()
	if p == nil {
		return nil
	}
	return &Path{ptr: p}
}

// Reference increases the reference count of the path.
func (p *Path) Reference() *Path {
	return &Path{ptr: C.plutovg_path_reference(p.ptr)}
}

// Destroy destroys the path and releases its resources.
func (p *Path) Destroy() {
	C.plutovg_path_destroy(p.ptr)
	p.ptr = nil
}

// GetReferenceCount returns the current reference count of the path.
func (p *Path) GetReferenceCount() int {
	return int(C.plutovg_path_get_reference_count(p.ptr))
}

// GetElements returns the path elements and their count.
func (p *Path) GetElements() ([]PathElement, int) {
	var elems *C.plutovg_path_element_t
	count := int(C.plutovg_path_get_elements(p.ptr, &elems))
	slice := (*[1 << 30]PathElement)(unsafe.Pointer(elems))[:count:count]
	return slice, count
}

// MoveTo moves the current point to the specified coordinates.
func (p *Path) MoveTo(x, y float32) {
	C.plutovg_path_move_to(p.ptr, C.float(x), C.float(y))
}

// LineTo adds a line from the current point to the specified coordinates.
func (p *Path) LineTo(x, y float32) {
	C.plutovg_path_line_to(p.ptr, C.float(x), C.float(y))
}

// QuadTo adds a quadratic Bézier curve to the path.
func (p *Path) QuadTo(x1, y1, x2, y2 float32) {
	C.plutovg_path_quad_to(p.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2))
}

// CubicTo adds a cubic Bézier curve to the path.
func (p *Path) CubicTo(x1, y1, x2, y2, x3, y3 float32) {
	C.plutovg_path_cubic_to(p.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2), C.float(x3), C.float(y3))
}

// ArcTo adds an elliptical arc to the path.
func (p *Path) ArcTo(rx, ry, angle float32, largeArc, sweep bool, x, y float32) {
	C.plutovg_path_arc_to(p.ptr, C.float(rx), C.float(ry), C.float(angle), C.bool(largeArc), C.bool(sweep), C.float(x), C.float(y))
}

// Close closes the current sub-path.
func (p *Path) Close() {
	C.plutovg_path_close(p.ptr)
}

// GetCurrentPoint returns the current point coordinates.
func (p *Path) GetCurrentPoint() (float32, float32) {
	var x, y C.float
	C.plutovg_path_get_current_point(p.ptr, &x, &y)
	return float32(x), float32(y)
}

// Reserve reserves space for the specified number of path elements.
func (p *Path) Reserve(count int) {
	C.plutovg_path_reserve(p.ptr, C.int(count))
}

// Reset resets the path to empty state.
func (p *Path) Reset() {
	C.plutovg_path_reset(p.ptr)
}

// AddRect adds a rectangle to the path.
func (p *Path) AddRect(x, y, w, h float32) {
	C.plutovg_path_add_rect(p.ptr, C.float(x), C.float(y), C.float(w), C.float(h))
}

// AddRoundRect adds a rounded rectangle to the path.
func (p *Path) AddRoundRect(x, y, w, h, rx, ry float32) {
	C.plutovg_path_add_round_rect(p.ptr, C.float(x), C.float(y), C.float(w), C.float(h), C.float(rx), C.float(ry))
}

// AddEllipse adds an ellipse to the path.
func (p *Path) AddEllipse(cx, cy, rx, ry float32) {
	C.plutovg_path_add_ellipse(p.ptr, C.float(cx), C.float(cy), C.float(rx), C.float(ry))
}

// AddCircle adds a circle to the path.
func (p *Path) AddCircle(cx, cy, r float32) {
	C.plutovg_path_add_circle(p.ptr, C.float(cx), C.float(cy), C.float(r))
}

// AddArc adds an arc to the path.
func (p *Path) AddArc(cx, cy, r, a0, a1 float32, ccw bool) {
	C.plutovg_path_add_arc(p.ptr, C.float(cx), C.float(cy), C.float(r), C.float(a0), C.float(a1), C.bool(ccw))
}

// AddPath adds another path to this path with optional transformation.
func (p *Path) AddPath(src *Path, matrix *Matrix) {
	C.plutovg_path_add_path(p.ptr, src.ptr, (*C.plutovg_matrix_t)(matrix))
}

// Transform applies a transformation matrix to the path.
func (p *Path) Transform(matrix *Matrix) {
	C.plutovg_path_transform(p.ptr, (*C.plutovg_matrix_t)(matrix))
}

// Clone creates a copy of the path.
func (p *Path) Clone() *Path {
	return &Path{ptr: C.plutovg_path_clone(p.ptr)}
}

// CloneFlatten creates a flattened copy of the path.
func (p *Path) CloneFlatten() *Path {
	return &Path{ptr: C.plutovg_path_clone_flatten(p.ptr)}
}

// CloneDashed creates a dashed copy of the path.
func (p *Path) CloneDashed(offset float32, dashes *float32, ndashes int) *Path {
	return &Path{ptr: C.plutovg_path_clone_dashed(p.ptr, C.float(offset), (*C.float)(dashes), C.int(ndashes))}
}

// Extents calculates the bounding rectangle of the path.
func (p *Path) Extents(extents *Rect, tight bool) float32 {
	return float32(C.plutovg_path_extents(p.ptr, (*C.plutovg_rect_t)(extents), C.bool(tight)))
}

// Length calculates the total length of the path.
func (p *Path) Length() float32 {
	return float32(C.plutovg_path_length(p.ptr))
}

// Parse parses a path from SVG path data string.
func (p *Path) Parse(data string, length int) bool {
	cstr := C.CString(data)
	defer C.free(unsafe.Pointer(cstr))
	return bool(C.plutovg_path_parse(p.ptr, cstr, C.int(length)))
}

// ================== FontFace ==================

// FontFaceLoadFromFile loads a font face from a file.
func FontFaceLoadFromFile(filename string, ttcindex int) *FontFace {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.plutovg_font_face_load_from_file(cstr, C.int(ttcindex))
	if ptr == nil {
		return nil
	}
	return &FontFace{ptr: ptr}
}

// FontFaceLoadFromData loads a font face from data.
func FontFaceLoadFromData(data unsafe.Pointer, length, ttcindex int, destroyFunc unsafe.Pointer, closure unsafe.Pointer) *FontFace {
	ptr := C.plutovg_font_face_load_from_data(data, C.uint(length), C.int(ttcindex), (*[0]byte)(destroyFunc), closure)
	if ptr == nil {
		return nil
	}
	return &FontFace{ptr: ptr}
}

// Reference increases the reference count of the font face.
func (f *FontFace) Reference() *FontFace {
	return &FontFace{ptr: C.plutovg_font_face_reference(f.ptr)}
}

// Destroy destroys the font face and releases its resources.
func (f *FontFace) Destroy() {
	C.plutovg_font_face_destroy(f.ptr)
	f.ptr = nil
}

// GetReferenceCount returns the current reference count of the font face.
func (f *FontFace) GetReferenceCount() int {
	return int(C.plutovg_font_face_get_reference_count(f.ptr))
}

// GetMetrics gets the font metrics for the specified size.
func (f *FontFace) GetMetrics(size float32, ascent, descent, lineGap *float32, extents *Rect) {
	C.plutovg_font_face_get_metrics(f.ptr, C.float(size), (*C.float)(ascent), (*C.float)(descent), (*C.float)(lineGap), (*C.plutovg_rect_t)(extents))
}

// GetGlyphMetrics gets the glyph metrics for the specified codepoint.
func (f *FontFace) GetGlyphMetrics(size float32, codepoint CodePoint, advanceWidth, leftSideBearing *float32, extents *Rect) {
	C.plutovg_font_face_get_glyph_metrics(f.ptr, C.float(size), codepoint, (*C.float)(advanceWidth), (*C.float)(leftSideBearing), (*C.plutovg_rect_t)(extents))
}

// GetGlyphPath gets the path for the specified glyph.
func (f *FontFace) GetGlyphPath(size, x, y float32, codepoint CodePoint, path *Path) float32 {
	return float32(C.plutovg_font_face_get_glyph_path(f.ptr, C.float(size), C.float(x), C.float(y), codepoint, path.ptr))
}

// TextExtents calculates the extents of the specified text.
func (f *FontFace) TextExtents(size float32, text unsafe.Pointer, length int, encoding TextEncoding, extents *Rect) float32 {
	return float32(C.plutovg_font_face_text_extents(f.ptr, C.float(size), text, C.int(length), encoding, (*C.plutovg_rect_t)(extents)))
}

// ================== Surface ==================

// SurfaceCreate creates a new surface with the specified width and height.
func SurfaceCreate(width, height int) *Surface {
	ptr := C.plutovg_surface_create(C.int(width), C.int(height))
	if ptr == nil {
		return nil
	}
	return &Surface{ptr: ptr}
}

// SurfaceCreateForData creates a surface using existing data buffer.
func SurfaceCreateForData(data unsafe.Pointer, width, height, stride int) *Surface {
	ptr := C.plutovg_surface_create_for_data((*C.uchar)(data), C.int(width), C.int(height), C.int(stride))
	if ptr == nil {
		return nil
	}
	return &Surface{ptr: ptr}
}

// SurfaceLoadFromImageFile loads a surface from an image file.
func SurfaceLoadFromImageFile(filename string) *Surface {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.plutovg_surface_load_from_image_file(cstr)
	if ptr == nil {
		return nil
	}
	return &Surface{ptr: ptr}
}

// SurfaceLoadFromImageData loads a surface from image data in memory.
func SurfaceLoadFromImageData(data unsafe.Pointer, length int) *Surface {
	ptr := C.plutovg_surface_load_from_image_data(data, C.int(length))
	if ptr == nil {
		return nil
	}
	return &Surface{ptr: ptr}
}

// SurfaceLoadFromImageBase64 loads a surface from base64 encoded image data.
func SurfaceLoadFromImageBase64(data string, length int) *Surface {
	cstr := C.CString(data)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.plutovg_surface_load_from_image_base64(cstr, C.int(length))
	if ptr == nil {
		return nil
	}
	return &Surface{ptr: ptr}
}

// Reference increases the reference count of the surface.
func (s *Surface) Reference() *Surface {
	return &Surface{ptr: C.plutovg_surface_reference(s.ptr)}
}

// Destroy destroys the surface and releases its resources.
func (s *Surface) Destroy() {
	C.plutovg_surface_destroy(s.ptr)
	s.ptr = nil
}

// GetReferenceCount returns the current reference count of the surface.
func (s *Surface) GetReferenceCount() int {
	return int(C.plutovg_surface_get_reference_count(s.ptr))
}

// GetData returns a pointer to the surface's pixel data.
func (s *Surface) GetData() unsafe.Pointer {
	return unsafe.Pointer(C.plutovg_surface_get_data(s.ptr))
}

// GetWidth returns the width of the surface in pixels.
func (s *Surface) GetWidth() int {
	return int(C.plutovg_surface_get_width(s.ptr))
}

// GetHeight returns the height of the surface in pixels.
func (s *Surface) GetHeight() int {
	return int(C.plutovg_surface_get_height(s.ptr))
}

// GetStride returns the stride (bytes per row) of the surface.
func (s *Surface) GetStride() int {
	return int(C.plutovg_surface_get_stride(s.ptr))
}

// Clear clears the surface with the specified color.
func (s *Surface) Clear(color *Color) {
	C.plutovg_surface_clear(s.ptr, (*C.plutovg_color_t)(color))
}

// WriteToPNG writes the surface to a PNG file.
func (s *Surface) WriteToPNG(filename string) bool {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	return bool(C.plutovg_surface_write_to_png(s.ptr, cstr))
}

// WriteToJPG writes the surface to a JPEG file with the specified quality.
func (s *Surface) WriteToJPG(filename string, quality int) bool {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	return bool(C.plutovg_surface_write_to_jpg(s.ptr, cstr, C.int(quality)))
}

// WriteToPNGStream writes the surface to a PNG stream using a callback function.
func (s *Surface) WriteToPNGStream(writeFunc unsafe.Pointer, closure unsafe.Pointer) bool {
	return bool(C.plutovg_surface_write_to_png_stream(s.ptr, (*[0]byte)(writeFunc), closure))
}

// WriteToJPGStream writes the surface to a JPEG stream using a callback function.
func (s *Surface) WriteToJPGStream(writeFunc unsafe.Pointer, closure unsafe.Pointer, quality int) bool {
	return bool(C.plutovg_surface_write_to_jpg_stream(s.ptr, (*[0]byte)(writeFunc), closure, C.int(quality)))
}

// ConvertARGBToRGBA converts pixel data from ARGB to RGBA format.
func ConvertARGBToRGBA(dst, src unsafe.Pointer, width, height, stride int) {
	C.plutovg_convert_argb_to_rgba((*C.uchar)(dst), (*C.uchar)(src), C.int(width), C.int(height), C.int(stride))
}

// ConvertRGBAToARGB converts pixel data from RGBA to ARGB format.
func ConvertRGBAToARGB(dst, src unsafe.Pointer, width, height, stride int) {
	C.plutovg_convert_rgba_to_argb((*C.uchar)(dst), (*C.uchar)(src), C.int(width), C.int(height), C.int(stride))
}

// ================== Color ==================

// InitRGB initializes a color with RGB values.
func (c *Color) InitRGB(r, g, b float32) {
	C.plutovg_color_init_rgb((*C.plutovg_color_t)(c), C.float(r), C.float(g), C.float(b))
}

// InitRGBA initializes a color with RGBA values.
func (c *Color) InitRGBA(r, g, b, a float32) {
	C.plutovg_color_init_rgba((*C.plutovg_color_t)(c), C.float(r), C.float(g), C.float(b), C.float(a))
}

// InitRGB8 initializes a color with 8-bit RGB values.
func (c *Color) InitRGB8(r, g, b int) {
	C.plutovg_color_init_rgb8((*C.plutovg_color_t)(c), C.int(r), C.int(g), C.int(b))
}

// InitRGBA8 initializes a color with 8-bit RGBA values.
func (c *Color) InitRGBA8(r, g, b, a int) {
	C.plutovg_color_init_rgba8((*C.plutovg_color_t)(c), C.int(r), C.int(g), C.int(b), C.int(a))
}

// InitRGBA32 initializes a color from a 32-bit RGBA value.
func (c *Color) InitRGBA32(value uint32) {
	C.plutovg_color_init_rgba32((*C.plutovg_color_t)(c), C.uint(value))
}

// InitARGB32 initializes a color from a 32-bit ARGB value.
func (c *Color) InitARGB32(value uint32) {
	C.plutovg_color_init_argb32((*C.plutovg_color_t)(c), C.uint(value))
}

// InitHSL initializes a color with HSL values.
func (c *Color) InitHSL(h, s, l float32) {
	C.plutovg_color_init_hsl((*C.plutovg_color_t)(c), C.float(h), C.float(s), C.float(l))
}

// InitHSLA initializes a color with HSLA values.
func (c *Color) InitHSLA(h, s, l, a float32) {
	C.plutovg_color_init_hsla((*C.plutovg_color_t)(c), C.float(h), C.float(s), C.float(l), C.float(a))
}

// ToRGBA32 converts the color to a 32-bit RGBA value.
func (c *Color) ToRGBA32() uint32 {
	return uint32(C.plutovg_color_to_rgba32((*C.plutovg_color_t)(c)))
}

// ToARGB32 converts the color to a 32-bit ARGB value.
func (c *Color) ToARGB32() uint32 {
	return uint32(C.plutovg_color_to_argb32((*C.plutovg_color_t)(c)))
}

// Parse parses a color from a string representation.
func (c *Color) Parse(str string) bool {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	return int(C.plutovg_color_parse((*C.plutovg_color_t)(c), cstr, C.int(len(str)))) == len(str)
}

// ================== Paint ==================

// PaintCreateRGB creates a paint with RGB color.
func PaintCreateRGB(r, g, b float32) *Paint {
	return &Paint{ptr: C.plutovg_paint_create_rgb(C.float(r), C.float(g), C.float(b))}
}

// PaintCreateRGBA creates a paint with RGBA color.
func PaintCreateRGBA(r, g, b, a float32) *Paint {
	return &Paint{ptr: C.plutovg_paint_create_rgba(C.float(r), C.float(g), C.float(b), C.float(a))}
}

// PaintCreateColor creates a paint from a color object.
func PaintCreateColor(color *Color) *Paint {
	return &Paint{ptr: C.plutovg_paint_create_color((*C.plutovg_color_t)(color))}
}

// PaintCreateLinearGradient creates a linear gradient paint.
func PaintCreateLinearGradient(x1, y1, x2, y2 float32, spread SpreadMethod, stops *GradientStop, nstops int, matrix *Matrix) *Paint {
	return &Paint{ptr: C.plutovg_paint_create_linear_gradient(C.float(x1), C.float(y1), C.float(x2), C.float(y2), spread, (*C.plutovg_gradient_stop_t)(stops), C.int(nstops), (*C.plutovg_matrix_t)(matrix))}
}

// PaintCreateRadialGradient creates a radial gradient paint.
func PaintCreateRadialGradient(cx, cy, cr, fx, fy, fr float32, spread SpreadMethod, stops *GradientStop, nstops int, matrix *Matrix) *Paint {
	return &Paint{ptr: C.plutovg_paint_create_radial_gradient(C.float(cx), C.float(cy), C.float(cr), C.float(fx), C.float(fy), C.float(fr), spread, (*C.plutovg_gradient_stop_t)(stops), C.int(nstops), (*C.plutovg_matrix_t)(matrix))}
}

// PaintCreateTexture creates a texture paint from a surface.
func PaintCreateTexture(surface *Surface, ttype TextureType, opacity float32, matrix *Matrix) *Paint {
	return &Paint{ptr: C.plutovg_paint_create_texture(surface.ptr, ttype, C.float(opacity), (*C.plutovg_matrix_t)(matrix))}
}

// Reference increases the reference count of the paint.
func (p *Paint) Reference() *Paint {
	return &Paint{ptr: C.plutovg_paint_reference(p.ptr)}
}

// Destroy destroys the paint and releases its resources.
func (p *Paint) Destroy() {
	C.plutovg_paint_destroy(p.ptr)
	p.ptr = nil
}

// GetReferenceCount returns the current reference count of the paint.
func (p *Paint) GetReferenceCount() int {
	return int(C.plutovg_paint_get_reference_count(p.ptr))
}

// ================== Canvas ==================

// CanvasCreate creates a new canvas with the specified surface.
func CanvasCreate(surface *Surface) *Canvas {
	ptr := C.plutovg_canvas_create(surface.ptr)
	if ptr == nil {
		return nil
	}
	return &Canvas{ptr: ptr}
}

// Reference increases the reference count of the canvas.
func (c *Canvas) Reference() *Canvas {
	return &Canvas{ptr: C.plutovg_canvas_reference(c.ptr)}
}

// Destroy destroys the canvas and releases its resources.
func (c *Canvas) Destroy() {
	C.plutovg_canvas_destroy(c.ptr)
	c.ptr = nil
}

// GetReferenceCount returns the current reference count of the canvas.
func (c *Canvas) GetReferenceCount() int {
	return int(C.plutovg_canvas_get_reference_count(c.ptr))
}

// GetSurface returns the surface associated with this canvas.
func (c *Canvas) GetSurface() *Surface {
	return &Surface{ptr: C.plutovg_canvas_get_surface(c.ptr)}
}

// Save pushes the current graphics state onto the state stack.
func (c *Canvas) Save() {
	C.plutovg_canvas_save(c.ptr)
}

// Restore pops the graphics state from the state stack.
func (c *Canvas) Restore() {
	C.plutovg_canvas_restore(c.ptr)
}

// SetRGB sets the current color using RGB values.
func (c *Canvas) SetRGB(r, g, b float32) {
	C.plutovg_canvas_set_rgb(c.ptr, C.float(r), C.float(g), C.float(b))
}

func (c *Canvas) SetHex(hex string) {

}

// SetRGBA sets the current color using RGBA values.
func (c *Canvas) SetRGBA(r, g, b, a float32) {
	C.plutovg_canvas_set_rgba(c.ptr, C.float(r), C.float(g), C.float(b), C.float(a))
}

// SetColor sets the current color using a color object.
func (c *Canvas) SetColor(color *Color) {
	C.plutovg_canvas_set_color(c.ptr, (*C.plutovg_color_t)(color))
}

// SetLinearGradient sets the current paint to a linear gradient.
func (c *Canvas) SetLinearGradient(x1, y1, x2, y2 float32, spread SpreadMethod, stops *GradientStop, nstops int, matrix *Matrix) {
	C.plutovg_canvas_set_linear_gradient(c.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2), spread, (*C.plutovg_gradient_stop_t)(stops), C.int(nstops), (*C.plutovg_matrix_t)(matrix))
}

// SetRadialGradient sets the current paint to a radial gradient.
func (c *Canvas) SetRadialGradient(cx, cy, cr, fx, fy, fr float32, spread SpreadMethod, stops *GradientStop, nstops int, matrix *Matrix) {
	C.plutovg_canvas_set_radial_gradient(c.ptr, C.float(cx), C.float(cy), C.float(cr), C.float(fx), C.float(fy), C.float(fr), spread, (*C.plutovg_gradient_stop_t)(stops), C.int(nstops), (*C.plutovg_matrix_t)(matrix))
}

// SetTexture sets the current paint to a texture.
func (c *Canvas) SetTexture(surface *Surface, ttype TextureType, opacity float32, matrix *Matrix) {
	C.plutovg_canvas_set_texture(c.ptr, surface.ptr, ttype, C.float(opacity), (*C.plutovg_matrix_t)(matrix))
}

// SetPaint sets the current paint.
func (c *Canvas) SetPaint(paint *Paint) {
	C.plutovg_canvas_set_paint(c.ptr, paint.ptr)
}

// GetPaint returns the current paint with the specified fallback color.
func (c *Canvas) GetPaint(color *Color) *Paint {
	return &Paint{ptr: C.plutovg_canvas_get_paint(c.ptr, (*C.plutovg_color_t)(color))}
}

// SetFont sets the current font face and size.
func (c *Canvas) SetFont(face *FontFace, size float32) {
	C.plutovg_canvas_set_font(c.ptr, face.ptr, C.float(size))
}

// SetFontFace sets the current font face.
func (c *Canvas) SetFontFace(face *FontFace) {
	C.plutovg_canvas_set_font_face(c.ptr, face.ptr)
}

// GetFontFace returns the current font face.
func (c *Canvas) GetFontFace() *FontFace {
	return &FontFace{ptr: C.plutovg_canvas_get_font_face(c.ptr)}
}

// SetFontSize sets the current font size.
func (c *Canvas) SetFontSize(size float32) {
	C.plutovg_canvas_set_font_size(c.ptr, C.float(size))
}

// GetFontSize returns the current font size.
func (c *Canvas) GetFontSize() float32 {
	return float32(C.plutovg_canvas_get_font_size(c.ptr))
}

// SetFillRule sets the current fill rule.
func (c *Canvas) SetFillRule(rule FillRule) {
	C.plutovg_canvas_set_fill_rule(c.ptr, rule)
}

// GetFillRule returns the current fill rule.
func (c *Canvas) GetFillRule() FillRule {
	return C.plutovg_canvas_get_fill_rule(c.ptr)
}

// SetOperator sets the current compositing operator.
func (c *Canvas) SetOperator(op Operator) {
	C.plutovg_canvas_set_operator(c.ptr, op)
}

// GetOperator returns the current compositing operator.
func (c *Canvas) GetOperator() Operator {
	return C.plutovg_canvas_get_operator(c.ptr)
}

// SetOpacity sets the current global opacity.
func (c *Canvas) SetOpacity(opacity float32) {
	C.plutovg_canvas_set_opacity(c.ptr, C.float(opacity))
}

// GetOpacity returns the current global opacity.
func (c *Canvas) GetOpacity() float32 {
	return float32(C.plutovg_canvas_get_opacity(c.ptr))
}

// SetLineWidth sets the current line width for stroking.
func (c *Canvas) SetLineWidth(w float32) {
	C.plutovg_canvas_set_line_width(c.ptr, C.float(w))
}

// GetLineWidth returns the current line width.
func (c *Canvas) GetLineWidth() float32 {
	return float32(C.plutovg_canvas_get_line_width(c.ptr))
}

// SetLineCap sets the current line cap style.
func (c *Canvas) SetLineCap(cap LineCap) {
	C.plutovg_canvas_set_line_cap(c.ptr, cap)
}

// GetLineCap returns the current line cap style.
func (c *Canvas) GetLineCap() LineCap {
	return C.plutovg_canvas_get_line_cap(c.ptr)
}

// SetLineJoin sets the current line join style.
func (c *Canvas) SetLineJoin(join LineJoin) {
	C.plutovg_canvas_set_line_join(c.ptr, join)
}

// GetLineJoin returns the current line join style.
func (c *Canvas) GetLineJoin() LineJoin {
	return C.plutovg_canvas_get_line_join(c.ptr)
}

// SetMiterLimit sets the current miter limit.
func (c *Canvas) SetMiterLimit(v float32) {
	C.plutovg_canvas_set_miter_limit(c.ptr, C.float(v))
}

// GetMiterLimit returns the current miter limit.
func (c *Canvas) GetMiterLimit() float32 {
	return float32(C.plutovg_canvas_get_miter_limit(c.ptr))
}

// SetDash sets the current dash pattern.
func (c *Canvas) SetDash(offset float32, dashes *float32, ndashes int) {
	C.plutovg_canvas_set_dash(c.ptr, C.float(offset), (*C.float)(dashes), C.int(ndashes))
}

// SetDashOffset sets the current dash offset.
func (c *Canvas) SetDashOffset(offset float32) {
	C.plutovg_canvas_set_dash_offset(c.ptr, C.float(offset))
}

// GetDashOffset returns the current dash offset.
func (c *Canvas) GetDashOffset() float32 {
	return float32(C.plutovg_canvas_get_dash_offset(c.ptr))
}

// SetDashArray sets the current dash array.
func (c *Canvas) SetDashArray(dashes *float32, ndashes int) {
	C.plutovg_canvas_set_dash_array(c.ptr, (*C.float)(dashes), C.int(ndashes))
}

// GetDashArray returns the current dash array.
func (c *Canvas) GetDashArray(dashes **float32) int {
	return int(C.plutovg_canvas_get_dash_array(c.ptr, (**C.float)(unsafe.Pointer(dashes))))
}

// Translate applies a translation transformation to the current matrix.
func (c *Canvas) Translate(tx, ty float32) {
	C.plutovg_canvas_translate(c.ptr, C.float(tx), C.float(ty))
}

// Scale applies a scaling transformation to the current matrix.
func (c *Canvas) Scale(sx, sy float32) {
	C.plutovg_canvas_scale(c.ptr, C.float(sx), C.float(sy))
}

// Shear applies a shear transformation to the current matrix.
func (c *Canvas) Shear(shx, shy float32) {
	C.plutovg_canvas_shear(c.ptr, C.float(shx), C.float(shy))
}

// Rotate applies a rotation transformation to the current matrix.
func (c *Canvas) Rotate(angle float32) {
	C.plutovg_canvas_rotate(c.ptr, C.float(angle))
}

// Transform applies a transformation matrix to the current matrix.
func (c *Canvas) Transform(matrix *Matrix) {
	C.plutovg_canvas_transform(c.ptr, (*C.plutovg_matrix_t)(matrix))
}

// ResetMatrix resets the current transformation matrix to identity.
func (c *Canvas) ResetMatrix() {
	C.plutovg_canvas_reset_matrix(c.ptr)
}

// SetMatrix sets the current transformation matrix.
func (c *Canvas) SetMatrix(matrix *Matrix) {
	C.plutovg_canvas_set_matrix(c.ptr, (*C.plutovg_matrix_t)(matrix))
}

// GetMatrix gets the current transformation matrix.
func (c *Canvas) GetMatrix(matrix *Matrix) {
	C.plutovg_canvas_get_matrix(c.ptr, (*C.plutovg_matrix_t)(matrix))
}

// Map transforms the given coordinates using the current matrix.
func (c *Canvas) Map(x, y float32) (float32, float32) {
	var xx, yy C.float
	C.plutovg_canvas_map(c.ptr, C.float(x), C.float(y), &xx, &yy)
	return float32(xx), float32(yy)
}

// MapPoint transforms a point using the current matrix.
func (c *Canvas) MapPoint(src, dst *Point) {
	C.plutovg_canvas_map_point(c.ptr, (*C.plutovg_point_t)(src), (*C.plutovg_point_t)(dst))
}

// MapRect transforms a rectangle using the current matrix.
func (c *Canvas) MapRect(src, dst *Rect) {
	C.plutovg_canvas_map_rect(c.ptr, (*C.plutovg_rect_t)(src), (*C.plutovg_rect_t)(dst))
}

// MoveTo moves the current point to the specified coordinates.
func (c *Canvas) MoveTo(x, y float32) {
	C.plutovg_canvas_move_to(c.ptr, C.float(x), C.float(y))
}

// LineTo adds a line from the current point to the specified coordinates.
func (c *Canvas) LineTo(x, y float32) {
	C.plutovg_canvas_line_to(c.ptr, C.float(x), C.float(y))
}

// QuadTo adds a quadratic Bézier curve to the current path.
func (c *Canvas) QuadTo(x1, y1, x2, y2 float32) {
	C.plutovg_canvas_quad_to(c.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2))
}

// CubicTo adds a cubic Bézier curve to the current path.
func (c *Canvas) CubicTo(x1, y1, x2, y2, x3, y3 float32) {
	C.plutovg_canvas_cubic_to(c.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2), C.float(x3), C.float(y3))
}

// ArcTo adds an elliptical arc to the current path.
func (c *Canvas) ArcTo(rx, ry, angle float32, largeArc, sweep bool, x, y float32) {
	C.plutovg_canvas_arc_to(c.ptr, C.float(rx), C.float(ry), C.float(angle), C.bool(largeArc), C.bool(sweep), C.float(x), C.float(y))
}

// Rect adds a rectangle to the current path.
func (c *Canvas) Rect(x, y, w, h float32) {
	C.plutovg_canvas_rect(c.ptr, C.float(x), C.float(y), C.float(w), C.float(h))
}

// RoundRect adds a rounded rectangle to the current path.
func (c *Canvas) RoundRect(x, y, w, h, rx, ry float32) {
	C.plutovg_canvas_round_rect(c.ptr, C.float(x), C.float(y), C.float(w), C.float(h), C.float(rx), C.float(ry))
}

// Ellipse adds an ellipse to the current path.
func (c *Canvas) Ellipse(cx, cy, rx, ry float32) {
	C.plutovg_canvas_ellipse(c.ptr, C.float(cx), C.float(cy), C.float(rx), C.float(ry))
}

// Circle adds a circle to the current path.
func (c *Canvas) Circle(cx, cy, r float32) {
	C.plutovg_canvas_circle(c.ptr, C.float(cx), C.float(cy), C.float(r))
}

// Arc adds an arc to the current path.
func (c *Canvas) Arc(cx, cy, r, a0, a1 float32, ccw bool) {
	C.plutovg_canvas_arc(c.ptr, C.float(cx), C.float(cy), C.float(r), C.float(a0), C.float(a1), C.bool(ccw))
}

// AddPath adds a path to the current path.
func (c *Canvas) AddPath(path *Path) {
	C.plutovg_canvas_add_path(c.ptr, path.ptr)
}

// NewPath clears the current path and starts a new path.
func (c *Canvas) NewPath() {
	C.plutovg_canvas_new_path(c.ptr)
}

// ClosePath closes the current subpath.
func (c *Canvas) ClosePath() {
	C.plutovg_canvas_close_path(c.ptr)
}

// GetCurrentPoint returns the current point coordinates.
func (c *Canvas) GetCurrentPoint() (float32, float32) {
	var x, y C.float
	C.plutovg_canvas_get_current_point(c.ptr, &x, &y)
	return float32(x), float32(y)
}

// GetPath returns a copy of the current path.
func (c *Canvas) GetPath() *Path {
	return &Path{ptr: C.plutovg_canvas_get_path(c.ptr)}
}

// FillExtents calculates the bounding box of the current path when filled.
func (c *Canvas) FillExtents(extents *Rect) {
	C.plutovg_canvas_fill_extents(c.ptr, (*C.plutovg_rect_t)(extents))
}

// StrokeExtents calculates the bounding box of the current path when stroked.
func (c *Canvas) StrokeExtents(extents *Rect) {
	C.plutovg_canvas_stroke_extents(c.ptr, (*C.plutovg_rect_t)(extents))
}

// ClipExtents calculates the current clip region extents.
func (c *Canvas) ClipExtents(extents *Rect) {
	C.plutovg_canvas_clip_extents(c.ptr, (*C.plutovg_rect_t)(extents))
}

// Fill fills the current path with the current paint.
func (c *Canvas) Fill() {
	C.plutovg_canvas_fill(c.ptr)
}

// Stroke strokes the current path with the current paint.
func (c *Canvas) Stroke() {
	C.plutovg_canvas_stroke(c.ptr)
}

// Clip clips the current region to the current path.
func (c *Canvas) Clip() {
	C.plutovg_canvas_clip(c.ptr)
}

// Paint paints the current paint to the entire canvas.
func (c *Canvas) Paint() {
	C.plutovg_canvas_paint(c.ptr)
}

// FillPreserve fills the current path with the current paint and preserves the path.
func (c *Canvas) FillPreserve() {
	C.plutovg_canvas_fill_preserve(c.ptr)
}

// StrokePreserve strokes the current path with the current paint and preserves the path.
func (c *Canvas) StrokePreserve() {
	C.plutovg_canvas_stroke_preserve(c.ptr)
}

// ClipPreserve clips the current region to the current path and preserves the path.
func (c *Canvas) ClipPreserve() {
	C.plutovg_canvas_clip_preserve(c.ptr)
}

// FillRect fills a rectangle with the current paint.
func (c *Canvas) FillRect(x, y, w, h float32) {
	C.plutovg_canvas_fill_rect(c.ptr, C.float(x), C.float(y), C.float(w), C.float(h))
}

// FillPath fills a path with the current paint.
func (c *Canvas) FillPath(path *Path) {
	C.plutovg_canvas_fill_path(c.ptr, path.ptr)
}

// StrokeRect strokes a rectangle with the current paint.
func (c *Canvas) StrokeRect(x, y, w, h float32) {
	C.plutovg_canvas_stroke_rect(c.ptr, C.float(x), C.float(y), C.float(w), C.float(h))
}

// StrokePath strokes a path with the current paint.
func (c *Canvas) StrokePath(path *Path) {
	C.plutovg_canvas_stroke_path(c.ptr, path.ptr)
}

// ClipRect clips the current region to a rectangle.
func (c *Canvas) ClipRect(x, y, w, h float32) {
	C.plutovg_canvas_clip_rect(c.ptr, C.float(x), C.float(y), C.float(w), C.float(h))
}

// ClipPath clips the current region to a path.
func (c *Canvas) ClipPath(path *Path) {
	C.plutovg_canvas_clip_path(c.ptr, path.ptr)
}

// AddGlyph adds a glyph to the current path.
func (c *Canvas) AddGlyph(codepoint CodePoint, x, y float32) float32 {
	return float32(C.plutovg_canvas_add_glyph(c.ptr, codepoint, C.float(x), C.float(y)))
}

// AddText adds text to the current path.
func (c *Canvas) AddText(text unsafe.Pointer, length int, encoding TextEncoding, x, y float32) float32 {
	return float32(C.plutovg_canvas_add_text(c.ptr, text, C.int(length), encoding, C.float(x), C.float(y)))
}

// FillText fills text with the current paint.
func (c *Canvas) FillText(text string, encoding TextEncoding, x, y float32) float32 {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	return float32(C.plutovg_canvas_fill_text(c.ptr, unsafe.Pointer(cstr), C.int(len(text)), encoding, C.float(x), C.float(y)))
}

// StrokeText strokes text with the current paint.
func (c *Canvas) StrokeText(text unsafe.Pointer, length int, encoding TextEncoding, x, y float32) float32 {
	return float32(C.plutovg_canvas_stroke_text(c.ptr, text, C.int(length), encoding, C.float(x), C.float(y)))
}

// ClipText clips the current region to text.
func (c *Canvas) ClipText(text unsafe.Pointer, length int, encoding TextEncoding, x, y float32) float32 {
	return float32(C.plutovg_canvas_clip_text(c.ptr, text, C.int(length), encoding, C.float(x), C.float(y)))
}

// FontMetrics gets the current font metrics.
func (c *Canvas) FontMetrics(ascent, descent, lineGap *float32, extents *Rect) {
	C.plutovg_canvas_font_metrics(c.ptr, (*C.float)(ascent), (*C.float)(descent), (*C.float)(lineGap), (*C.plutovg_rect_t)(extents))
}

// GlyphMetrics gets the metrics for a specific glyph.
func (c *Canvas) GlyphMetrics(codepoint CodePoint, advanceWidth, leftSideBearing *float32, extents *Rect) {
	C.plutovg_canvas_glyph_metrics(c.ptr, codepoint, (*C.float)(advanceWidth), (*C.float)(leftSideBearing), (*C.plutovg_rect_t)(extents))
}

// TextExtents calculates the extents of text.
func (c *Canvas) TextExtents(text unsafe.Pointer, length int, encoding TextEncoding, extents *Rect) float32 {
	return float32(C.plutovg_canvas_text_extents(c.ptr, text, C.int(length), encoding, (*C.plutovg_rect_t)(extents)))
}
