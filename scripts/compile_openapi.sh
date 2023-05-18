#!/bin/bash

set -ex

# requires nodejs installed: https://nodejs.org
npx @redocly/cli build-docs openapi.yaml -o openapi.html