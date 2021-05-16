exec 5<>stream.fifo
cat <&5 >/dev/null & cat_pid=$!
sleep 1
kill "$cat_pid"
