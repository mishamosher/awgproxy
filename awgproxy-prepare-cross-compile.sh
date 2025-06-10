#!/bin/sh

go mod tidy -v
mkdir android_ndk
wget --no-verbose https://dl.google.com/android/repository/android-ndk-r27c-linux.zip -O android_ndk/ndk.zip
unzip android_ndk/ndk.zip -d android_ndk
mv "android_ndk/$(ls -1 -d android_ndk/*/ | cut -f 2 -d '/' | head -1)" android_ndk/extracted