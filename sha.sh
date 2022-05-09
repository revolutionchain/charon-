#!/bin/sh
if [ -n "$GIT_SHA" ]; then
	echo $GIT_SHA
else
	git rev-parse HEAD
fi
