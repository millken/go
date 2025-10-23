package main

/*
#cgo CFLAGS: -I/workspace/Codes/github.com/millken/go/sdl3/include
#cgo LDFLAGS: /workspace/Codes/github.com/millken/go/sdl3/libs/libSDL3.a /workspace/Codes/github.com/millken/go/sdl3/libs/libSDL3_ttf.a -lstdc++ -lm -pthread -ldl -lfreetype -lharfbuzz
// 如果不需要静态链接或平台不同，请根据平台调整 LDFLAGS（例如 macOS 可能需要 -framework CoreFoundation/其它 framework，或去掉 -pthread/-ldl）。
// 确保将静态库和头文件放在上面所示的 vendor 路径，或修改为你的实际路径。

#include <stdlib.h>
#include <SDL3/SDL.h>
#include <SDL3_ttf/SDL_ttf.h>

// 声明全局字体指针（去掉 static，以便 cgo 可以访问/调用辅助函数）
TTF_Font *globalFont = NULL;

// 初始化字体
int initFont(const char *fontPath, int fontSize) {
	if (!TTF_Init()) {
		return -1;
	}
    globalFont = TTF_OpenFont(fontPath, fontSize);
    return globalFont ? 0 : -2;
}

// 关闭字体并退出 TTF 子系统
void closeFont() {
    if (globalFont) {
        TTF_CloseFont(globalFont);
        globalFont = NULL;
    }
    TTF_Quit();
}

// 渲染文本到纹理
SDL_Texture* createTextTexture(SDL_Renderer *renderer, const char *text, SDL_Color fg, SDL_Color bg) {
    if (!globalFont) return NULL;

	// Use blended rendering so the surface has alpha and no solid background
	SDL_Surface *surface = TTF_RenderText_Blended(globalFont, text, 0, fg);
	if (!surface) return NULL;

	SDL_Texture *texture = SDL_CreateTextureFromSurface(renderer, surface);
	if (texture) {
		SDL_SetTextureBlendMode(texture, SDL_BLENDMODE_BLEND);
	}
    SDL_DestroySurface(surface);
    return texture;
}

// Helper: get texture size from SDL_Texture struct
int go_SDL_TextureSize(SDL_Texture *texture, int *w, int *h) {
	if (!texture) return -1;
	if (w) *w = texture->w;
	if (h) *h = texture->h;
	return 0;
}

// 包装：返回 SDL_RENDERER_PRESENTVSYNC 的值（有些 SDL3 宏/内联在 cgo 中解析有问题）

// 创建渲染器并设置是否启用 VSync（使用 SDL3 的 properties API）
SDL_Renderer* createRendererWithVSync(SDL_Window *window, int present_vsync) {
	SDL_PropertiesID props = SDL_CreateProperties();
	if (!props) return NULL;
	// 将窗口指针与属性关联
	SDL_SetPointerProperty(props, SDL_PROP_RENDERER_CREATE_WINDOW_POINTER, window);
	// 设置是否启用 VSync（非零表示启用）
	SDL_SetNumberProperty(props, SDL_PROP_RENDERER_CREATE_PRESENT_VSYNC_NUMBER, present_vsync);
	SDL_Renderer *renderer = SDL_CreateRendererWithProperties(props);
	return renderer;
}

// SDL3: SDL_QueryTexture may not exist; use go_SDL_TextureSize helper instead.

// 包装：SDL_RenderCopy
// Wrapper: call SDL_RenderTexture (SDL3) converting integer SDL_Rect to SDL_FRect
int go_SDL_RenderTexture(SDL_Renderer *renderer, SDL_Texture *texture, const SDL_Rect *srcrect, const SDL_Rect *dstrect) {
	SDL_FRect srcf, dstf;
	SDL_FRect *srcfp = NULL, *dstfp = NULL;
	if (srcrect) {
		srcf.x = srcrect->x;
		srcf.y = srcrect->y;
		srcf.w = srcrect->w;
		srcf.h = srcrect->h;
		srcfp = &srcf;
	}
	if (dstrect) {
		dstf.x = dstrect->x;
		dstf.y = dstrect->y;
		dstf.w = dstrect->w;
		dstf.h = dstrect->h;
		dstfp = &dstf;
	}
	return SDL_RenderTexture(renderer, texture, srcfp, dstfp) ? 1 : 0;
}

// Event helpers
int event_is_quit(SDL_Event *e) {
	return e->type == SDL_EVENT_QUIT;
}

int event_is_escape(SDL_Event *e) {
	if (e->type == SDL_EVENT_KEY_DOWN) {
		return e->key.key == SDLK_ESCAPE;
	}
	return 0;
}

int event_is_enter(SDL_Event *e) {
	if (e->type == SDL_EVENT_KEY_DOWN) {
		return e->key.key == SDLK_RETURN || e->key.key == SDLK_KP_ENTER;
	}
	return 0;
}
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func main() {
	// 初始化SDL
	if !bool(C.SDL_Init(C.SDL_INIT_VIDEO)) {
		panic("SDL初始化失败")
	}
	defer C.SDL_Quit()

	// 创建窗口（为调试去掉透明属性）
	title := C.CString("透明窗口")
	defer C.free(unsafe.Pointer(title))
	// SDL3 SDL_CreateWindow(signature): (const char *title, int w, int h, SDL_WindowFlags flags)
	window := C.SDL_CreateWindow(title, C.int(800), C.int(600), C.SDL_WINDOW_BORDERLESS)
	if window == nil {
		panic("窗口创建失败")
	}
	defer C.SDL_DestroyWindow(window)

	// 创建渲染器（使用 SDL3 properties API，启用 VSync）
	renderer := C.createRendererWithVSync(window, 1)
	if renderer == nil {
		panic("渲染器创建失败")
	}
	defer C.SDL_DestroyRenderer(renderer)

	// 设置渲染器混合模式（支持透明度）
	C.SDL_SetRenderDrawBlendMode(renderer, C.SDL_BLENDMODE_BLEND)

	// 初始化字体：尝试常见路径
	candidates := []string{
		"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
	}
	chosen := ""
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			chosen = p
			break
		}
	}
	if chosen == "" {
		fmt.Println("没有找到字体文件，文本渲染可能失败；请安装或修改字体路径")
		chosen = candidates[1]
	}
	fontPath := C.CString(chosen)
	defer C.free(unsafe.Pointer(fontPath))
	if C.initFont(fontPath, 48) != 0 {
		// print SDL errors for diagnosis
		fmt.Println("initFont failed, SDL error:", C.GoString(C.SDL_GetError()))
		panic("字体加载失败")
	}
	// 使用封装的 closeFont 清理字体并调用 TTF_Quit
	defer C.closeFont()

	// 创建文本纹理
	text := C.CString("你好，中国")
	defer C.free(unsafe.Pointer(text))
	fgColor := C.SDL_Color{r: 0, g: 0, b: 0, a: 255}       // 黑色文字
	bgColor := C.SDL_Color{r: 255, g: 255, b: 255, a: 255} // 白色背景
	texture := C.createTextTexture(renderer, text, fgColor, bgColor)
	if texture == nil {
		fmt.Println("createTextTexture returned NULL, SDL error:", C.GoString(C.SDL_GetError()))
		panic("文本纹理创建失败")
	}
	defer C.SDL_DestroyTexture(texture)

	// 获取文本尺寸
	var w, h C.int
	C.go_SDL_TextureSize(texture, &w, &h)

	// 主循环
	running := true
	for running {
		var event C.SDL_Event
		for C.SDL_PollEvent(&event) {
			if C.event_is_quit(&event) != 0 {
				running = false
			}
			if C.event_is_escape(&event) != 0 {
				running = false
			}
			if C.event_is_enter(&event) != 0 {
				running = false
			}
		}

		// 清空渲染器（完全透明）
		C.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 0)
		C.SDL_RenderClear(renderer)

		// 调试：绘制一个半透明红色矩形，验证混合是否生效
		var dbgRect C.SDL_FRect
		dbgRect.x = 10.0
		dbgRect.y = 10.0
		dbgRect.w = 200.0
		dbgRect.h = 100.0
		C.SDL_SetRenderDrawColor(renderer, 255, 0, 0, 128)
		C.SDL_RenderFillRect(renderer, &dbgRect)

		// 计算文本位置并使用整数 SDL_Rect 渲染（兼容性更好）
		rect := C.SDL_Rect{
			x: C.int((800 - int(w)) / 2),
			y: C.int((600 - int(h)) / 2),
			w: w,
			h: h,
		}

		// 渲染文本（使用包装的 SDL_RenderTexture）
		C.go_SDL_RenderTexture(renderer, texture, nil, &rect)

		// 更新屏幕
		C.SDL_RenderPresent(renderer)
		C.SDL_Delay(16)
	}
}
