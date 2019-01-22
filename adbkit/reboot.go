package adbkit

func (c Client) Reboot(serial string) bool {
	change, _ := c.Transport(serial)
	// TODO:: 接下来的问题是切换连接态
	if change {
		return true
	}
	return false
}
