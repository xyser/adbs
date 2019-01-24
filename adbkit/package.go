package adbkit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type Package struct {
	ApplicationId string `json:"application_id"`
	ApkPath       string `json:"apk_path"`
}

// 清理包缓存
// fmt.Println(adbkit.New("127.0.0.1", 5037).Reboot("emulator-5554"))
func (c Client) Clear(serial, pkg string) error {
	conn, err := c.Transport(serial)
	if err != nil {
		return err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command := "shell:pm clear " + pkg
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return err
	}

	resp := <-readChan

	if string(resp[0:4]) == OKAY {
		return nil
	} else if string(resp[0:4]) == FAIL {
		return errors.New("adb fail response: " + string(resp[4:]))
	}

	return errors.New("adb response: " + string(resp))
}

// 获取包名
// fmt.Println(adbkit.New("127.0.0.1", 5037).List("emulator-5554"))
func (c Client) List(serial string) (packages []Package, err error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return nil, err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command := "shell:pm list packages -f 2>/dev/null"
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return nil, err
	}

	resp := <-readChan
	if string(resp[0:4]) == OKAY {
		var re = regexp.MustCompile(`(?m)(\S+):(\S+)=(\S+)`)
		for _, match := range re.FindAllStringSubmatch(string(resp[0:4]), -1) {
			if len(match) >= 4 {
				packages = append(packages, Package{ApkPath: match[2], ApplicationId: match[3]})
			}
		}
		return packages, nil
	} else if string(resp[0:4]) == FAIL {
		return nil, errors.New("adb fail response: " + string(resp[4:]))
	}

	return nil, errors.New("adb response: " + string(resp))
}

func (c Client) Features(serial string) (features []string, err error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return nil, err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command := "shell:pm list features 2>/dev/null"
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return nil, err
	}

	resp := <-readChan

	if string(resp[0:4]) == OKAY {
		for _, lint := range strings.Split(string(resp[4:]), "\n") {
			feature := strings.Split(lint, ":")
			if len(feature) > 1 {
				features = append(features, feature[1])
			}
		}
		return features, nil
	} else if string(resp[0:4]) == FAIL {
		return nil, errors.New("adb fail response: " + string(resp[4:]))
	}

	return nil, errors.New("adb response: " + string(resp))
}

// 获取包的路径
// path,err := adbkit.New("127.0.0.1", 5037).GetPath("emulator-5554", "com.android.smoketest")
func (c Client) GetPath(serial, pkg string) (path string, err error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return "", err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command := fmt.Sprintf("shell:pm path %s 2>/dev/null", pkg)
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return "", err
	}

	resp := <-readChan
	if string(resp[0:4]) == OKAY {
		if len(resp) > 12 {
			return strings.TrimRight(string(resp[12:]), "\n"), nil
		}
		return "", errors.New(string(resp[4:]))
	} else if string(resp[0:4]) == FAIL {
		return "", errors.New("adb fail response: " + string(resp[4:]))
	}

	return "", errors.New("adb response: " + string(resp))
}

// 安装远端应用
func (c Client) Install(serial, path string) (bool, error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return false, err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command := fmt.Sprintf("shell:pm install -r %s", path)
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return false, err
	}

	resp := <-readChan
	if string(resp[0:4]) == OKAY {
		switch true {
		case strings.Contains(string(resp[4:]), "Success"):
			return true, nil
		case strings.Contains(string(resp[4:]), "Failure"):
			return false, errors.New("adb fail response: " + string(resp[4:]))
		}
		return false, errors.New(string(resp[4:]))
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb fail response: " + string(resp[4:]))
	}

	return false, errors.New("adb response: " + string(resp))
}

func (c Client) UnInstall(serial, pkg string) (bool, error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return false, err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command := fmt.Sprintf("shell:pm uninstall %s", pkg)
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return false, err
	}

	resp := <-readChan
	if string(resp[0:4]) == OKAY {
		switch true {
		case strings.Contains(string(resp[4:]), "Success"):
			return true, nil
		case strings.Contains(string(resp[4:]), "Failure"):
			return false, errors.New("adb fail response: " + string(resp[4:]))
		case strings.Contains(string(resp[4:]), "Unknown package"):
			return false, errors.New("unknown package: " + pkg)
		}
		return false, errors.New(string(resp[4:]))
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb fail response: " + string(resp[4:]))
	}

	return false, errors.New("adb response: " + string(resp))
}

// 打开 Activity or Service
func (c Client) Start(serial, command string, args map[string]string) (bool, error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return false, err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command = fmt.Sprintf("shell:am %s", command)
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return false, err
	}

	resp := <-readChan
	if string(resp[0:4]) == OKAY {
		var re = regexp.MustCompile(`(?m)^Error: (.*)$`)
		match := re.FindStringSubmatch(string(resp[0:4]))
		if len(match) > 1 {
			return false, errors.New(match[1])
		}
		return true, nil
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb fail response: " + string(resp[4:]))
	}

	return false, errors.New("adb response: " + string(resp))
}
