#!/bin/sh
#===============================================================
# A load extension demonstration
#===============================================================

# Startup Xvfb
Xvfb -ac :99 -screen 0 1280x1024x16 > /dev/null 2>&1 &

sleep 1

# Show process
ps

export DISPLAY=:99.0

./starter -debug -workspace out
