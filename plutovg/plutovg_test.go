package plutovg

import (
	"fmt"
	"testing"
)

/*
在 PlutoVG（以及大多数 2D 图形库）中，**Surface** 和 **Canvas** 的关系如下：

---

### Surface

- **Surface**（表面）是一个像素缓冲区，通常对应一张图片或一块内存区域，用于存储最终的绘制结果（像素数据）。
- 你可以把 Surface 理解为“画布”或“纸张”，所有的绘制最终都体现在 Surface 上。
- Surface 负责像素的分配、释放、保存为 PNG/JPG 文件等。

---

### Canvas

- **Canvas**（画板/绘图上下文）是一个绘图环境，提供了各种绘图操作（如画线、画圆、填充、绘制文本等）。
- Canvas 负责解释你的绘图命令，并把结果绘制到它所关联的 Surface 上。
- Canvas 维护当前的绘图状态（如颜色、线宽、变换矩阵、当前路径等）。

---

### 它们的关系

- **Canvas 总是依附于一个 Surface**。你创建 Canvas 时，必须传入一个 Surface。
- 所有对 Canvas 的绘图操作，最终都会作用到它关联的 Surface 上。
- 一个 Surface 可以有多个 Canvas（但通常一个 Surface 只用一个 Canvas）。
- 你可以通过 `canvas.GetSurface()` 获取 Canvas 关联的 Surface。

---

### 代码示例

```go
surface := plutovg.SurfaceCreate(400, 300) // 创建像素缓冲区
canvas := plutovg.CanvasCreate(surface)    // 创建绘图上下文，绑定到 surface

canvas.SetRGB(1, 0, 0)
canvas.FillRect(0, 0, 100, 100)            // 实际像素写入 surface

surface.WriteToPNG("out.png")              // 保存 surface 的像素数据
```

---

**总结**：
- Surface 是像素存储区，负责“承载”图像数据
- Canvas 是绘图环境，负责“操作”这些像素
- Canvas 的所有操作都作用于它关联的 Surface

这样你就可以灵活地管理和复用像素数据与绘图状态。
*/
func TestColor(t *testing.T) {
	color := MakeColor(0.5, 0.5, 0.5, 1.0)
	if color.r != 0.5 || color.g != 0.5 || color.b != 0.5 || color.a != 1.0 {
		t.Errorf("Expected color (0.5, 0.5, 0.5, 1.0), got (%f, %f, %f, %f)", color.r, color.g, color.b, color.a)
	}

	str := "rgba(255, 0, 0, 1)"
	cp := color.Parse(str)
	if !cp {
		t.Errorf("Expected color parse to succeed, got %v", cp)
	}
	fmt.Printf("Parsed color: (%f, %f, %f, %f)\n", color.r, color.g, color.b, color.a)
}
