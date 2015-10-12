#!/bin/bash
set -e

PFX=clojure-stdlib
#PKG=https://clojuredocs.org/core-library/vars
PKG=http://bert:5000/core-library/vars

# --random-wait
# --wait: don't spam the host
# -E: convert extensions, adds .html/.css
# -k: convert links to relative ones
# -H: allow spanning hosts
# -D: include the following list of domains
# -m: mirror the given site
# -l: limit the recursion level
# -P: local directory prefix for downloaded files
wget --wait=0 \
     -E \
     -e robots=off \
     -k \
     -H \
     -D bert,192.168.1.10,gravatar.com,avatars.githubusercontent.com \
     -m \
     -P$PFX \
     $PKG
