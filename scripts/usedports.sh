# Use netstat to list all ports that are currently bound

netstat -aln | egrep ^tcp | fgrep LISTEN |
awk '{print $4}' | egrep -o '[0-9]+$' |
sort -n | uniq
