package meminfo

import (
	"fmt"
	"os/exec"
)

/**
adb shell dumpsys meminfo -d
Applications Memory Usage (in Kilobytes):
Uptime: 10495398 Realtime: 10495398

Total PSS by process:
     70,146K: system (pid 1636)
     59,674K: com.android.systemui (pid 1758)
     23,106K: zygote (pid 1306)
     22,201K: com.android.launcher3 (pid 2121 / activities)
     21,074K: com.android.phone (pid 1901)
     19,333K: com.dingdayu.smscat:xg_service_v3 (pid 2538)
     15,801K: com.dingdayu.helloandroid (pid 3384 / activities)
     12,462K: com.dingdayu.smscat (pid 3493)
     11,533K: android.process.acore (pid 2081)
     10,392K: com.android.inputmethod.latin (pid 1750)
      7,757K: com.android.deskclock (pid 1733)
      7,717K: com.android.packageinstaller (pid 3442)
      5,361K: com.android.gallery3d (pid 3310)
      4,434K: android.ext.services (pid 1973)
      4,350K: com.android.printspooler (pid 3328)
      4,193K: com.android.keychain (pid 3346)
      4,156K: com.android.defcontainer (pid 3288)
      3,595K: com.svox.pico (pid 3368)
      3,068K: logd (pid 1250)
      3,040K: audioserver (pid 1307)
      2,918K: mediaserver (pid 1315)
      2,647K: media.extractor (pid 1314)
      2,483K: media.codec (pid 1312)
      2,269K: cameraserver (pid 1308)
      1,953K: surfaceflinger (pid 1303)
      1,920K: mediadrmserver (pid 1313)
      1,576K: wpa_supplicant (pid 1732)
      1,442K: vold (pid 1259)
      1,375K: netd (pid 1316)
      1,323K: drmserver (pid 1309)
      1,240K: sdcard (pid 1856)
      1,236K: sdcard (pid 1795)
      1,079K: /init (pid 1)
      1,077K: hostapd (pid 1536)
      1,021K: adbd (pid 1331)
        916K: keystore (pid 1311)
        739K: rild (pid 1317)
        708K: installd (pid 1310)
        703K: fingerprintd (pid 1319)
        671K: gatekeeperd (pid 1320)
        640K: ueventd (pid 926)
        638K: lmkd (pid 1300)
        535K: logcat (pid 3302)
        532K: servicemanager (pid 1302)
        520K: healthd (pid 1297)
        501K: dumpsys (pid 3517)
        487K: dnsmasq (pid 1538)
        468K: sh (pid 1304)
        444K: perfprofd (pid 1324)
        444K: sh (pid 1484)
        361K: libxguardian.so (pid 2583)
        331K: ipv6proxy (pid 1517)
        278K: debuggerd (pid 1258)
        238K: debuggerd:signaller (pid 1261)

Total PSS by OOM adjustment:
     64,927K: Native
         23,106K: zygote (pid 1306)
          3,068K: logd (pid 1250)
          3,040K: audioserver (pid 1307)
          2,918K: mediaserver (pid 1315)
          2,647K: media.extractor (pid 1314)
          2,483K: media.codec (pid 1312)
          2,269K: cameraserver (pid 1308)
          1,953K: surfaceflinger (pid 1303)
          1,920K: mediadrmserver (pid 1313)
          1,576K: wpa_supplicant (pid 1732)
          1,442K: vold (pid 1259)
          1,375K: netd (pid 1316)
          1,323K: drmserver (pid 1309)
          1,240K: sdcard (pid 1856)
          1,236K: sdcard (pid 1795)
          1,079K: /init (pid 1)
          1,077K: hostapd (pid 1536)
          1,021K: adbd (pid 1331)
            916K: keystore (pid 1311)
            739K: rild (pid 1317)
            708K: installd (pid 1310)
            703K: fingerprintd (pid 1319)
            671K: gatekeeperd (pid 1320)
            640K: ueventd (pid 926)
            638K: lmkd (pid 1300)
            535K: logcat (pid 3302)
            532K: servicemanager (pid 1302)
            520K: healthd (pid 1297)
            501K: dumpsys (pid 3517)
            487K: dnsmasq (pid 1538)
            468K: sh (pid 1304)
            444K: perfprofd (pid 1324)
            444K: sh (pid 1484)
            361K: libxguardian.so (pid 2583)
            331K: ipv6proxy (pid 1517)
            278K: debuggerd (pid 1258)
            238K: debuggerd:signaller (pid 1261)
     70,146K: System
         70,146K: system (pid 1636)
     80,748K: Persistent
         59,674K: com.android.systemui (pid 1758)
         21,074K: com.android.phone (pid 1901)
     15,801K: Foreground
         15,801K: com.dingdayu.helloandroid (pid 3384 / activities)
      4,434K: Visible
          4,434K: android.ext.services (pid 1973)
     10,392K: Perceptible
         10,392K: com.android.inputmethod.latin (pid 1750)
     19,333K: A Services
         19,333K: com.dingdayu.smscat:xg_service_v3 (pid 2538)
     22,201K: Home
         22,201K: com.android.launcher3 (pid 2121 / activities)
     61,124K: Cached
         12,462K: com.dingdayu.smscat (pid 3493)
         11,533K: android.process.acore (pid 2081)
          7,757K: com.android.deskclock (pid 1733)
          7,717K: com.android.packageinstaller (pid 3442)
          5,361K: com.android.gallery3d (pid 3310)
          4,350K: com.android.printspooler (pid 3328)
          4,193K: com.android.keychain (pid 3346)
          4,156K: com.android.defcontainer (pid 3288)
          3,595K: com.svox.pico (pid 3368)

Total PSS by category:
     64,586K: Native
     46,719K: Dalvik
         30,283K: .Heap
         11,430K: .LOS
          4,926K: .LinearAlloc
          4,414K: .Zygote
          2,972K: .GC
          1,276K: .IndirectRef
            592K: .NonMoving
            176K: .JITCache
     44,440K: .so mmap
     42,272K: .dex mmap
     41,699K: .oat mmap
     27,691K: .apk mmap
     24,426K: .art mmap
     18,944K: Ashmem
     10,082K: Other mmap
      9,913K: Unknown
      9,350K: Dalvik Other
      6,480K: Stack
      1,867K: .ttf mmap
        477K: Other dev
        160K: .jar mmap
          0K: Cursor
          0K: Gfx dev
          0K: EGL mtrack
          0K: GL mtrack
          0K: Other mtrack

Total RAM: 1,550,680K (status normal)
 Free RAM: 1,197,288K (   61,124K cached pss +   216,008K cached kernel +   920,156K free)
 Used RAM:   337,922K (  287,982K used pss +    49,940K kernel)
 Lost RAM:    15,470K
   Tuning: 384 (large 384), oom   184,320K, restore limit    61,440K (high-end-gfx)
*/

func List() (str string, err error) {
	cmd := exec.Command("adb", "shell", "dumpsys", "meminfo")

	fmt.Println(cmd)
	return "", nil
}
