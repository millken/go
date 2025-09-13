github.com/ungerik/go-cairo
https://cairographics.org/releases/cairo-1.18.4.tar.xz
cd cairo-1.18.4

meson setup build \
  --prefix=/usr/local \
  --buildtype=release \
  --default-library=static \
  -Dpng=enabled \
  -Dfreetype=enabled \
  -Dfontconfig=enabled \
  -Dtests=disabled \
  -Dgtk_doc=false

meson compile -C build