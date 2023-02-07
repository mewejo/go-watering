#/bin/bash

curl -s https://api.github.com/repos/mewejo/go-watering/releases/latest \
    | grep "arm64\"" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qi - -O go-watering-latest

chmod +x go-watering-latest

