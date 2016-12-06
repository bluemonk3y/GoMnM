# #!/bin/sh
echo "Running...."
cat - | while read x ; do echo `date` $x ; done