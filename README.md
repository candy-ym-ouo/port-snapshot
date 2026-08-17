# port-snapshot

采集当前监听端口、PID、进程名和启动用户，并支持端口范围过滤与 CSV 导出。

```sh
go run . -min 8000 -max 9000
go run . -csv snapshot.csv
```

macOS/BSD 使用 `lsof`，Linux 使用 `ss`，Windows 使用 `netstat`。部分进程信息需要管理员权限；扫描期间退出的进程会保留为未知字段并给出警告。
