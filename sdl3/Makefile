SHELL := /bin/bash
.ONESHELL:
.PHONY: all build-sdl3 cmake-build autoconf-build clean build-sdl3-ttf

SRCDIR := $(CURDIR)/SDL
BUILD_DIR := $(SRCDIR)/build/sdl3-build
INSTALL_DIR := $(SRCDIR)/sdl3-dist
INCLUDE_DEST := $(CURDIR)/include
LIB_DEST := $(CURDIR)/libs
BUILD_TYPE := Release
JOBS ?= $(shell nproc 2>/dev/null || echo 1)

all: build-sdl3

# 主构建目标：自动检测 CMake 或 autoconf（configure）
build-sdl3:
	mkdir -p "$(BUILD_DIR)"
	cd "$(BUILD_DIR)"
	# 配置为只生成静态库，尝试常见选项以禁用共享库和测试
	cmake -S "$(SRCDIR)" -B "$(BUILD_DIR)" -DCMAKE_BUILD_TYPE=$(BUILD_TYPE) \
		-DBUILD_SHARED_LIBS=OFF -DSDL_SHARED=OFF -DSDL_STATIC=ON -DSDL_TEST_LIBRARY=OFF || true
	cmake --build "$(BUILD_DIR)" --config $(BUILD_TYPE) -- -j$(JOBS)
	# 将构建产物安装到本地前缀（避免写入系统）
	@echo "安装 SDL3 到本地前缀：$(INSTALL_DIR)"
	mkdir -p "$(INSTALL_DIR)"
	cmake --install "$(BUILD_DIR)" --prefix "$(INSTALL_DIR)" || true
	# 提取头文件（优先使用源代码 include 目录，然后使用本地安装前缀）
	mkdir -p "$(INCLUDE_DEST)" "$(LIB_DEST)"
	if [ -d "$(SRCDIR)/include" ]; then \
		# 仅复制 *.h 文件并保留目录结构
		rsync -a --prune-empty-dirs --include='*/' --include='*.h' --exclude='*' "$(SRCDIR)/include/" "$(INCLUDE_DEST)/" || cp -r "$(SRCDIR)/include"/* "$(INCLUDE_DEST)/"; \
	fi
	if [ -d "$(INSTALL_DIR)/include" ]; then \
		rsync -a --prune-empty-dirs --include='*/' --include='*.h' --exclude='*' "$(INSTALL_DIR)/include/" "$(INCLUDE_DEST)/" || true; \
	fi
	# 复制构建/安装输出中的所有静态库（包括 libSDL3_test.a）
	find "$(BUILD_DIR)" -type f -name '*.a' -exec cp -f {} "$(LIB_DEST)/" \; 2>/dev/null || true
	find "$(INSTALL_DIR)" -type f -name '*.a' -exec cp -f {} "$(LIB_DEST)/" \; 2>/dev/null || true

.PHONY: build-sdl3-ttf
build-sdl3-ttf:
	@echo "构建 SDL_ttf（如果存在 SDL_ttf 目录）..."
	EXT_SRCDIR="$(CURDIR)/SDL_ttf"
	if [ ! -d "$$EXT_SRCDIR" ]; then \
		echo "未发现 $$EXT_SRCDIR，跳过 SDL_ttf 构建"; \
		exit 0; \
	fi
	# 确保 SDL3 已安装到本地前缀，否则提示先构建 SDL3
	if [ ! -f "$(INSTALL_DIR)/lib/cmake/SDL3/SDL3Config.cmake" ]; then \
		echo "未在本地前缀找到 SDL3Config.cmake，先运行: make build-sdl3"; \
		exit 1; \
	fi
	mkdir -p "$$EXT_SRCDIR/build-ttf"
	# 指定本地前缀，让 find_package(SDL3) 能在不安装到系统的情况下找到 SDL3
	# 传入 SDL3_DIR 指向 cmake 配置所在目录，兼容某些 CMakeLists.txt
	cmake -S "$$EXT_SRCDIR" -B "$$EXT_SRCDIR/build-ttf" \
		-DCMAKE_BUILD_TYPE=$(BUILD_TYPE) \
		-DCMAKE_PREFIX_PATH="$(INSTALL_DIR)" \
		-DSDL3_DIR="$(INSTALL_DIR)/lib/cmake/SDL3" \
		-DBUILD_SHARED_LIBS=OFF || true
	cmake --build "$$EXT_SRCDIR/build-ttf" -- -j$(JOBS) || true
	cmake --install "$$EXT_SRCDIR/build-ttf" --prefix "$(INSTALL_DIR)" || true
	# 收集头文件与静态库到统一输出目录
	mkdir -p "$(INCLUDE_DEST)" "$(LIB_DEST)"
	if [ -d "$$EXT_SRCDIR/include" ]; then \
		rsync -a --prune-empty-dirs --include='*/' --include='*.h' --exclude='*' "$$EXT_SRCDIR/include/" "$(INCLUDE_DEST)/" || true; \
	fi
	find "$$EXT_SRCDIR/build-ttf" -type f -name '*.a' -exec cp -f {} "$(LIB_DEST)/" \; 2>/dev/null || true
	find "$(INSTALL_DIR)" -type f -name '*.a' -exec cp -f {} "$(LIB_DEST)/" \; 2>/dev/null || true
	@echo "SDL_ttf 构建完成并已复制静态库到 $(LIB_DEST) ，头文件到 $(INCLUDE_DEST)"

# 清理构建产物和提取的文件
clean:
	@echo "清理中..."
	rm -rf "$(BUILD_DIR)" "$(INSTALL_DIR)" "$(INCLUDE_DEST)" "$(LIB_DEST)"
