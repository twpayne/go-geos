#!/bin/bash

comm -23 <(rg -INo '\bGEOS[A-Z_a-z]+_r\b' /usr/include/geos_c.h | sort | uniq) <(rg -INo '\bGEOS[A-Z_a-z]+_r\b' *.c *.go | sort | uniq)