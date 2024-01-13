# Basic usage
```
/*
Command:
go build -o cli_todolist.exe main.go

Excecution:
cli_todolist.exe -list
cli_todolist.exe -add YOUR_TASK
cli_todolist.exe -complete=ID
cli_todolist.exe -del=ID
*/
```
# What's Next?

1. Add Cleanup flag: remove all completed tasks --> done 2024/1/11
```
cli_todolist.exe -cleanup
```
2. Add Urgent Level flag: color difference from common task --> done 2024/1/13
```
cli_todolist.exe -urgent -add YOUR_TASK
```
3. Notification (Key Feature): Urgent Level will be notified every N minutes --> done 2024/1/13


Build Monitor executable file:
```
go build -o todoListMonitor.exe monitor.go
```

Start background monitoring(Window Terminal):
```
start /b todoListMonitor.exe
```

Stop background monitoring:

1. Find PID
```
tasklist | findstr "todoListMonitor.exe"
```

2. Kill PID (add /F if can't not be killed)
```
taskkill /PID <PID>
taskkill /F /PID <PID>
```

# Refs
https://github.com/joefazee/go-toto-app/tree/main
https://github.com/alexeyco/simpletable