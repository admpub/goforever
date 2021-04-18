export DISTPATH=../../dist/
export PKGPATH=./
mkdir $DISTPATH

linux_amd64() {
    export GOOS=linux
    export GOARCH=amd64
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH} -trimpath -ldflags="-s -w" $PKGPATH
}

linux_386() {
    export GOOS=linux
    export GOARCH=386
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH} -trimpath -ldflags="-s -w" $PKGPATH
}

linux_arm5() {
    export GOOS=linux
    export GOARM=5
    export GOARCH=arm
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH}${GOARM} -trimpath -ldflags="-s -w" $PKGPATH
}

linux_arm6() {
    export GOOS=linux
    export GOARM=6
    export GOARCH=arm
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH}${GOARM} -trimpath -ldflags="-s -w" $PKGPATH
}

linux_arm7() {
    export GOOS=linux
    export GOARM=7
    export GOARCH=arm
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH}${GOARM} -trimpath -ldflags="-s -w" $PKGPATH
}

linux_arm64() {
    export GOOS=linux
    export GOARM=
    export GOARCH=arm64
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH}${GOARM} -trimpath -ldflags="-s -w" $PKGPATH
}

darwin_amd64() {
    export GOOS=darwin
    export GOARCH=amd64
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH} -trimpath -ldflags="-s -w" $PKGPATH
}

windows_amd64() {
    export GOOS=windows
    export GOARCH=amd64
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH}.exe -trimpath -ldflags="-s -w" $PKGPATH
}

windows_386() {
    export GOOS=windows
    export GOARCH=386
    go build -o ${DISTPATH}forever_${GOOS}_${GOARCH}.exe -trimpath -ldflags="-s -w" $PKGPATH
}

case "$1" in
    "linux_amd64")
        linux_amd64
        ;;
    "linux_386")
        linux_386
        ;;
    "linux_arm5")
        linux_arm5
        ;;
    "linux_arm6")
        linux_arm6
        ;;
    "linux_arm7")
        linux_arm7
        ;;
    "linux_arm8"|"linux_arm64")
        linux_arm64
        ;;
    "darwin_amd64")
        darwin_amd64
        ;;
    "windows_amd64")
        windows_amd64
        ;;
    "windows_386")
        windows_386
        ;;
    *)
        linux_amd64
        linux_386
        linux_arm5
        linux_arm6
        linux_arm7
        linux_arm64
        darwin_amd64
        windows_amd64
        windows_386
        ;;
esac