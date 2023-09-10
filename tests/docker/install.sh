#!/bin/bash
apt update
read -p "apt install flags: " PROMPT
apt install $PROMPT neovim
