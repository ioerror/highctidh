#!/bin/bash -x

export SOURCE_DATE_EPOCH=`git log -1 --pretty=%ct`
mkdir -p dist/source
mkdir -p build/tmp
export TMP=build/tmp
export VENV=venv${PYVER}
if [ ! -d $VENV ];
then
    python${PYVER} -m venv $VENV
else
    echo "venv exists: $VENV" 
fi

source $VENV/bin/activate && CC=${CC} TMP=$TMP pip${PYVER} --cache-dir=build/ install --upgrade pip
source $VENV/bin/activate && CC=${CC} TMP=$TMP pip${PYVER} --cache-dir=build/ install pytest pytest-xdist
source $VENV/bin/activate && CC=${CC} TMP=$TMP pip${PYVER} --cache-dir=build/ install --no-clean .
