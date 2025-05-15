package thorvg

import "unsafe"

type Result byte

const (
	ResultSuccess          Result = 0   // 正确执行
	ResultInvalidArgument  Result = 1   // 参数错误
	ResultInsufficientCond Result = 2   // 条件不足
	ResultFailedAllocation Result = 3   // 内存分配失败
	ResultMemoryCorruption Result = 4   // 内存损坏
	ResultNotSupported     Result = 5   // 不支持
	ResultUnknown          Result = 255 // 其他未知情况
)

type Colorspace byte

const (
	ColorspaceABGR8888  Colorspace = 0
	ColorspaceARGB8888  Colorspace = 1
	ColorspaceABGR8888S Colorspace = 2
	ColorspaceARGB8888S Colorspace = 3
	ColorspaceUnknown   Colorspace = 255
)

type MaskMethod byte

const (
	MaskNone         MaskMethod = 0  // 无遮罩
	MaskAlpha        MaskMethod = 1  // Alpha遮罩
	MaskInverseAlpha MaskMethod = 2  // 反Alpha遮罩
	MaskLuma         MaskMethod = 3  // 灰度遮罩
	MaskInverseLuma  MaskMethod = 4  // 反灰度遮罩
	MaskAdd          MaskMethod = 5  // 叠加
	MaskSubtract     MaskMethod = 6  // 相减
	MaskIntersect    MaskMethod = 7  // 交集
	MaskDifference   MaskMethod = 8  // 差异
	MaskLighten      MaskMethod = 9  // 取最大透明度
	MaskDarken       MaskMethod = 10 // 取最小透明度
)

type Type byte

const (
	TypeUndef      Type = 0  // 未定义类型
	TypeShape      Type = 1  // 形状类型
	TypeScene      Type = 2  // 场景类型
	TypePicture    Type = 3  // 图片类型
	TypeText       Type = 4  // 文本类型
	TypeLinearGrad Type = 10 // 线性渐变类型
	TypeRadialGrad Type = 11 // 放射渐变类型
)

type BlendMethod byte

const (
	BlendNormal     BlendMethod = 0  // 普通混合（默认）
	BlendMultiply   BlendMethod = 1  // 正片叠底
	BlendScreen     BlendMethod = 2  // 滤色
	BlendOverlay    BlendMethod = 3  // 叠加
	BlendDarken     BlendMethod = 4  // 变暗
	BlendLighten    BlendMethod = 5  // 变亮
	BlendColorDodge BlendMethod = 6  // 颜色减淡
	BlendColorBurn  BlendMethod = 7  // 颜色加深
	BlendHardLight  BlendMethod = 8  // 强光
	BlendSoftLight  BlendMethod = 9  // 柔光
	BlendDifference BlendMethod = 10 // 差值
	BlendExclusion  BlendMethod = 11 // 排除
	BlendHue        BlendMethod = 12 // 色相（保留，未支持）
	BlendSaturation BlendMethod = 13 // 饱和度（保留，未支持）
	BlendColor      BlendMethod = 14 // 颜色（保留，未支持）
	BlendLuminosity BlendMethod = 15 // 明度（保留，未支持）
	BlendAdd        BlendMethod = 16 // 线性加深
	BlendHardMix    BlendMethod = 17 // 硬混合（保留，未支持）
)

type PathCommand uint8

const (
	PathClose   PathCommand = 0 // 结束当前子路径并闭合
	PathMoveTo  PathCommand = 1 // 移动到新起点
	PathLineTo  PathCommand = 2 // 画直线到指定点
	PathCubicTo PathCommand = 3 // 画三次贝塞尔曲线到指定点
)

type FillRule byte

const (
	FillRuleNonZero FillRule = 0 // 非零环绕规则
	FillRuleEvenOdd FillRule = 1 // 奇偶规则
)

type StrokeFill int

const (
	StrokeFillPad     StrokeFill = 0 // 剩余区域用最近的渐变端点色填充
	StrokeFillReflect StrokeFill = 1 // 渐变区域外反射填充
	StrokeFillRepeat  StrokeFill = 2 // 渐变区域外重复填充
)

func toCStr(s string) uintptr {
	cstr := make([]byte, len(s)+1)
	copy(cstr, s)
	return uintptr(unsafe.Pointer(&cstr[0]))
}
