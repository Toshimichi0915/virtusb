# Maintainer: Toshimichi0915 <main@toshimichi.net>
pkgname=virtusb
pkgver=1.0.1
pkgrel=1
pkgdesc='GUI to attach and detach host USB devices to libvirt VMs'
arch=('x86_64' 'aarch64')
url='https://github.com/Toshimichi0915/virtusb'
license=('MIT')
# GLFW dlopens the X11/Wayland libraries, so they don't show up as NEEDED.
depends=('libvirt' 'glibc' 'libglvnd' 'wayland' 'libxkbcommon'
         'libx11' 'libxcursor' 'libxext' 'libxi' 'libxinerama' 'libxrandr' 'libxrender')
optdepends=('libdecor: client-side window decorations on Wayland (e.g. GNOME)')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$pkgname-$pkgver"
  export CGO_CPPFLAGS="$CPPFLAGS"
  export CGO_CFLAGS="$CFLAGS"
  export CGO_CXXFLAGS="$CXXFLAGS"
  export CGO_LDFLAGS="$LDFLAGS"
  export GOFLAGS='-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw'
  go build -o virtusb .
}

check() {
  cd "$pkgname-$pkgver"
  go test .
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 virtusb "$pkgdir/usr/bin/virtusb"
  install -Dm644 net.toshimichi.virtusb.desktop -t "$pkgdir/usr/share/applications"
  install -Dm644 icon.svg "$pkgdir/usr/share/icons/hicolor/scalable/apps/net.toshimichi.virtusb.svg"
  install -Dm644 LICENSE -t "$pkgdir/usr/share/licenses/$pkgname"
}
