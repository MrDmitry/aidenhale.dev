#!/bin/bash
read -p "Install location: " LOC
mkdir $LOC
dd if=/dev/urandom of=$LOC/payload  bs=1M  count=1152
read -p "Happy Y/N: " HAPPY
echo $HAPPY > $LOC/happy
