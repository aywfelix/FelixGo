package shell

import (
	"fmt"
	"runtime"
	"strings"
)

func UpdateSystemDate(dateTime string) bool {
	system := runtime.GOOS
	switch system {
	case "windows":
		{
			_, err1 := ShellExec(`date  ` + strings.Split(dateTime, " ")[0])
			_, err2 := ShellExec(`time  ` + strings.Split(dateTime, " ")[1])
			if err1 != nil || err2 != nil {
				fmt.Println("更新系统时间错误:请用管理员身份启动程序!")
				return false
			}
			return true
		}
	case "linux":
		{
			_, err1 := ShellExec(`date -s  "` + dateTime + `"`)
			if err1 != nil {
				fmt.Println("更新系统时间错误:", err1.Error())
				return false
			}
			return true
		}
	case "darwin":
		{
			_, err1 := ShellExec(`date -s  "` + dateTime + `"`)
			if err1 != nil {
				fmt.Println("更新系统时间错误:", err1.Error())
				return false
			}
			return true
		}
	}
	return false
}

// func main() {
// 	dateTime := "2022-01-05 16:43:00"
// 	UpdateSystemDate(dateTime)
// }
