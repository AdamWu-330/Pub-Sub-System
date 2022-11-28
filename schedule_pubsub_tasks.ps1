$trigger = New-ScheduledTaskTrigger -Once -At 20:38 -RepetitionInterval (New-TimeSpan -Minutes 15)
$action = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\pub\publish_bike_stations.exe"
Register-ScheduledTask -Action $action -Trigger $trigger -TaskName "pub-bike-stations"

$trigger_bike_status = New-ScheduledTaskTrigger -Once -At 20:41 -RepetitionInterval (New-TimeSpan -Minutes 15)
$action_bike_status = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\pub\publish_bike_status.exe"
Register-ScheduledTask -Action $action_bike_status -Trigger $trigger_bike_status -TaskName "pub-bike-status"

$trigger_camera = New-ScheduledTaskTrigger -Once -At 20:45 -RepetitionInterval (New-TimeSpan -Minutes 15)
$action_camera = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\pub\publish_camera.exe"
Register-ScheduledTask -Action $action_camera -Trigger $trigger_camera -TaskName "pub-camera"