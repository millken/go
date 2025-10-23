package ui

import (
	"fmt"
	"testing"

	"github.com/millken/glay"
	"github.com/millken/go/plutovg"
)

func TestUI(t *testing.T) {
	width, height := 400, 120
	surface := plutovg.SurfaceCreate(width, height)
	defer surface.Destroy()
	canvas := plutovg.CanvasCreate(surface)
	// canvas.Scale(float32(width)*1.5, float32(height)*1.5)
	defer canvas.Destroy()
	var context glay.Context
	err := context.Initialize(glay.Config{
		Layout: glay.Dimensions{Width: 400, Height: 120},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = context.BeginLayout()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("生成一个UI界面，UI中间放置一个圆角按钮", func(t *testing.T) {
		// Try with explicit background rendering
		err = context.Clay(glay.ElementDeclaration{
			ID:              glay.ID("OuterContainer"),
			BackgroundColor: glay.Color{R: 43, G: 41, B: 51, A: 255},
			Layout: glay.LayoutConfig{
				LayoutDirection: glay.TopToBottom,
				Sizing: glay.Sizing{
					Width:  glay.SizingAxis{Type: glay.SizingGrow},
					Height: glay.SizingAxis{Type: glay.SizingGrow},
				},
				Padding:  glay.PaddingAll(16),
				ChildGap: 16,
				ChildAlignment: glay.ChildAlignment{
					X: glay.AlignXCenter,
					Y: glay.AlignYCenter,
				},
			},
		}, func() error {
			return context.Clay(glay.ElementDeclaration{
				ID:              glay.ID("Button"),
				BackgroundColor: glay.Color{R: 98, G: 102, B: 244, A: 255},
				CornerRadius: glay.CornerRadius{
					TopLeft:     8,
					TopRight:    8,
					BottomRight: 8,
					BottomLeft:  8,
				},
				Layout: glay.LayoutConfig{
					Sizing: glay.Sizing{
						Width:  glay.SizingAxis{Type: glay.SizingPercent, Percent: 0.3},
						Height: glay.SizingAxis{Type: glay.SizingPercent, Percent: 0.4},
					},
				},
			})
		})
		if err != nil {
			t.Fatal(err)
		}

		cmds, err := context.EndLayout()
		if err != nil {
			t.Fatal(err)
		}

		GlayRender(canvas, cmds)

		if !surface.WriteToPNG("ui_button_test.png") {
			t.Fatal("Failed to write PNG")
		}
	})
}

func TestStepByStep(t *testing.T) {
	width, height := 400, 120
	surface := plutovg.SurfaceCreate(width, height)
	defer surface.Destroy()
	canvas := plutovg.CanvasCreate(surface)
	defer canvas.Destroy()

	var context glay.Context
	err := context.Initialize(glay.Config{
		Layout: glay.Dimensions{Width: 400, Height: 120},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = context.BeginLayout()
	if err != nil {
		t.Fatal(err)
	}

	// Test: Container + Button, but container should be green to clearly distinguish
	err = context.Clay(glay.ElementDeclaration{
		ID:              glay.ID("GreenContainer"),
		BackgroundColor: glay.Color{R: 0, G: 255, B: 0, A: 255}, // Green container
		Layout: glay.LayoutConfig{
			LayoutDirection: glay.TopToBottom,
			Sizing: glay.Sizing{
				Width:  glay.SizingAxis{Type: glay.SizingGrow},
				Height: glay.SizingAxis{Type: glay.SizingGrow},
			},
			Padding: glay.PaddingAll(20), // Larger padding to see container
		},
	}, func() error {
		return context.Clay(glay.ElementDeclaration{
			ID:              glay.ID("RedButton"),
			BackgroundColor: glay.Color{R: 255, G: 0, B: 0, A: 255}, // Red button
			Layout: glay.LayoutConfig{
				Sizing: glay.Sizing{
					Width:  glay.SizingAxis{Type: glay.SizingFixed},
					Height: glay.SizingAxis{Type: glay.SizingFixed},
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	cmds, err := context.EndLayout()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Total commands: %d\n", len(cmds))
	GlayRender(canvas, cmds)

	if !surface.WriteToPNG("step_test.png") {
		t.Fatal("Failed to write PNG")
	}
}
