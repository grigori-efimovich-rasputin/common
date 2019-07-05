#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'


genny -in math_helper.generic -out math_helper.go gen "T=int,int32,int64,uint,uint32,uint64,float32,float64"

