#!/bin/bash
PFX=jquery
PKG=http://api.jquery.com/

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
     -P$PFX \
     $PKG

mkdir -p $PFX/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts
pushd $PFX/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts && wget http://api.jquery.com/jquery-wp-content/themes/jquery/css/fonts/fontawesome-webfont.eot && popd
pushd $PFX/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts && wget http://api.jquery.com/jquery-wp-content/themes/jquery/css/fonts/fontawesome-webfont.ttf && popd
pushd $PFX/api.jquery.com/jquery-wp-content/themes/jquery/css/fonts && wget http://api.jquery.com/jquery-wp-content/themes/jquery/css/fonts/fontawesome-webfont.woff && popd

find $PFX/api.jquery.com/jquery-wp-content/themes/jquery/css -name 'base.css*' -exec sed -iE 's$http://api.jquery.com$../../../../../api.jquery.com$g' {} \;

find $PFX -name '*.html*' -exec sed -iE 's:<script>window.jQuery:<script><![CDATA[window.jQuery:g' {} \;
find $PFX -name '*.html*' -exec sed -iE 's:)</script>:)]]></script>:g' {} \;
