#!/bin/bash

PACK_NAME=man-pages
SOURCE=http://man7.org/linux/man-pages/dir_all_alphabetic.html
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
     --no-check-certificate \
     -D man7.org \
     -m \
     -l 1 \
     -np \
     --restrict-file-names=windows \
     -p \
     -P$DOWNLOAD_DIR \
     $SOURCE

echo ''                                              >> `find $DOWNLOAD_DIR/*/linux -name "style.css"`
echo '.nav-bar { display: none !important; }'        >> `find $DOWNLOAD_DIR/*/linux -name "style.css"`
echo '.man-search-box { display: none !important; }' >> `find $DOWNLOAD_DIR/*/linux -name "style.css"`
echo '.start-footer { display: none !important; }'   >> `find $DOWNLOAD_DIR/*/linux -name "style.css"`
echo '.footer { display: none !important; }'         >> `find $DOWNLOAD_DIR/*/linux -name "style.css"`

read -r -d '' PACK_DESC << EOM
The online man pages are are maintained and (c) <a href="http://man7.org/mtk/index.html">Michael Kerrisk</a>.<br />
Man pages can be downloaded <a href="https://www.kernel.org/doc/man-pages/download.html">here</a>.<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer man -name $PACK_NAME -desc "$PACK_DESC" -version `date --iso` -source $DOWNLOAD_DIR -dest /tmp/packs
