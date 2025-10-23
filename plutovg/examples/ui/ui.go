// Package ui provides PlutoVG rendering functionality for Clay UI layouts.
package ui

import (
	"fmt"

	"github.com/millken/glay"
	"github.com/millken/go/plutovg"
)

func GlayRender(canvas *plutovg.Canvas, cmds []glay.RenderCommand) error {
	for _, cmd := range cmds {
		render(canvas, cmd)
	}
	return nil
}
func render(canvas *plutovg.Canvas, cmd glay.RenderCommand) {
	bbox := cmd.BoundingBox

	// Debug: Print render command info to understand layout results
	fmt.Printf("Render command: Type=%d, BBox={X:%.1f, Y:%.1f, W:%.1f, H:%.1f}\n",
		cmd.CommandType, bbox.X, bbox.Y, bbox.Width, bbox.Height) // NOTE: LayoutConfig (Sizing, Padding, ChildGap, etc.) is processed by Clay during layout calculation.
	// The render function receives the final computed bounding boxes and should only handle visual rendering.
	// Layout configuration affects:
	// 1. Element positioning (bbox.X, bbox.Y)
	// 2. Element sizing (bbox.Width, bbox.Height)
	// 3. Child element arrangement (handled by Clay's layout engine)

	switch cmd.CommandType {
	case glay.RenderCommandTypeRectangle:
		canvas.Save()
		canvas.Translate(float32(bbox.X), float32(bbox.Y))

		renderData := cmd.RenderData.(glay.RectangleRenderData)
		bc := renderData.BackgroundColor
		var bg plutovg.Color
		bg.InitRGBA(float32(bc.R)/255, float32(bc.G)/255, float32(bc.B)/255, float32(bc.A)/255)
		canvas.SetColor(&bg)

		cr := renderData.CornerRadius
		width := float32(bbox.Width)
		height := float32(bbox.Height)

		// Create rounded rectangle path using arcs (like Cairo implementation)
		canvas.NewPath()

		// Check if any corner has radius
		hasRoundedCorners := cr.TopLeft > 0 || cr.TopRight > 0 || cr.BottomRight > 0 || cr.BottomLeft > 0

		if !hasRoundedCorners {
			// Simple rectangle
			canvas.Rect(0, 0, width, height)
		} else {
			// Rounded rectangle using precise path construction
			rTl := float32(cr.TopLeft)
			rTr := float32(cr.TopRight)
			rBr := float32(cr.BottomRight)
			rBl := float32(cr.BottomLeft)

			// Start from top-left corner
			canvas.MoveTo(rTl, 0)

			// Top edge and top-right corner
			canvas.LineTo(width-rTr, 0)
			if rTr > 0 {
				canvas.Arc(width-rTr, rTr, rTr, plutovg.PI*1.5, plutovg.PI*2, false)
			}

			// Right edge and bottom-right corner
			canvas.LineTo(width, height-rBr)
			if rBr > 0 {
				canvas.Arc(width-rBr, height-rBr, rBr, 0, plutovg.PI*0.5, false)
			}

			// Bottom edge and bottom-left corner
			canvas.LineTo(rBl, height)
			if rBl > 0 {
				canvas.Arc(rBl, height-rBl, rBl, plutovg.PI*0.5, plutovg.PI, false)
			}

			// Left edge and top-left corner
			canvas.LineTo(0, rTl)
			if rTl > 0 {
				canvas.Arc(rTl, rTl, rTl, plutovg.PI, plutovg.PI*1.5, false)
			}

			canvas.ClosePath()
		}
		canvas.Fill()

		canvas.Restore()
	case glay.RenderCommandTypeBorder:
		canvas.Save()
		canvas.Translate(float32(bbox.X), float32(bbox.Y))

		renderData := cmd.RenderData.(glay.BorderRenderData)
		borderColor := renderData.Color
		var color plutovg.Color
		color.InitRGBA(float32(borderColor.R)/255, float32(borderColor.G)/255, float32(borderColor.B)/255, float32(borderColor.A)/255)
		canvas.SetColor(&color)

		cr := renderData.CornerRadius
		width := float32(bbox.Width)
		height := float32(bbox.Height)

		// Half radius for border arcs (like Cairo implementation)
		topLeftRadius := float32(cr.TopLeft) / 2.0
		topRightRadius := float32(cr.TopRight) / 2.0
		bottomRightRadius := float32(cr.BottomRight) / 2.0
		bottomLeftRadius := float32(cr.BottomLeft) / 2.0

		// Draw top border
		if renderData.Width.Top > 0 {
			canvas.SetLineWidth(float32(renderData.Width.Top))
			canvas.NewPath()

			// Left half-arc for top-left corner (225° to 270°)
			if topLeftRadius > 0 {
				canvas.Arc(topLeftRadius, topLeftRadius, topLeftRadius,
					plutovg.Deg2Rad(225), plutovg.Deg2Rad(270), false)
			} else {
				canvas.MoveTo(0, 0)
			}

			// Line to right half-arc
			canvas.LineTo(width-topRightRadius, 0)

			// Right half-arc for top-right corner (270° to 305°)
			if topRightRadius > 0 {
				canvas.Arc(width-topRightRadius, topRightRadius, topRightRadius,
					plutovg.Deg2Rad(270), plutovg.Deg2Rad(305), false)
			}

			canvas.Stroke()
		}

		// Draw right border
		if renderData.Width.Right > 0 {
			canvas.SetLineWidth(float32(renderData.Width.Right))
			canvas.NewPath()

			// Top half-arc for top-right corner (305° to 350°)
			if topRightRadius > 0 {
				canvas.Arc(width-topRightRadius, topRightRadius, topRightRadius,
					plutovg.Deg2Rad(305), plutovg.Deg2Rad(350), false)
			} else {
				canvas.MoveTo(width, 0)
			}

			// Line to bottom half-arc
			canvas.LineTo(width, height-bottomRightRadius)

			// Bottom half-arc for bottom-right corner (0° to 45°)
			if bottomRightRadius > 0 {
				canvas.Arc(width-bottomRightRadius, height-bottomRightRadius, bottomRightRadius,
					plutovg.Deg2Rad(0), plutovg.Deg2Rad(45), false)
			}

			canvas.Stroke()
		}

		// Draw bottom border
		if renderData.Width.Bottom > 0 {
			canvas.SetLineWidth(float32(renderData.Width.Bottom))
			canvas.NewPath()

			// Right half-arc for bottom-right corner (45° to 90°)
			if bottomRightRadius > 0 {
				canvas.Arc(width-bottomRightRadius, height-bottomRightRadius, bottomRightRadius,
					plutovg.Deg2Rad(45), plutovg.Deg2Rad(90), false)
			} else {
				canvas.MoveTo(width, height)
			}

			// Line to left half-arc
			canvas.LineTo(bottomLeftRadius, height)

			// Left half-arc for bottom-left corner (90° to 135°)
			if bottomLeftRadius > 0 {
				canvas.Arc(bottomLeftRadius, height-bottomLeftRadius, bottomLeftRadius,
					plutovg.Deg2Rad(90), plutovg.Deg2Rad(135), false)
			}

			canvas.Stroke()
		}

		// Draw left border
		if renderData.Width.Left > 0 {
			canvas.SetLineWidth(float32(renderData.Width.Left))
			canvas.NewPath()

			// Bottom half-arc for bottom-left corner (135° to 180°)
			if bottomLeftRadius > 0 {
				canvas.Arc(bottomLeftRadius, height-bottomLeftRadius, bottomLeftRadius,
					plutovg.Deg2Rad(135), plutovg.Deg2Rad(180), false)
			} else {
				canvas.MoveTo(0, height)
			}

			// Line to top half-arc
			canvas.LineTo(0, topLeftRadius)

			// Top half-arc for top-left corner (180° to 225°)
			if topLeftRadius > 0 {
				canvas.Arc(topLeftRadius, topLeftRadius, topLeftRadius,
					plutovg.Deg2Rad(180), plutovg.Deg2Rad(225), false)
			}

			canvas.Stroke()
		}

		canvas.Restore()
	case glay.RenderCommandTypeText:
		canvas.Save()
		canvas.Translate(float32(bbox.X), float32(bbox.Y))

		// Set default text color (black)
		var color plutovg.Color
		color.InitRGBA(0, 0, 0, 1)
		canvas.SetColor(&color)

		// Set font size (default)
		fontSize := float32(16)
		canvas.SetFontSize(fontSize)

		// Render text
		canvas.FillText("Button", 0, 0, 20) // Placeholder text with UTF-8 encoding

		canvas.Restore()
	case glay.RenderCommandTypeImage:
		canvas.Save()
		canvas.Translate(float32(bbox.X), float32(bbox.Y))

		// Load and render image if source is provided
		// Note: This would require implementing image loading
		// For now, just draw a placeholder rectangle
		var placeholderColor plutovg.Color
		placeholderColor.InitRGBA(0.8, 0.8, 0.8, 1.0) // Light gray
		canvas.SetColor(&placeholderColor)
		canvas.FillRect(0, 0, float32(bbox.Width), float32(bbox.Height))

		canvas.Restore()
	default:
		// Unknown command type - draw debug rectangle
		canvas.Save()
		canvas.Translate(float32(bbox.X), float32(bbox.Y))

		// Draw a debug outline for unknown elements
		var debugColor plutovg.Color
		debugColor.InitRGBA(1.0, 0.0, 1.0, 0.3) // Magenta with transparency
		canvas.SetColor(&debugColor)
		canvas.SetLineWidth(1.0)
		canvas.StrokeRect(0, 0, float32(bbox.Width), float32(bbox.Height))

		canvas.Restore()
	}
}
