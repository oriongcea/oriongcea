#!/bin/bash
step=8 #间隔的秒数，不能大于60

cd /data/server/api/cli

#while true
#do
    for process in excrun/runshowdonetraden excrun/runshowquotation excrun/runshowdepth excrun/runshowklinesave excrun/runshowkline excrun/excdone transrun/index useractive/calShareSl useractive/refund useractive/authSend useractive/convertedToTotal useractive/remove
    do
        pid=$(ps x | grep -w $process| grep -v grep | grep -v PPID | awk '{printf $1}' )

        if [ ! -z "$pid"  ]
        then
          echo "jc->$process"
          else
            echo "->$process start"
            nohup php cli.php $process > runshow.txt &
        fi
    done
#    sleep $step
#done
exit


