#!/bin/bash

rm -rf ./libquickjs.a
rm -rf ./pkgs/quickjs-master

cd ./pkgs

unzip ./quickjs-master.zip

cd ./quickjs-master

cp ./quickjs.h ../../quickjs.h
cp ./quickjs-libc.h ../../quickjs-libc.h

