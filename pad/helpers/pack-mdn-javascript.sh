#!/usr/bin/env bash

PACK_NAME=mdn-javascript
SOURCE=https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference
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
# -np: no parent
wget --random-wait \
     -E \
     -e robots=off \
     -k \
     -H \
     --no-check-certificate \
     -D developer.mozilla.org,developer.cdn.mozilla.net,gravatar.com,secure.gravatar.com,i2.wp.com \
     -m \
     -l 2 \
     -np \
     -P$DOWNLOAD_DIR \
     $SOURCE

find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E "s:<span class='hidden'>contributors</span>: :g" {} \;
find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E 's:<script type="text/javascript">:<script type="text/javascript"><![CDATA[:g' {} \;
find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E 's:<script>:<script><![CDATA[:g' {} \;
find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E 's:</script>:]]></script>:g' {} \;

echo '' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo '' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo '#main-header { display: none !important; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'div.article-meta  { display: none !important; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'div.contributor-avatars  { display: none !important; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'div.global-notice  { display: none !important; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'div.wiki-block.contributors  { display: none !important; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'div.column-half { float: none !important; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'footer { visibility: hidden; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'div#toc { display: none; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`
echo 'div#wiki-left { display: none; }' >> `find $DOWNLOAD_DIR -name "mdn.*.css"`

read -r -d '' PACK_DESC << EOM
© 2005-2015 Mozilla Developer Network and individual contributors.<br />
Content is available under <a href="https://developer.mozilla.org/en-US/docs/MDN/About#Copyrights_and_licenses">these licenses</a>.<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer mdn -name $PACK_NAME -desc "$PACK_DESC" -version `date --iso` -source $DOWNLOAD_DIR -dest /tmp/packs

find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E 's:<script type="text/javascript"><!\[CDATA\[:<script type="text/javascript">:g' {} \;
find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E 's:<script><\!\[CDATA\[:<script>:g' {} \;
find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E 's:\]\]></script>:</script>:g' {} \;
