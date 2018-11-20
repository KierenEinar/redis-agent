pname=`whoami`
pids=`ps -fu $pname | grep redis-agent | grep -v grep | awk '{print $2}'`
for pid in $pids
do
    kill -9 $pid
done

