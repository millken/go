package main

import (
	"github.com/millken/go/plutovg"
)

func main() {
	width := 150
	height := 150

	centerX := float32(width) / 2
	centerY := float32(height) / 2
	faceRadius := float32(70)
	mouthRadius := float32(50)
	eyeRadius := float32(10)
	eyeOffsetX := float32(25)
	eyeOffsetY := float32(20)
	eyeX := centerX - eyeOffsetX
	eyeY := centerY - eyeOffsetY

	surface := plutovg.SurfaceCreate(width, height)
	defer surface.Destroy()

	canvas := plutovg.CanvasCreate(surface)
	defer canvas.Destroy()

	// 画脸
	canvas.Save()
	canvas.Arc(centerX, centerY, faceRadius, 0, plutovg.TWO_PI, false)
	canvas.SetRGB(1, 1, 0)
	canvas.FillPreserve()
	canvas.SetRGB(0, 0, 0)
	canvas.SetLineWidth(1)
	canvas.Stroke()
	canvas.Restore()

	// 画眼睛
	canvas.Save()
	canvas.Arc(eyeX, eyeY, eyeRadius, 0, plutovg.TWO_PI, false)
	canvas.Arc(centerX+eyeOffsetX, eyeY, eyeRadius, 0, plutovg.TWO_PI, false)
	canvas.SetRGB(0, 0, 0)
	canvas.Fill()
	canvas.Restore()

	// 画嘴
	canvas.Save()
	canvas.Arc(centerX, centerY, mouthRadius, 0, plutovg.PI, false)
	canvas.SetRGB(0, 0, 0)
	canvas.SetLineWidth(5)
	canvas.Stroke()
	canvas.Restore()

	// 保存PNG
	canvas.GetSurface().WriteToPNG("smiley.png")
}
