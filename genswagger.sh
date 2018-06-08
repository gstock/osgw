#!/bin/bash

swagger generate spec -o ./swagger.json --scan-models && swagger serve -F=swagger swagger.json