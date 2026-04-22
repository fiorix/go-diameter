#!/bin/bash
# Install open5gs/freeDiameter from source into the sdme rootfs.
# Pinned commit keeps the rootfs reproducible.

set -euo pipefail

FD_REPO="${FD_REPO:-https://github.com/open5gs/freeDiameter.git}"
FD_REF="${FD_REF:-main}"    # open5gs fork has no version tags
SRC_DIR=/tmp/freediameter
PREFIX=/usr/local
GO_VERSION="${GO_VERSION:-1.25.0}"

# Install a recent Go toolchain (Ubuntu 24.04 only ships 1.22; go.mod
# requires 1.25). Fetch the official tarball instead of relying on
# GOTOOLCHAIN auto-download, which needs DNS the test pod lacks.
ARCH="$(uname -m)"
case "${ARCH}" in
    aarch64) GO_ARCH=arm64 ;;
    x86_64)  GO_ARCH=amd64 ;;
    *) echo "unknown arch ${ARCH}" >&2; exit 1 ;;
esac
curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" \
    -o /tmp/go.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf /tmp/go.tar.gz
rm -f /tmp/go.tar.gz
ln -sf /usr/local/go/bin/go /usr/local/bin/go

rm -rf "${SRC_DIR}"
git clone --depth 1 --branch "${FD_REF}" "${FD_REPO}" "${SRC_DIR}" \
  || git clone "${FD_REPO}" "${SRC_DIR}"

cd "${SRC_DIR}"
# If --branch missed a plain commit, pin it explicitly.
if [ -n "${FD_COMMIT:-}" ]; then
    git checkout "${FD_COMMIT}"
fi

mkdir build
cd build
# Keep ALL_EXTENSIONS off so per-extension BUILD_<NAME> flags work.
# Disable extensions whose deps (json-schema-validator, MySQL client,
# PostgreSQL, SWIG) we don't want to install; they aren't used by any
# test scenario.
cmake -DCMAKE_INSTALL_PREFIX="${PREFIX}" \
      -DCMAKE_BUILD_TYPE=RelWithDebInfo \
      -DDISABLE_SCTP:BOOL=ON \
      -DBUILD_TESTING:BOOL=OFF \
      -DBUILD_DICT_JSON:BOOL=OFF \
      -DBUILD_APP_ACCT:BOOL=OFF \
      -DBUILD_APP_DIAMEAP:BOOL=OFF \
      -DBUILD_APP_RADGW:BOOL=OFF \
      -DBUILD_APP_SIP:BOOL=OFF \
      -DBUILD_APP_REDIRECT:BOOL=OFF \
      -DBUILD_DBG_INTERACTIVE:BOOL=OFF \
      -DBUILD_TEST_APP:BOOL=ON \
      -DBUILD_ACL_WL:BOOL=ON \
      -DBUILD_RT_DEFAULT:BOOL=ON \
      -DBUILD_DBG_MSG_DUMPS:BOOL=ON \
      -DBUILD_DICT_S6A:BOOL=ON \
      -DBUILD_DICT_DCCA_3GPP:BOOL=ON \
      ..
make -j"$(nproc)"
make install
ldconfig

# Clean up build tree to shrink the rootfs.
cd /
rm -rf "${SRC_DIR}"

# Create dirs the config templates expect.
install -d -m 0755 /etc/freeDiameter
install -d -m 0755 /etc/go-diameter
install -d -m 0755 /var/log/go-diameter-tests
