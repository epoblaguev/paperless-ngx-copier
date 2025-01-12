#!/bin/bash

cargo install cross

declare -a targets=("x86_64-pc-windows-gnu" "x86_64-unknown-linux-gnu" "aarch64-unknown-linux-gnu")

for target in "${targets[@]}"
do
    echo "$target"
    cross build --release --target "$target"
done