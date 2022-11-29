$trigger_bike_stations = New-ScheduledTaskTrigger -Once -At 20:38 -RepetitionInterval (New-TimeSpan -Minutes 15)
$action_bike_stations = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\pub\publish_bike_stations.exe"
Register-ScheduledTask -Action $action_bike_stations -Trigger $trigger_bike_stations -TaskName "pub-bike-stations"

$trigger_bike_status = New-ScheduledTaskTrigger -Once -At 20:41 -RepetitionInterval (New-TimeSpan -Minutes 15)
$action_bike_status = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\pub\publish_bike_status.exe"
Register-ScheduledTask -Action $action_bike_status -Trigger $trigger_bike_status -TaskName "pub-bike-status"

$trigger_camera = New-ScheduledTaskTrigger -Once -At 20:45 -RepetitionInterval (New-TimeSpan -Minutes 15)
$action_camera = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\pub\publish_camera.exe"
Register-ScheduledTask -Action $action_camera -Trigger $trigger_camera -TaskName "pub-camera"

$trigger_generic = New-ScheduledTaskTrigger -Once -At 20:45 -RepetitionInterval (New-TimeSpan -Hours 1)
$action_generic = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\pub\publish_generic.exe"
Register-ScheduledTask -Action $action_generic -Trigger $trigger_generic -TaskName "pub-generic"

$trigger_save_bike_station = New-ScheduledTaskTrigger -Once -At 20:45
$action_save_bike_station = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\save_to_db\save_to_db_bike_stations.exe"
Register-ScheduledTask -Action $action_save_bike_station -Trigger $trigger_save_bike_station -TaskName "save-bike-station"

$trigger_save_bike_status = New-ScheduledTaskTrigger -Once -At 20:45
$action_save_bike_status = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\save_to_db\save_to_db_bike_status.exe"
Register-ScheduledTask -Action $action_save_bike_status -Trigger $trigger_save_bike_status -TaskName "save-bike-status"

$trigger_save_camera = New-ScheduledTaskTrigger -Once -At 20:45
$action_save_camera = New-ScheduledTaskAction -Execute "C:\Users\Adam\Pub-Sub-System\save_to_db\save_to_db_camera.exe"
Register-ScheduledTask -Action $action_save_camera -Trigger $trigger_save_camera -TaskName "save-camera"