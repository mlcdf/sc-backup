#! /usr/bin/env bash
set -e

PLATFORM=$(python3 -c 'import platform


def arch() -> str:
    if platform.machine() == "x86_64":
        return "amd64"

    if "arm" in platform.machine():
        if "64" in platform.architecture()[0]:
            return "arm64"
        else:
            return "arm"


def os() -> str:
    return platform.system().lower()


print(f"{os()}-{arch()}")')

# Download latest release of sc-backup
LOCATION=$(curl -s https://api.github.com/repos/mlcdf/sc-backup/releases/latest | grep browser_download_url | grep $PLATFORM | cut -d '"' -f 4)
echo $LOCATION
curl -L "${LOCATION}" -o sc-backup

chmod +x sc-backup

./sc-backup --collection mlcdf -o example-output
./sc-backup --list https://www.senscritique.com/liste/Vu_au_cinema/363578 -o example-output
