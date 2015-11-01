#!/bin/bash
PACK_NAME=django
SOURCE=http://django.readthedocs.org/en/1.6.x/
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
     -D django.readthedocs.org \
     -p \
     --restrict-file-names=windows \
     -m \
     -l 1 \
     -P$DOWNLOAD_DIR \
     $SOURCE

echo "#hd, .injected { display: none; }" >> `find $DOWNLOAD_DIR -name "default.css*"`
echo "#ft { visibility: hidden; margin-top: 5em; }" >> `find $DOWNLOAD_DIR -name "default.css*"`

read -r -d '' PACK_DESC << EOM
readthedocs.org's © 2010 - 2015:
<a href="http://ericholscher.com/">Eric Holscher</a>,
<a href="http://charlesleifer.com/">Charles Leifer</a>, and
<a href="http://bobbygrace.info/">Bobby Grace</a> for the 2010 <a href="http://djangodash.com/">Django Dash</a>.<br />
Django's © 2005-2015:
<a href="https://www.djangoproject.com/foundation/"> Django Software Foundation</a> and individual contributors.<br />
Django is a <a href="https://www.djangoproject.com/trademarks/">registered trademark</a> of the Django Software Foundation.<br />
Python's <a href="https://docs.python.org/2/copyright.html">copyright</a> is owned by the Python Software Foundation.<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer django -name $PACK_NAME -desc "$PACK_DESC" -version 1.6 -source $DOWNLOAD_DIR -dest /tmp/packs
