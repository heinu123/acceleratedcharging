#!/system/bin/sh
SKIPUNZIP=0
rm -rf /data/thermal/
set_perm_recursive $MODPATH 0 0 0777 0777
set_perm $MODPATH/system/bin/charge-current 0 0 0777 0777
