package picasso_test

import (
	"image"
	"os"
	"testing"

	"github.com/dnsoa/go/assert"
	"github.com/millken/go/picasso"
)

const (
	TEST_WIDTH  = 640
	TEST_HEIGHT = 480
)

func TestVersion(t *testing.T) {
	r := assert.New(t)
	init := picasso.Initialize()
	r.Equal(picasso.True, init)
	picassoVersion := picasso.Version()
	r.Equal(int32(28000), picassoVersion)
	data := picasso.MakeByte(TEST_WIDTH * TEST_HEIGHT * 4)
	canvas := picasso.CanvasCreateWithData(data.Byte(), picasso.ColorFormatRgb, TEST_WIDTH, TEST_HEIGHT, TEST_WIDTH*4)
	r.NotNil(canvas)
	picasso.CanvasUnref(canvas)
	r.Equal(picasso.StatusSucceed, picasso.LastStatus())
}

var testImg *image.RGBA = image.NewRGBA(image.Rect(0, 0, TEST_WIDTH, TEST_HEIGHT))

func TestFlowers(t *testing.T) {
	r := assert.New(t)
	init := picasso.Initialize()
	r.Equal(picasso.True, init)
	data := picasso.MakeByte(TEST_WIDTH * TEST_HEIGHT * 4)
	canvas := picasso.CanvasCreateWithData(data.Byte(), picasso.ColorFormatRgba, TEST_WIDTH, TEST_HEIGHT, TEST_WIDTH*4)
	context := picasso.ContextCreate(canvas, nil)
	path := picasso.PathCreate()
	picasso.PathMoveTo(path, &picasso.Point{X: 0, Y: 0})
	picasso.PathLineTo(path, &picasso.Point{X: 100, Y: 100})
	picasso.PathLineTo(path, &picasso.Point{X: 200, Y: 0})
	picasso.PathLineTo(path, &picasso.Point{X: 300, Y: 100})
	picasso.PathLineTo(path, &picasso.Point{X: 400, Y: 0})
	picasso.PathSubClose(path)
	picasso.SetSourceColor(context, &picasso.Color{R: 1, G: 0, B: 0, A: 1})
	picasso.SetPath(context, path)
	picasso.Fill(context)
	picasso.PathUnref(path)
	picasso.ContextUnref(context)
	picasso.CanvasUnref(canvas)
	r.Equal(picasso.StatusSucceed, picasso.LastStatus())
	picasso.Shutdown()
	os.WriteFile("test.txt", data.Get(), os.ModePerm)
}
