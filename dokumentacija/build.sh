#!/usr/bin/bash

pdflatex -interaction nonstopmode -shell-escape Aplikacija_za_preporuke.tex
biber Aplikacija_za_preporuke
pdflatex -interaction nonstopmode -shell-escape Aplikacija_za_preporuke.tex
pdflatex -interaction nonstopmode -shell-escape Aplikacija_za_preporuke.tex
