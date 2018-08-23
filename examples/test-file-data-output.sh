#!/bin/bash
INPUT_FILE=$1
set -x 

echo "Testing $INPUT_FILE"

# check the output contains the right strings
grep -q "sweaters Expected: 15 actual: 15" $INPUT_FILE || exit 1
grep -q "Json File .jsondata.sample expected: sample actual: sample" $INPUT_FILE || exit 1
grep -q "Text File .script expected scriptdata actual: scriptdata" $INPUT_FILE || exit 1

exit 0