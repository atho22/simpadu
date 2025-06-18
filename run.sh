
if ! pgrep -f 'binary_app'
then
    nohup /var/www/html/simpadu/main && echo "baru running" > ./out_test.txt
else
    echo "udah di-run" > ./out_test.txt
fi
