# Pub-Sub-System
Backend server code for CVST pub-sub, fetching source data and send to users in real time

./client:
contains the user program for receiving the subscribe content. The exe file is provided on the Frontend UI for the user to download. To run it, "go run ./client/subscriber_receive.go <user email>"
 
./pub:
 scripts for publishing different topics, bikes, camera, yorkopendata are known topics, another generic script is available for publishing other user added new topics. For any of this script, it publishes the topic(s) to RabbitMQ a) message exchange for real-time publishing, and b) topic work queue for saving to MongoDB. To run any of these scripts, run "go run ./pub/<publish_script.go>"
  
./save_to_db:
 contains work scripts for saving topics to MongoDB. to run any of these, run "go run ./save_to_dv/<save_script.go>"

