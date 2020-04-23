#!/bin/bash
protoc -I. --go_out=plugins=micro:. \
		*.proto
