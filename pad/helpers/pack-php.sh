#!/bin/bash

PACK_NAME=php
SOURCE=http://nz2.php.net/manual/en/
DOWNLOAD_DIR=d/$PACK_NAME

# rm -vrf $DOWNLOAD_DIR
# mkdir -p $DOWNLOAD_DIR

# --random-wait
# --wait: don't spam the host
# -E: convert extensions, adds .html/.css
# -k: convert links to relative ones
# -H: allow spanning hosts
# -D: include the following list of domains
# -m: mirror the given site
# -l: limit the recursion level
# -P: local directory prefix for downloaded files
# -I "/manual/en","/images","/fonts","/cached.php" \
# --no-parent \
# wget --random-wait \
# 			-E \
# 			-e robots=off \
# 			-k \
# 			-H \
# 			-D nz2.php.net \
# 			-m \
# 			-p \
# 			--reject '*vote=up*','*vote=down*','*add-note.php*' \
# 			--accept-regex '/manual/en/.*|/cached.php*|/images*|/fonts*' \
# 			--restrict-file-names=windows \
# 			-P$DOWNLOAD_DIR \
# 			$SOURCE

find $DOWNLOAD_DIR/*/manual/* -maxdepth 0 -not -name "en" -exec rm -rf {} \;
echo ""                                    >> `find $DOWNLOAD_DIR/ -name "*theme-medium.css"`
echo "#head-nav { display: none }"         >> `find $DOWNLOAD_DIR/ -name "*theme-medium.css"`
echo "#breadcrumbs { display: none }"      >> `find $DOWNLOAD_DIR/ -name "*theme-medium.css"`
echo "aside.layout-menu { display: none }" >> `find $DOWNLOAD_DIR/ -name "*theme-medium.css"`
echo "div.page-tools { display: none }"    >> `find $DOWNLOAD_DIR/ -name "*theme-medium.css"`

read -r -d '' PACK_DESC << EOM
Copyright © 2001-2015 The PHP Group.<br />
More copyright information available <a href="https://secure.php.net/copyright.php">here</a>.<br />
PHP is available for download <a href="https://secure.php.net/downloads.php">here</a>.<br />
This is the a mirror of the online documentation with user contributions and is rather large when extracted (~500MB).<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer php -name $PACK_NAME -desc "$PACK_DESC" -version `date --iso` -source $DOWNLOAD_DIR -dest /tmp/packs
