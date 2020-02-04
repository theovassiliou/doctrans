#!/bin/bash
echo `pwd`
me=`basename -s .sh "$0"`
$me -a COUNT -r -l warn &