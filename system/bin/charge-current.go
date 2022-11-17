package main

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

var (
    temperature     string
    batterylevel     string
    temperaturewall    string
    speed       string
    rmthermal       string
    timesleep       string
    file    string
    times       int
    start     int
    stop     int
)

func shell(command string, su bool ) string { //调用shell执行命令(root权限)
    var output []byte
    var err error
    var cmd *exec.Cmd
    if su == true {
        cmd = exec.Command("su", "-c", command)
    } else {
        cmd = exec.Command("bash", "-c", command)
    }
    if output, err = cmd.CombinedOutput(); err == nil {
    }
    return string(output)
}

func uninstalltemperaturecontrol() { //删除温控
    if rmthermal == "true" {
        shell("echo \"#\" >/data/vendor/thermal/config/{thermal-phone.conf,thermal-4k.conf,thermal-app.conf,thermal-arvr.conf,thermal-camera.conf,thermal-charge.conf,thermal-class0.conf,thermal-hp-normal.conf,thermal-huanji.conf,thermal-mgame.conf,thermal-navigatstart.conf,thermal-nolimits.conf,thermal-normal.conf,thermal-per-camera.conf,thermal-per-class0.conf,thermal-per-huanji.conf,thermal-per-navigatstart.conf,thermal-per-normal.conf,thermal-per-phone.conf,thermal-per-video.conf,thermal-phone.conf,thermal-tgame.conf,thermal-video.conf,thermal-videochat.conf,thermal-yuanshen.conf,thermald-devices.conf,thermal-scene.conf}",true)
    }
    shell("echo '" + speed + "' > /data/adb/modules/acceleratedcharging/" + file,true) //写入充电电流到模块缓存文件
    shell("mount /data/adb/modules/acceleratedcharging/" + file + " /sys/class/power_supply/battery/" + file,true) //通过mount命令挂载充电电流速度
    shell("setprop ctl.stop mi_thermald",true)
    shell("setprop ctl.restart mi_thermald",true)
}

func installtemperaturecontrol() {
    shell("rm -rf /data/vendor/thermal/config/*.conf",true)
    shell("umount /sys/class/power_supply/battery/" + file,true)
    shell("setprop ctl.stop mi_thermald",true)
    shell("setprop ctl.restart mi_thermald",true)
}



func sleeps(times int) { //硬核休眠
    sum := 1
    for sum <= times {
        sum = sum + 1
        time.Sleep(time.Second)
    }
}

func runlog(text string) {
    shell("echo \"$(TZ=Asia/Shanghai date \"+%Y-%m-%d %H:%M:%S\")\n" + text + "\" >>/data/adb/modules/acceleratedcharging/charge-current.log",true);
    fmt.Println(text)
}
func main() {
    //读取命令行参数
    speed = os.Args[1]
    temperaturewall = os.Args[2]
    timesleep = os.Args[3]
    rmthermal = os.Args[4]
    file = os.Args[5]
    timesleep, err := strconv.Atoi(timesleep) //将string类型转为int类型
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    //初始化变量
    start = 0
    stop = 0
    shell("rm -rf /data/adb/modules/acceleratedcharging/charge-current.log && touch /data/adb/modules/acceleratedcharging/charge-current.log",true);
    for true { //循环
        var batterydata = shell("dumpsys battery",true)
        temperature = shell("cat /sys/class/power_supply/battery/temp",true);
        var dl = strings.Contains(batterydata, "status: 2")
        if dl { //判断是否在充电
            if temperature > temperaturewall {
                installtemperaturecontrol()
                start = 0
                stop = 1
                runlog("电池温度达到设置温度墙 已恢复快充")
            } else {
                if strings.Contains(batterydata, "level: 100") {
                    installtemperaturecontrol() //恢复
                    start = 1
                    stop = 0
                    runlog("已满电 已恢复快充")
                } else {
                    if start == 0 {
                        uninstalltemperaturecontrol() //删除温控 修改充电速度
                        start = 1
                        stop = 0
                        runlog("已修改快充")
                    }
                }
            }
        } else {
            if strings.Contains(batterydata, "level: 100") {
                runlog("已满电")
            } else {
                if stop == 0 {
                    installtemperaturecontrol() //恢复
                    start = 0
                    stop = 1
                    runlog("已恢复快充")
                }
            }
        }
        runlog("循环结束")
        sleeps(timesleep) //休眠
    }
}
