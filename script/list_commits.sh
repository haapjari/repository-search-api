#!/usr/bin/env bash
# Source: https://gist.github.com/textarcana/1306223

git log \
    --pretty=format:'{%n  "commit": "%H",%n  "author": "%aN",%n  "date": "%ad",%n},' \
    $@ | \
    perl -pe 'BEGIN{print "["}; END{print "]\n"}' | \
    perl -pe 's/},]/}]/'