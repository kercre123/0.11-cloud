#!/bin/bash

if [[ $1 == "" ]]; then
    echo "You must provide an architecture to compile to."
    echo "Supported: amd64, arm64"
    exit 1
fi

if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root. sudo ./setup.sh"
    exit 1
fi

UNAME=$(uname -a)
ROOT="/root"

if [[ ${UNAME} == *"Darwin"* ]]; then
    if [[ -f /usr/local/Homebrew/bin/brew ]] || [[ -f /opt/Homebrew/bin/brew ]]; then
        TARGET="darwin"
        ROOT="$HOME"
        echo "macOS detected."
        if [[ ! -f /usr/local/go/bin/go ]]; then
            if [[ -f /usr/local/bin/go ]]; then
                mkdir -p /usr/local/go/bin
                ln -s /usr/local/bin/go /usr/local/go/bin/go
            else
                echo "Go was not found. You must download it from https://go.dev/dl/ for your macOS."
                exit 1
            fi
        fi
    else
        echo "macOS detected, but 'brew' was not found. Install it with the following command and try running setup.sh again:"
        echo '/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"'
        exit 1
    fi
    elif [[ -f /usr/bin/apt ]]; then
    TARGET="debian"
    echo "Debian-based Linux detected."
    elif [[ -f /usr/bin/pacman ]]; then
    TARGET="arch"
    echo "Arch Linux detected."
    elif [[ -f /usr/bin/dnf ]]; then
    TARGET="fedora"
    echo "Fedora/openSUSE detected."
else
    echo "This OS is not supported. This script currently supports Linux with either apt, pacman, or dnf."
    if [[ ! "$1" == *"--bypass-target-check"* ]]; then
        echo "If you would like to get the required packages yourself, you may bypass this by running setup.sh with the --bypass-target-check flag"
        echo "The following packages are required (debian apt in this case): wget openssl net-tools libsox-dev libopus-dev make iproute2 xz-utils libopusfile-dev pkg-config gcc curl g++ unzip avahi-daemon git"
        exit 1
    fi
fi

if [[ "${UNAME}" == *"x86_64"* ]]; then
    ARCH="x86_64"
    echo "amd64 architecture confirmed."
    elif [[ "${UNAME}" == *"aarch64"* ]] || [[ "${UNAME}" == *"arm64"* ]]; then
    ARCH="aarch64"
    echo "aarch64 architecture confirmed."
    elif [[ "${UNAME}" == *"armv7l"* ]]; then
    ARCH="armv7l"
    echo "armv7l (32-bit) WARN: The Coqui and VOSK bindings are broken for this platform at the moment, so please choose Picovoice when the script asks. wire-pod is designed for 64-bit systems."
    STT=""
else
    echo "Your CPU architecture not supported. This script currently supports x86_64, aarch64, and armv7l."
    exit 1
fi

if [[ ${TARGET} == "debian" ]]; then
    apt update -y
    apt install -y wget openssl net-tools libsox-dev libopus-dev make iproute2 xz-utils libopusfile-dev pkg-config gcc curl g++ unzip avahi-daemon git libasound2-dev libsodium-dev golang-go
elif [[ ${TARGET} == "arch" ]]; then
    pacman -Sy --noconfirm
    sudo pacman -S --noconfirm wget openssl net-tools sox opus make iproute2 opusfile curl unzip avahi git libsodium go pkg-config
elif [[ ${TARGET} == "fedora" ]]; then
    dnf update
    dnf install -y wget openssl net-tools sox opus make opusfile curl unzip avahi git libsodium-devel go
elif [[ ${TARGET} == "darwin" ]]; then
    sudo -u $SUDO_USER brew update
    sudo -u $SUDO_USER brew install wget pkg-config opus opusfile
fi

echo "getting vosk stuff..."

mkdir -p deploy
cd deploy

if [[ ! -d model ]]; then
    mkdir voskModel
    cd voskModel
    wget https://alphacephei.com/vosk/models/vosk-model-en-us-0.22-lgraph.zip
    unzip *.zip
    rm *.zip
    mv * ../model
    cd ../
    rm -r voskModel
fi

if [[ ! -d vosk ]]; then
    mkdir voskGet
    cd voskGet
    if [[ $1 == "arm64" ]]; then
        wget https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-linux-aarch64-0.3.45.zip
    elif [[ $1 == "amd64" ]]; then
        wget https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-linux-x86_64-0.3.45.zip
    fi
    unzip *.zip
    rm *.zip
    mv * ../vosk
    cd ../
    rm -r voskGet
fi

if [[ ! -d opus ]]; then
    mkdir opus
    cd ..
    BPREFIX="$(pwd)/deploy/opus"
    mkdir opusBuild
    cd opusBuild
    if [[ $1 == "arm64" ]]; then
        PODHOST=aarch64-linux-gnu
    else
        PODHOST=x86_64-linux-gnu
    fi
    git clone https://github.com/xiph/opus
    cd opus
    git checkout 0dc559f060db0d62d95f424e3fd26a5f673b2f6b
    ./autogen.sh
    autoreconf -i
    ./autogen.sh
    ./configure --host=${PODHOST} --prefix=$BPREFIX
    make -j8
    make install
    cd ../../
fi

echo "Creating server deployment tar..."

rm -f 011-cloud

CC=aarch64-linux-gnu-gcc GCC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=$1 CGO_LDFLAGS="-L$(pwd)/deploy/vosk -L$(pwd)/deploy/opus/lib" CGO_CFLAGS="-I$(pwd)/deploy/vosk -I$(pwd)/deploy/opus/include" go build -tags nolibopusfile -o deploy/011-cloud main.go

cp -r webroot deploy/
cp install.sh deploy/
cp 011-cloud.sh deploy/
cp stttest.pcm deploy/
cp -r intent-data deploy/
cp weather-map.json deploy/
mkdir deploy/secrets

tar -zcvf deploy.tar.gz -C deploy .

#rm -rf deploy