for i in {0..250}
do
    /home/ubuntu/Pub-Sub-System/update_subscriber_queues test_topic$i
done
echo "successfully subscribed test users to the newly created test topics"

