#!/bin/bash

PACK_NAME=react
SOURCE=https://facebook.github.io/react/docs/getting-started.html
DOWNLOAD_DIR=d/$PACK_NAME

rm -vrf $DOWNLOAD_DIR
mkdir -p $DOWNLOAD_DIR

# --random-wait
# --wait: don't spam the host
# -E: convert extensions, adds .html/.css
# -k: convert links to relative ones
# -H: allow spanning hosts
# -D: include the following list of domains
# -m: mirror the given site
# -l: limit the recursion level
# -P: local directory prefix for downloaded files
wget --random-wait \
     -E \
     -e robots=off \
     -k \
     -H \
     -D facebook.github.io \
     -l 1 \
     -m \
     --restrict-file-names=windows \
     -p \
     -P$DOWNLOAD_DIR \
     $SOURCE

find $DOWNLOAD_DIR -name blog -type d -exec rm -vrf {} \;
find $DOWNLOAD_DIR -name flux -type d -exec rm -vrf {} \;
find $DOWNLOAD_DIR -name videos.html -type f -exec rm -vf {} \;
echo ".nav-main { visibility: hidden !important; }" >> `find $DOWNLOAD_DIR -name "react-native.css"`
echo ".nav-main { visibility: hidden !important; }" >> `find $DOWNLOAD_DIR -name "react.css"`

read -r -d '' PACK_DESC << EOM
© 2013–2015 Facebook Inc.<br />
Documentation licensed under <a href="https://creativecommons.org/licenses/by/4.0/">CC BY 4.0</a>.<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer react -name $PACK_NAME -desc "$PACK_DESC" -version 0.14.0 -source $DOWNLOAD_DIR -dest /tmp/packs
