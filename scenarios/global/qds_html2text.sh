#!/bin/bash
echo `pwd`
me=`basename -s .sh "$0"`
$me -a HTML2TEXT -r -l warn &