## topic 2 
Redirects the project's output to the specified log file
Collect information from log files to kafka
Set different policies for different interfaces, e.g. add alerts when too many requests are triggered by the same IP address.
Add panic and other related log keyword detection, panic add an alarm message
Add interface survival detection, periodically send requests to determine whether the interface is alive or not

## topic 3
Use redis instead of go's map to store data
The key of redis can be deleted automatically when it expires
go's map does not automatically shrink memory, even if the key is deleted
Save hourly total sales data to redis
Every minute, the expired minute data is subtracted from the key
Each sale of a car gives the data +1, so that you do not need to iterate through sixty keys every time you get the sales rate
Add different error levels. This allows the monitoring system to select different alarm strategies based on different error level keywords

