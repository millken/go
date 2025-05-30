#if defined(__APPLE__)
    #define GL_SILENCE_DEPRECATION // Suppress OpenGL deprecation warnings on macOS
    #include <OpenGL/gl3.h>      // Or <OpenGL/gl.h> if you compiled NanoVG for GL2
                                 // This should match the NANOVG_GL_HEADER in your Makefile
#else
    // Add appropriate OpenGL headers for other platforms if needed
    // e.g., #include <GL/gl.h> for Linux
#endif

#include <stdio.h>
#define NANOVG_IMPLEMENTATION     // <--- 包含 NanoVG 核心实现
#define NANOVG_GL3_IMPLEMENTATION // <--- 包含 GL3 后端实现
#include "nanovg.h"
#include "nanovg_gl.h"

int main(int argc, char const *argv[]) {
    NVGcontext* vg = NULL;
    int nvgFlags = NVG_NO_FONTSTASH | NVGL_DEBUG;
    vg = nvglCreate(nvgFlags);

    if (vg == NULL) {
        printf("nvglCreate() returned NULL. This is expected without an active OpenGL context.\n");
    } else {
        // This would be unexpected in this minimal test but indicates the function call itself worked.
        printf("nvglDelete() succeeded (unexpectedly, as no GL context was explicitly set up).\n");
        nvglDelete(vg); // Match the delete function to the create function
    }

    printf("\nNanoVG test program finished.\n");
    printf("If you see this message and there were no 'undefined symbol' errors during linking,\n");
    printf("the static library and its linkage with OpenGL are likely set up correctly.\n");

    return 0;
}