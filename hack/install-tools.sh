#!/usr/bin/env bash

cat .tool-versions | while read line
do
  echo $line | awk '{print $1}' | xargs -I{} asdf plugin add {}
  if [ $? -ne 0 ]; then
    continue
  fi
done

asdf install