#!/bin/bash

rm -rf ./quickjs.h
rm -rf ./quickjs-libc.h
rm -rf ./libquickjs.a
rm -rf ./pkgs/quickjs-master

cd ./pkgs

unzip ./quickjs-master.zip

cd ./quickjs-master

rm -rf ./quickjs.h
rm -rf ./quickjs.c

cp ../quickjs.h ./quickjs.h
cp ../quickjs.c ./quickjs.c

cp ./quickjs.h ../../quickjs.h
cp ./quickjs-libc.h ../../quickjs-libc.h

make
cp ./libquickjs.a ../../libquickjs.a

cd ../

rm -rf ./quickjs-master
rm -rf ./__MACOSX

cd ../