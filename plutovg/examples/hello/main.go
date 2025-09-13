package main

import (
	"github.com/millken/go/plutovg"
)

func main() {
	// 创建画布和表面
	width, height := 400, 120
	surface := plutovg.SurfaceCreate(width, height)
	defer surface.Destroy()
	canvas := plutovg.CanvasCreate(surface)
	// canvas.Scale(float32(width)*1.5, float32(height)*1.5)
	defer canvas.Destroy()

	// 填充白色背景
	var bg = plutovg.NewColor()
	bg.InitRGB(1, 1, 1)
	// 或者使用 MakeColor 创建颜色
	bg = plutovg.MakeColor(1, 1, 1, 1)
	bg.Parse("rgb(182, 29, 29)")
	canvas.SetColor(bg)
	canvas.FillRect(0, 0, float32(width), float32(height))

	// 加载字体（假设有字体文件）
	font := plutovg.FontFaceLoadFromFile("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 0)
	if font == nil {
		panic("Font not found")
	}
	defer font.Destroy()
	canvas.SetFont(font, 36)

	// 设置黑色
	var fg plutovg.Color
	fg.Parse("rgb(255, 255, 0)")
	canvas.SetColor(&fg)

	// 绘制文本
	text := "Hello, World!"
	canvas.FillText(text, plutovg.TextEncoding(0), 40, 70)

	// 保存为 PNG
	if !surface.WriteToPNG("hello.png") {
		panic("Failed to write PNG")
	}
}
