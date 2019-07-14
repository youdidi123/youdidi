#!/bin/bash
PROCESSNAME=youdidi

#ESTAB=`netstat -an |grep -i estab |grep :8764 |grep -v "127.0.0.1"|wc -l`

PID=`ps -ef|grep ${PROCESSNAME}|grep -v 'grep'|awk '{print $2}'`

killpid()
{
	if [ "$PID" ]; then
		echo "${PID} will be killed!"
		i=0
		while [ $i -lt 10 ]
		do
		    PID=`ps -ef|grep ${PROCESSNAME}|grep -v 'grep'|awk '{print $2}'`
		    if [ "$PID" ];then
		        kill -9 $PID
		        sleep 3s
			((i++))
		    else
		        echo "$PID Killed succefully!"
		        break
		    fi
		done
		fi

	if [ "$PID" ]; then
		kill -9 $PID
		sleep 3
		if [ "$PID" ]; then
			echo $PID 进程无法结束.请检查服务器状态
			exit 1
		fi
	fi
}


killpid

docker start mymysql

./youdidi -d >logs/run.log &

ps -ef|grep ${PROCESSNAME}
