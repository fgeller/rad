#!/bin/bash
PFX=react
PKG=https://facebook.github.io/react/docs/getting-started.html
# --random-wait
# --wait: don't spam the host
# -E: convert extensions, adds .html/.css
# -k: convert links to relative ones
# -H: allow spanning hosts
# -D: include the following list of domains
# -m: mirror the given site
# -l: limit the recursion level
# -P: local directory prefix for downloaded files
wget --wait=0.5 \
     -E \
     -e robots=off \
     -k \
     -H \
     -D facebook.github.io \
     -l 1 \
     -m \
     -P$PFX \
     $PKG

find $PFX -name blog -type d -exec rm -vrf {} \;
find $PFX -name flux -type d -exec rm -vrf {} \;
find $PFX -name videos.html -type f -exec rm -vf {} \;
echo ".nav-main { visibility: hidden !important; }" >> `find $PFX -name "react-native.css"`
echo ".nav-main { visibility: hidden !important; }" >> `find $PFX -name "react.css"`
