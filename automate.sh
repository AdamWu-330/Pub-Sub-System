date >> /home/ubuntu/cronlog.txt

cd /home/ubuntu/Pub-Sub-System && /usr/local/go/bin/go run ./pub/publish_bike_stations.go >> /home/ubuntu/cronlog.txt
cd /home/ubuntu/Pub-Sub-System && /usr/local/go/bin/go run ./pub/publish_bike_status.go >> /home/ubuntu/cronlog.txt
cd /home/ubuntu/Pub-Sub-System && /usr/local/go/bin/go run ./pub/publish_camera.go >> /home/ubuntu/cronlog.txt

cd /home/ubuntu/Pub-Sub-System && timeout 30 /usr/local/go/bin/go run ./save_to_db/save_to_db_bike_stations.go >> /home/ubuntu/cronlog.txt
cd /home/ubuntu/Pub-Sub-System && timeout 30 /usr/local/go/bin/go run ./save_to_db/save_to_db_bike_status.go >> /home/ubuntu/cronlog.txt
cd /home/ubuntu/Pub-Sub-System && timeout 60 /usr/local/go/bin/go run ./save_to_db/save_to_db_camera.go >> /home/ubuntu/cronlog.txt

date >> /home/ubuntu/cronlog.txt