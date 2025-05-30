NANOVG_SRC_DIR = ../nanovgXC/src

# 临时C文件名，用于整合NanoVG实现
NANOVG_IMPL_C = nanovg_combined_impl.c

# 选择 NanoVG 的 OpenGL 后端实现宏
# 可选项: NANOVG_GL2_IMPLEMENTATION, NANOVG_GL3_IMPLEMENTATION,
#         NANOVG_GLES2_IMPLEMENTATION, NANOVG_GLES3_IMPLEMENTATION
# NANOVG_GL3_IMPLEMENTATION 适用于较新的桌面 OpenGL
NANOVG_GL_BACKEND_DEFINE = NANOVG_GL3_IMPLEMENTATION

ifeq ($(NANOVG_GL_BACKEND_DEFINE),NANOVG_GL3_IMPLEMENTATION)
NANOVG_GL_HEADER = <OpenGL/gl3.h>
# 如果使用 GL3，通常还需要包含 gl3ext.h
NANOVG_GL_EXT_HEADER = <OpenGL/gl3ext.h>
else ifeq ($(NANOVG_GL_BACKEND_DEFINE),NANOVG_GL2_IMPLEMENTATION)
NANOVG_GL_HEADER = <OpenGL/gl.h>
NANOVG_GL_EXT_HEADER = # GL2 通常不需要额外的 ext 头文件
else ifeq ($(NANOVG_GL_BACKEND_DEFINE),NANOVG_GLES2_IMPLEMENTATION)
NANOVG_GL_HEADER = <OpenGLES/ES2/gl.h>
NANOVG_GL_EXT_HEADER = <OpenGLES/ES2/glext.h>
else ifeq ($(NANOVG_GL_BACKEND_DEFINE),NANOVG_GLES3_IMPLEMENTATION)
NANOVG_GL_HEADER = <OpenGLES/ES3/gl.h>
NANOVG_GL_EXT_HEADER = <OpenGLES/ES3/glext.h>
else
$(error Unknown NANOVG_GL_BACKEND_DEFINE: $(NANOVG_GL_BACKEND_DEFINE))
endif


test_nanovg:
	clang++ -arch arm64 -I../nanovgXC/src test_nanovg.c -L. -o test_nanovg -framework OpenGL -framework Foundation


# 通用编译参数
NANOVG_CFLAGS_COMMON = -fPIC -std=c99 -O2 -I$(NANOVG_SRC_DIR)

# 为 macOS arm64 架构编译 NanoVG 静态库
libnanovg_darwin_arm64.a: $(NANOVG_SRC_DIR)/nanovg.h $(NANOVG_SRC_DIR)/nanovg_gl.h
	@echo "Building NanoVG static library for darwin_arm64..."
	@echo '#define NANOVG_IMPLEMENTATION' > $(NANOVG_IMPL_C)
	@echo '#include "nanovg.h"' >> $(NANOVG_IMPL_C)
	@echo "" >> $(NANOVG_IMPL_C)
	@echo '// Include necessary OpenGL headers before nanovg_gl.h' >> $(NANOVG_IMPL_C)
	@echo '#include $(NANOVG_GL_HEADER)' >> $(NANOVG_IMPL_C)
ifeq ($(NANOVG_GL_BACKEND_DEFINE),NANOVG_GL3_IMPLEMENTATION)
	@echo '#define GL_SILENCE_DEPRECATION // For macOS to silence GL deprecation warnings' >> $(NANOVG_IMPL_C)
	@echo '#include $(NANOVG_GL_EXT_HEADER)' >> $(NANOVG_IMPL_C)
else ifneq ($(NANOVG_GL_EXT_HEADER),)
	@echo '#include $(NANOVG_GL_EXT_HEADER)' >> $(NANOVG_IMPL_C)
endif
	@echo "" >> $(NANOVG_IMPL_C)
	@echo '#define $(NANOVG_GL_BACKEND_DEFINE)' >> $(NANOVG_IMPL_C)
	@echo '#include "nanovg_gl.h"' >> $(NANOVG_IMPL_C)
	# 如果你使用 nanovg_gl_utils.h 并且它需要一个实现宏:
	# @echo '#define NANOVG_GL_UTILS_IMPLEMENTATION_MACRO_HERE' >> $(NANOVG_IMPL_C)
	# @echo '#include "nanovg_gl_utils.h"' >> $(NANOVG_IMPL_C)

	clang $(NANOVG_CFLAGS_COMMON) -arch arm64 -c $(NANOVG_IMPL_C) -o nanovg_darwin_arm64.o
	ar rcs libnanovg_darwin_arm64.a nanovg_darwin_arm64.o
	rm nanovg_darwin_arm64.o $(NANOVG_IMPL_C)
	@echo "Successfully built libnanovg_darwin_arm64.a"


# 为 macOS amd64 架构编译 NanoVG 静态库
libnanovg_darwin_amd64.a: $(NANOVG_SRC_DIR)/nanovg.h $(NANOVG_SRC_DIR)/nanovg_gl.h
	@echo "Building NanoVG static library for darwin_amd64..."
	@echo '#define NANOVG_IMPLEMENTATION' > $(NANOVG_IMPL_C)
	@echo '#include "nanovg.h"' >> $(NANOVG_IMPL_C)
	@echo '#define $(NANOVG_GL_BACKEND_DEFINE)' >> $(NANOVG_IMPL_C)
	@echo '#include "nanovg_gl.h"' >> $(NANOVG_IMPL_C)

	clang $(NANOVG_CFLAGS_COMMON) -arch x86_64 -c $(NANOVG_IMPL_C) -o nanovg_darwin_amd64.o
	ar rcs libnanovg_darwin_amd64.a nanovg_darwin_amd64.o
	rm nanovg_darwin_amd64.o $(NANOVG_IMPL_C)
	@echo "Successfully built libnanovg_darwin_amd64.a"

# 创建 macOS 通用 (arm64 + amd64) NanoVG 静态库
libnanovg_darwin_universal.a: libnanovg_darwin_arm64.a libnanovg_darwin_amd64.a
	@echo "Creating universal NanoVG static library for macOS..."
	lipo -create libnanovg_darwin_arm64.a libnanovg_darwin_amd64.a -output libnanovg_darwin_universal.a
	@echo "Successfully built libnanovg_darwin_universal.a"

# 一个便捷的目标来构建通用的 macOS NanoVG 库
build_nanovg_macos: libnanovg_darwin_universal.a

# 清理 NanoVG 编译产物
clean_nanovg_macos:
	@echo "Cleaning NanoVG macOS build artifacts..."
	rm -f libnanovg_darwin_arm64.a libnanovg_darwin_amd64.a libnanovg_darwin_universal.a
	rm -f nanovg_darwin_arm64.o nanovg_darwin_amd64.o $(NANOVG_IMPL_C)
	@echo "NanoVG macOS build artifacts cleaned."

# 你可能还想为 Linux 和 Windows 添加类似的编译目标，使用对应的交叉编译器和标志。
# 例如，使用 zig cc 进行交叉编译：
# NANOVG_CFLAGS_LINUX = -fPIC -std=c99 -O2 -I$(NANOVG_SRC_DIR)
# libnanovg_linux_amd64.a:
#	...
#	zig cc -target x86_64-linux-gnu $(NANOVG_CFLAGS_LINUX) -c $(NANOVG_IMPL_C) -o nanovg_linux_amd64.o
#	ar rcs libnanovg_linux_amd64.a nanovg_linux_amd64.o
#	...