#!/bin/bash
find ./views -name "*.tmpl" -exec grep -h "class=" {} \; | \
  grep -o 'class="[^"]*"' | \
  sed 's/class="//g' | \
  sed 's/"//g' > tailwind-classes.txt
