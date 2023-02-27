#!/usr/bin/python3
import os
import sys

os_arch = {
    "darwin": [
        "amd64"
    ],
    "linux": [
        "amd64"
    ],
    "windows": [
        "amd64"
    ]
}
pkg_name = "bsd"


def build_all(name: str):
    for os_, arch_s in os_arch.items():
        for arch in arch_s:
            if os_ == "windows":
                cmd = f"CGO_ENABLED=0 GOOS={os_} GOARCH={arch} go build -o build/{name}_{os_}_{arch}.exe"
                print(cmd)
                os.system(cmd)
            else:
                cmd = f"CGO_ENABLED=0 GOOS={os_} GOARCH={arch} go build -o build/{name}_{os_}_{arch}"
                print(cmd)
                os.system(cmd)


if __name__ == '__main__':
    build_all(pkg_name)