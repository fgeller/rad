#!/bin/bash

read -r -d '' PACK_DESC << EOM
© <a href="https://docs.python.org/2/copyright.html">Copyright</a> 1990-2015, Python Software Foundation.<br />
Documentation downloaded from <a href="https://docs.python.org/2/download.html">here</a>.<br />
Python is available for download <a href="https://www.python.org/download/releases/2.7/">here</a>.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

PACK_NAME=python27
DOWNLOAD_DIR=d/python-2.7.10-docs-html
../pad -indexer py27 -name $PACK_NAME -desc "$PACK_DESC" -version 2.7.10 -source $DOWNLOAD_DIR -dest /tmp/packs
