libRGFW_arm64.dylib:
	cp ../RGFW/RGFW.h RGFW.c
	clang  -c RGFW.c  -arch arm64 -D RGFW_IMPLEMENTATION -fPIC -D RGFW_EXPORT
	rm RGFW.c
	clang  -shared -arch arm64 RGFW.o -framework CoreVideo -framework Cocoa -framework OpenGL -framework IOKit -std=c99 -o libRGFW_darwin_arm64.dylib
	rm RGFW.o

libRGFW_amd64.dylib:
	cp ../RGFW/RGFW.h RGFW.c
	clang  -c RGFW.c  -arch x86_64 -D RGFW_IMPLEMENTATION -fPIC -D RGFW_EXPORT
	rm RGFW.c
	clang  -shared -arch x86_64 RGFW.o -framework CoreVideo -framework Cocoa -framework OpenGL -framework IOKit -std=c99 -o libRGFW_darwin_amd64.dylib
	rm RGFW.o

libRGFW_windows_amd64.dll:
	cp ../RGFW/RGFW.h RGFW.c
	docker run --rm -v $(PWD):/work -e CROSS_TRIPLE=x86_64-w64-mingw32 -w /work \
        startergo/crossbuild:latest \
        x86_64-w64-mingw32-gcc -shared RGFW.c -o libRGFW_windows_amd64.dll -D RGFW_IMPLEMENTATION -D RGFW_EXPORT -std=c99 -fPIC -lgdi32 -lopengl32
	rm RGFW.c

libRGFW_linux_amd64.so:
	cp ../RGFW/RGFW.h RGFW.c
	docker exec rgfw_linux_arm64_build \
    	cc -shared RGFW.c -o libRGFW_linux_amd64.so -D RGFW_IMPLEMENTATION -D RGFW_EXPORT -std=c99 -fPIC -ldl -lpthread -lX11 -lXrandr -lGL
	rm RGFW.c

libRGFW_linux_arm64.so:
	rm -f libRGFW_linux_arm64.so
	cp ../RGFW/RGFW.h RGFW.c
	docker exec rgfw_linux_arm64_build \
		gcc -shared RGFW.c -o libRGFW_linux_arm64.so -D RGFW_IMPLEMENTATION -D RGFW_EXPORT -std=c99 -fPIC -ldl -lpthread -lX11 -lXrandr -lGL
	rm RGFW.c
