#!/bin/bash

OS=$(uname)
if [[ "$OS" == "Darwin" ]]; then
  bin/gonews-mac
elif [[ "$OS" == "Linux" ]]; then
  bin/gonews
elif [[ "$OS" == "CYGWIN"* || "$OS" == "MINGW"* || "$OS" == "MSYS"* ]]; then
  bin/gonews.exe
else
  echo "Unknown OS: $OS"
  exit 1
fi
