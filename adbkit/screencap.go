package adbkit

import (
	"bytes"
	"errors"
	"io/ioutil"
	"regexp"
)

// Screencap 用于获取屏幕截图(未压缩)
func (c Client) Screencap(serial string) (buf []byte, err error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return
	}

	// 写入命令
	_, err = conn.Write(EncodeCommend("shell:screencap -p 2>/dev/null"))
	if err != err {
		return
	}

	stat := make([]byte, 4)
	_, err = conn.Read(stat)
	if err != nil {
		return nil, errors.New("adb response: Fail")
	}

	if string(stat) == OKAY {
		buf, _ = ioutil.ReadAll(conn)
		// 此处用于 将 \r\n 替换为 \n,
		// https://stackoverflow.com/questions/13578416/read-binary-stdout-data-from-adb-shell
		buf = bytes.Replace(buf, []byte("\x0D\x0A"), []byte("\x0A"), -1)
		return
	} else if string(stat) == FAIL {
		return nil, errors.New("adb response: Fail")
	}

	return nil, errors.New("adb response: " + string(stat))
}

// ScreenSize 用于获取屏幕大小尺寸
func (c Client) ScreenSize(serial string) (width int, height int, err error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return
	}

	// 写入命令
	_, err = conn.Write(EncodeCommend("shell:dumpsys window 2>/dev/null"))
	if err != err {
		return
	}

	stat := make([]byte, 4)
	_, err = conn.Read(stat)
	if err != nil {
		err = errors.New("adb response: Fail")
		return
	}

	if string(stat) == OKAY {
		rsRE := regexp.MustCompile(`\s*mRestrictedScreen=\(\d+,\d+\) (?P<w>\d+)x(?P<h>\d+)`)
		out, _ := ioutil.ReadAll(conn)
		matches := rsRE.FindStringSubmatch(string(out))
		if len(matches) == 0 {
			return 0, 0, errors.New("get shape(width,height) from device error")
		}
		return Atoi(matches[1]), Atoi(matches[2]), nil
	} else if string(stat) == FAIL {
		return 0, 0, errors.New("adb response: Fail")
	}
	return 0, 0, errors.New("adb response: " + string(stat))
}
