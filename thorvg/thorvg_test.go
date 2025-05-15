package thorvg_test

import (
	"io"
	"os"
	"testing"

	"github.com/dnsoa/go/assert"
	"github.com/millken/go/thorvg"
)

const (
	TEST_DIR = "./resources/"
)

func TestVersion(t *testing.T) {
	r := assert.New(t)
	major, minor, micro, version, res := thorvg.EngineVersion()
	t.Logf("Version: %d.%d.%d, %s", major, minor, micro, version)
	r.Equal(thorvg.ResultSuccess, res)
}

func TestPicture(t *testing.T) {
	r := assert.New(t)
	t.Run("Load RAW Data", func(t *testing.T) {
		file, err := os.Open(TEST_DIR + "/rawimage_200x300.raw")
		r.NoError(err)
		defer file.Close()
		data := make([]byte, 200*300)
		n, err := io.ReadFull(file, data)
		r.NoError(err)
		r.Equal(200*300, n)
		pic := thorvg.NewPicture()
		r.NotNil(pic)

		//Negative cases
		r.Equal(thorvg.ResultInvalidArgument, pic.LoadPictureRaw(nil, 200, 300, thorvg.ColorspaceABGR8888, false))
		r.Equal(thorvg.ResultInvalidArgument, pic.LoadPictureRaw(data, 0, 0, thorvg.ColorspaceABGR8888, false))
		r.Equal(thorvg.ResultInvalidArgument, pic.LoadPictureRaw(data, 200, 0, thorvg.ColorspaceUnknown, false))
		r.Equal(thorvg.ResultInvalidArgument, pic.LoadPictureRaw(data, 0, 300, thorvg.ColorspaceABGR8888, true))

		//Positive cases
		r.Equal(thorvg.ResultSuccess, pic.LoadPictureRaw(data, 200, 300, thorvg.ColorspaceABGR8888, false))
		r.Equal(thorvg.ResultSuccess, pic.LoadPictureRaw(data, 200, 300, thorvg.ColorspaceABGR8888, true))

		w, h, res := pic.GetPictureSize()
		r.Equal(thorvg.ResultSuccess, res)
		r.Equal(float32(200), w)
		r.Equal(float32(300), h)
	})

	t.Run("Load RAW file and render", func(t *testing.T) {
		r.Equal(thorvg.ResultSuccess, thorvg.EngineInit(0))
		canvas := thorvg.NewSwcanvas()
		r.NotNil(canvas)
		buffer := [100 * 100]uint32{}
		r.Equal(thorvg.ResultSuccess, canvas.SwcanvasSetTarget(buffer[:], 100, 100, 100, thorvg.ColorspaceABGR8888))
		file, err := os.Open(TEST_DIR + "/rawimage_200x300.raw")
		r.NoError(err)
		defer file.Close()
		data := make([]byte, 200*300)
		n, err := io.ReadFull(file, data)
		r.NoError(err)
		r.Equal(200*300, n)

		pic := thorvg.NewPicture()
		r.NotNil(pic)
		r.Equal(thorvg.ResultSuccess, pic.LoadPictureRaw(data, 200, 300, thorvg.ColorspaceABGR8888, false))
		r.Equal(thorvg.ResultSuccess, pic.SetPictureSize(100, 150))
		r.Equal(thorvg.ResultSuccess, canvas.Push(pic))
		data = nil
		r.Equal(thorvg.ResultSuccess, thorvg.EngineTerm()) //TODO:check why fail
	})
}

func TestText(t *testing.T) {
	r := assert.New(t)
	t.Run("Load SVG Data", func(t *testing.T) {
		text := thorvg.NewText()
		r.NotNil(text)
		tt, res := text.GetType()
		r.Equal(thorvg.ResultSuccess, res)
		r.Equal(thorvg.TypeText, tt)
	})
	t.Run("Load TTF Data from a file", func(t *testing.T) {
		thorvg.EngineInit(0)
		text := thorvg.NewText()
		r.NotNil(text)
		r.Equal(thorvg.ResultInsufficientCond, thorvg.UnloadFont(TEST_DIR+"/invalid.ttf"))
		r.Equal(thorvg.ResultSuccess, thorvg.LoadFont(TEST_DIR+"/Arial.ttf"))
		r.Equal(thorvg.ResultInvalidArgument, thorvg.LoadFont(TEST_DIR+"/invalid.ttf"))
		r.Equal(thorvg.ResultSuccess, thorvg.UnloadFont(TEST_DIR+"/Arial.ttf"))
		r.Equal(thorvg.ResultInvalidArgument, thorvg.LoadFont(""))
		r.Equal(thorvg.ResultSuccess, thorvg.LoadFont(TEST_DIR+"/NanumGothicCoding.ttf"))
		thorvg.EngineTerm()
	})

	t.Run("Load TTF Data from a buffer", func(t *testing.T) {
		thorvg.EngineInit(0)
		text := thorvg.NewText()
		r.NotNil(text)
		data, err := os.ReadFile(TEST_DIR + "/Arial.ttf")
		r.NoError(err)
		svg := "<svg height=\"1000\" viewBox=\"0 0 600 600\" ></svg>"

		//load
		r.Equal(thorvg.ResultInvalidArgument, thorvg.LoadFontData("", string(data), "", false))
		r.Equal(thorvg.ResultNotSupported, thorvg.LoadFontData("ArialSvg", svg, "unknown", false))
		thorvg.EngineTerm()
		r.Equal(thorvg.ResultSuccess, thorvg.LoadFontData("ArialUnknown", string(data), "unknown", false))
		r.Equal(thorvg.ResultSuccess, thorvg.LoadFontData("ArialTtf", string(data), "ttf", true))
		r.Equal(thorvg.ResultSuccess, thorvg.LoadFontData("Arial", string(data), "", false))

		//unload
		r.Equal(thorvg.ResultInsufficientCond, thorvg.UnloadFont("ArialUnknown"))
		r.Equal(thorvg.ResultInsufficientCond, thorvg.UnloadFont("ArialSvg"))
		thorvg.EngineTerm()
	})

	t.Run("Text Font", func(t *testing.T) {
		thorvg.EngineInit(0)
		text := thorvg.NewText()
		r.NotNil(text)
		r.Equal(thorvg.ResultSuccess, thorvg.LoadFont(TEST_DIR+"/Arial.ttf"))
		r.Equal(thorvg.ResultSuccess, text.TextSetFont("Arial", 80, ""))
		r.Equal(thorvg.ResultSuccess, text.TextSetFont("Arial", 1, ""))
		r.Equal(thorvg.ResultSuccess, text.TextSetFont("Arial", 50, ""))
		r.Equal(thorvg.ResultSuccess, text.TextSetFont("", 50, ""))
		r.Equal(thorvg.ResultInsufficientCond, text.TextSetFont("InvalidFont", 80, ""))
		thorvg.EngineTerm()
	})

	t.Run("Text Basic", func(t *testing.T) {
		thorvg.EngineInit(0)
		text := thorvg.NewText()
		r.NotNil(text)
		canvas := thorvg.NewSwcanvas()
		r.NotNil(canvas)
		buffer := [100 * 100]uint32{}
		r.Equal(thorvg.ResultSuccess, canvas.SwcanvasSetTarget(buffer[:], 100, 100, 100, thorvg.ColorspaceABGR8888))

		r.Equal(thorvg.ResultSuccess, thorvg.LoadFont(TEST_DIR+"/Arial.ttf"))
		r.Equal(thorvg.ResultSuccess, text.TextSetFont("Arial", 80, ""))
		r.Equal(thorvg.ResultSuccess, text.TextSetText(""))
		r.Equal(thorvg.ResultSuccess, text.TextSetText("ABCDEFGHIJIKLMOPQRSTUVWXYZ"))
		r.Equal(thorvg.ResultSuccess, text.TextSetText("THORVG Text"))
		r.Equal(thorvg.ResultSuccess, text.TextSetFillColor(uint8(255), uint8(255), uint8(255)))
		r.Equal(thorvg.ResultSuccess, canvas.Push(text))
	})
	t.Run("Text with composite glyphs", func(t *testing.T) {
		thorvg.EngineInit(0)
		text := thorvg.NewText()
		r.NotNil(text)
		canvas := thorvg.NewSwcanvas()
		r.NotNil(canvas)
		buffer := [100 * 100]uint32{}
		r.Equal(thorvg.ResultSuccess, canvas.SwcanvasSetTarget(buffer[:], 100, 100, 100, thorvg.ColorspaceABGR8888))

		r.Equal(thorvg.ResultSuccess, thorvg.LoadFont(TEST_DIR+"/Arial.ttf"))
		r.Equal(thorvg.ResultSuccess, text.TextSetFont("Arial", 80, ""))
		r.Equal(thorvg.ResultSuccess, text.TextSetText("\xc5\xbb\x6f\xc5\x82\xc4\x85\x64\xc5\xba \xc8\xab"))
		r.Equal(thorvg.ResultSuccess, text.TextSetFillColor(uint8(255), uint8(255), uint8(255)))
		r.Equal(thorvg.ResultSuccess, canvas.Push(text))
	})
}

func TestSaver(t *testing.T) {
	r := assert.New(t)
	t.Run("Saver Creation", func(t *testing.T) {
		saver := thorvg.NewSaver()
		r.NotNil(saver)
	})

	t.Run("Save a lottie into gif", func(t *testing.T) {
		thorvg.EngineInit(0)
		animation := thorvg.NewAnimation()
		r.NotNil(animation)

		picture := animation.GetPicture()
		r.Equal(thorvg.ResultSuccess, picture.LoadPicture(TEST_DIR+"test.json"))
		r.Equal(thorvg.ResultSuccess, picture.SetPictureSize(100, 100))
		saver := thorvg.NewSaver()
		r.Equal(thorvg.ResultSuccess, saver.SaveAnimation(animation, TEST_DIR+"test.gif", 0, 32))
		r.Equal(thorvg.ResultSuccess, saver.Sync())
	})
}
