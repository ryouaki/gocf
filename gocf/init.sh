#!/bin/bash

rm -rf ./libquickjs.a
rm -rf ./pkgs/quickjs-master

cd ./pkgs

unzip ./quickjs-master.zip

cd ./quickjs-master

cp ./quickjs.h ../../quickjs.h
cp ./quickjs-libc.h ../../quickjs-libc.h

make
cp ./libquickjs.a ../../libquickjs.a

cd ../

rm -rf ./quickjs-master

cd ../

go install


# go env -w CGO_ENABLED=1