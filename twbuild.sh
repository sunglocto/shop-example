#!/bin/bash
./tailwind_class_finder.sh && tailwindcss -o ./static/css/main.css --content tailwind-classes.txt
