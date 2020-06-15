#!/bin/bash

for process in excrun/runshowdonetraden excrun/runshowquotation excrun/runshowdepth excrun/runshowklinesave excrun/runshowkline excrun/excdone
do
    pid=$(ps x | grep $process| grep -v grep | grep -v PPID | awk '{printf $1}' )

    if [ ! -z "$pid"  ]
    then
      echo $pid
      kill  $pid
    fi

    echo "->$process start"
    nohup php cli.php $process > runshow.txt &
done

exit


