#!/bin/bash

PACK_NAME=jquery
SOURCE=http://api.jquery.com/
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
     -D api.jquery.com,use.typekit.net,ajax.googleapis.com \
     -m \
     -l 1 \
     --restrict-file-names=windows \
     -p \
     -P$DOWNLOAD_DIR \
     $SOURCE

mkdir -p $DOWNLOAD_DIR/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts
pushd $DOWNLOAD_DIR/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts && wget http://api.jquery.com/jquery-wp-content/themes/jquery/css/fonts/fontawesome-webfont.eot && popd
pushd $DOWNLOAD_DIR/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts && wget http://api.jquery.com/jquery-wp-content/themes/jquery/css/fonts/fontawesome-webfont.ttf && popd
pushd $DOWNLOAD_DIR/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts && wget http://api.jquery.com/jquery-wp-content/themes/jquery/css/fonts/fontawesome-webfont.woff && popd
find $DOWNLOAD_DIR/api.jquery.com/jquery-wp-content/themes/jquery/css -name 'base.css*' -exec sed -iE 's$http://api.jquery.com$../../../../../api.jquery.com$g' {} \;
find $DOWNLOAD_DIR -name '*.html*' -exec sed -iE 's:<script>window.jQuery:<script><![CDATA[window.jQuery:g' {} \;
find $DOWNLOAD_DIR -name '*.html*' -exec sed -iE 's:)</script>:)]]></script>:g' {} \;
echo 'div#logo-events.constrain.clearfix { display: none; }' >> `find $DOWNLOAD_DIR -name '*style.css*'`
echo 'section#global-nav { display: none; }' >> `find $DOWNLOAD_DIR -name '*style.css*'`
echo 'nav#main.constrain.clearfix { display: none; }' >> `find $DOWNLOAD_DIR -name '*style.css*'`
echo 'footer.clearfix.simple { display: none; }' >> `find $DOWNLOAD_DIR -name '*style.css*'`

read -r -d '' PACK_DESC << EOM
Copyright 2015 <a href="https://jquery.org/team/">The jQuery Foundation</a>. <a href="https://jquery.org/license/">jQuery License</a><br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer jquery -name $PACK_NAME -desc "$PACK_DESC" -version 2.1.4 -source $DOWNLOAD_DIR -dest /tmp/packs
