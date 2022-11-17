#!/system/bin/sh

. ${0%/*}/config.ini

function mk_thermal_folder(){
    resetprop -p sys.thermal.data.path /data/vendor/thermal/
    resetprop -p vendor.sys.thermal.data.path /data/vendor/thermal/
    chattr -R -i -a /data/vendor/thermal
    if [ ! -d /data/vendor/thermal ];then
        rm -rf /data/vendor/thermal
        mkdir -p /data/vendor/thermal/config
    fi
    chmod 0771 /data/vendor/thermal
    chmod 644 /data/vendor/thermal/config/*
    chown root:system /data/vendor/thermal
    chown root:system /data/vendor/thermal/config/*
    chcon u:object_r:vendor_data_file:s0 /data/vendor/thermal
    chcon u:object_r:vendor_data_file:s0 /data/vendor/thermal/config/*
}

#检查温控二进制文件！
function check_thermal_control_file(){
    find /system /system_ext /vendor /product -iname 'mi_thermald' -type f -o -iname 'thermal-engine' -type f -o -iname 'thermalserviced' -type f 2>/dev/null | while read file ;do
    size="$(du -k $file | awk '{print $1}' | tr -cd '[0-9]'  )"
    details="$(cat $file 2>/dev/null | sed 's/[[:space:]]//g;s|/n||g' )"
    if test -f "$file" -a "$size" -ge "1" -a "$details" == "" ;then
        echo "温控二进制命令缺失！"
        exit 1
    fi
    done
}

#避免冻结电量和性能
function enable_miui_powerkeeper(){
    if test "$( pm list package | grep -w 'com.miui.powerkeeper' | wc -l)" -gt "0" ;then
        pm enable com.miui.powerkeeper >/dev/null 2>&1
        pm unsuspend com.miui.powerkeeper >/dev/null 2>&1
        pm unhide com.miui.powerkeeper >/dev/null 2>&1
        pm install-existing --user 0 com.miui.powerkeeper >/dev/null 2>&1
    fi
}

#重新启用电量与性能
function call_cloud_conf_release(){
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.cloudcontrol.CloudUpdateReceiver >/dev/null 2>&1
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.cloudcontrol.CloudUpdateJobService >/dev/null 2>&1
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.ui.CloudInfoActivity >/dev/null 2>&1
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.statemachine.PowerStateMachineService >/dev/null 2>&1
    am broadcast --user 0 -a update_profile com.miui.powerkeeper/com.miui.powerkeeper.cloudcontrol.CloudUpdateReceiver >/dev/null 2>&1
}

check_thermal_control_file
enable_miui_powerkeeper
call_cloud_conf_release
mk_thermal_folder

nohup charge-current $speed $temperaturewall $timesleep $rmthermal $file &