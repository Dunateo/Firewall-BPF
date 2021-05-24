#!/bin/bash
cd /home/
rm app.txt prog.txt
cd /snap/
ls -d */ >> /home/prog.txt
cd /opt/
ls -d */ >> /home/prog.txt
echo "firefox" >> /home/prog.txt
sed -e 's/\///' /home/prog.txt -i
sed '/bin/d' /home/prog.txt -i
sed '/core/d' /home/prog.txt -i
for line in $(cat /home/prog.txt)
do 
	var=`ps -edf | grep $line | wc -l`
	if [ $var != 1 ] 
	then 
		echo "$line" >> /home/app.txt
	fi 
done
