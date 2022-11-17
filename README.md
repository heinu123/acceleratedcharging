# acceleratedcharging
MIUI全机型充电加速(888系列可用)

充电时自动修改温控(实时生效)和强行修改充电电流实现满速充电(跑满120w)

温度墙45度

要设置温度墙或者充电电流请修改模块目录下的config.ini文件

# go编译命令:
```bash
go build -ldflags="-s -w" -o ./system/bin/charge-current ./system/bin/charge-current.go
```

# termux一键下载&编译&安装模块:
```bash
git clone https://github.com/heinu123/acceleratedcharging acceleratedcharging && cd acceleratedcharging && go build -ldflags="-s -w" -o ./system/bin/charge-current ./system/bin/charge-current.go && rm -rf ./system/bin/charge-current.go && zip -q -r acceleratedcharging.zip * && su -c magisk --install-module ./acceleratedcharging.zip
```
