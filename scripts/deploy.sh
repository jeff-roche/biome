#!/bin/bash

# Fetch the latest tags
git fetch --tags

CURRENT=$(svu current)
NEXT=$(svu next)

if [ $CURRENT != $NEXT ]
then
  echo "Tagging with" $NEXT
  git tag $NEXT
  git tag latest
  git push --tags

  # Do the release
  VERSION=$NEXT goreleaser --rm-dist
else
  echo "No new version detected. Skipping release."
fi