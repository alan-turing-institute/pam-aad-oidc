rm -rf build
make build
make -e DESTDIR=./build -e PREFIX=/usr install
mkdir build/DEBIAN
cp -r debian/control build/DEBIAN/control
dpkg-deb --build build
mv build.deb pam-aad-oidc.deb
lintian pam-aad-oidc.deb
