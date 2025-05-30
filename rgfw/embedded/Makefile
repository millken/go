libRGFW_arm64.dylib:
	cp ../RGFW/RGFW.h RGFW.c
	clang  -c RGFW.c  -arch arm64 -D RGFW_IMPLEMENTATION -fPIC -D RGFW_EXPORT
	rm RGFW.c
	clang  -shared -arch arm64 RGFW.o -framework CoreVideo -framework Cocoa -framework OpenGL -framework IOKit -std=c99 -o libRGFW_darwin_arm64.dylib
	rm RGFW.o

libRGFW_arm64.a:
	cp ../RGFW/RGFW.h RGFW.c
	clang -c RGFW.c -arch arm64 -D RGFW_IMPLEMENTATION -fPIC -D RGFW_EXPORT -std=c99 -o RGFW_darwin_arm64.o
	rm RGFW.c
	ar rcs libRGFW_darwin_arm64.a RGFW_darwin_arm64.o
	rm RGFW_darwin_arm64.o

libRGFW_amd64.dylib:
	cp ../RGFW/RGFW.h RGFW.c
	clang  -c RGFW.c  -arch x86_64 -D RGFW_IMPLEMENTATION -fPIC -D RGFW_EXPORT
	rm RGFW.c
	clang  -shared -arch x86_64 RGFW.o -framework CoreVideo -framework Cocoa -framework OpenGL -framework IOKit -std=c99 -o libRGFW_darwin_amd64.dylib
	rm RGFW.o

libRGFW_windows_amd64.dll:
	cp ../RGFW/RGFW.h RGFW.c
	docker exec rgfw_linux_arm64_build \
		x86_64-w64-mingw32-gcc -shared RGFW.c -o libRGFW_windows_amd64.dll -D RGFW_IMPLEMENTATION -D RGFW_EXPORT -std=c99 -fPIC -lgdi32 -lopengl32 
	rm RGFW.c

windows_arm64.dll:
	unlink libRGFW_windows_arm64.dll
	cp ../RGFW/RGFW.h RGFW.c
	docker exec rgfw_linux_arm64_build \
		aarch64-w64-mingw32-gcc -shared RGFW.c -o libRGFW_windows_arm64.dll -D RGFW_IMPLEMENTATION -D RGFW_EXPORT -std=c99 -fPIC -lgdi32 -lopengl32 
	rm RGFW.c

linux_amd64.so:
	cp ../RGFW/RGFW.h RGFW.c
	rm -f libRGFW_linux_amd64.so
	cc -shared RGFW.c -o libRGFW_linux_amd64.so -Wall -Werror -Wextra -Wpedantic -Wconversion -Wsign-conversion -Wshadow -Wpointer-arith -Wvla -Wcast-align -Wstrict-overflow -Wstrict-aliasing -Wredundant-decls -Winit-self -Wmissing-noreturn -D __linux__ -D RGFW_DEBUG -D RGFW_IMPLEMENTATION -D RGFW_EXPORT -std=c99 -fPIC -ldl -lpthread -lX11 -lXrandr -lGL
	rm RGFW.c

libRGFW_linux_arm64.so:
	rm -f libRGFW_linux_arm64.so
	cp ../RGFW/RGFW.h RGFW.c
	docker exec rgfw_linux_arm64_build \
		gcc -shared RGFW.c -o libRGFW_linux_arm64.so -D RGFW_IMPLEMENTATION -D RGFW_EXPORT -std=c99 -fPIC -ldl -lpthread -lX11 -lXrandr -lGL
	rm RGFW.c
