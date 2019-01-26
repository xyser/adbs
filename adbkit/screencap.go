package adbkit

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Screencap 用于获取屏幕截图(未压缩)
func (c Client) Screencap(serial string) ([]byte, error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return nil, err
	}

	// 写入命令
	_, err = conn.Write(EncodeCommend("shell:screencap -p 2>/dev/null"))
	if err != err {
		return nil, err
	}

	stat := make([]byte, 4)
	_, err = conn.Read(stat)
	if err != nil {
		return nil, errors.New("adb response: Fail")
	}

	if string(stat) == OKAY {
		buf := bytes.NewBuffer([]byte{})
		for {
			temp := make([]byte, 1024)
			n, err := conn.Read(temp)
			if err == io.EOF {
				break
			}
			if n > 0 {
				buf.Write(temp[0:n])
			}
		}

		con := []byte(strings.Replace(buf.String(), "\x0D\x0A", "\x0A", -1))
		err = ioutil.WriteFile("/go/adbs/1.png", con, 0644)
		return con, err
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
	//_, err = conn.Write(EncodeCommend("shell:cat /sdcard/a.png"))
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
			err = errors.New("get shape(width,height) from device error")
			return
		}
		width, _ = strconv.Atoi(matches[1])
		height, _ = strconv.Atoi(matches[2])
		return
	}
	return
}
