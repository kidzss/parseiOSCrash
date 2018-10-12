package main

import(
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

func init()  {
	GetSymbolicatecrashTool()
}

func CrashFileValid(excuteFile string,dsymFile string) bool {
	if len(excuteFile) < 4 {
		return false
	}

	excuteCmd := excuteFile + "/" + excuteFile[0:len(excuteFile) - 4]
	dsymCmd := dsymFile

	var excListStr []string
	var dsymListStr []string

	//获取app的udid
	cmd := exec.Command("dwarfdump","--uuid",excuteCmd)
	if excResult, err := cmd.Output(); err != nil {
		fmt.Println()
		panic(err)
	} else {
		fmt.Println(string(excResult))
		excListStr = getUdidList(string(excResult))
	}

	//获取dsym的udid
	cmd = exec.Command("dwarfdump","--uuid",dsymCmd)
	if dsymResult, err := cmd.Output(); err != nil {
		panic(err)
	} else {
		dsymListStr = getUdidList(string(dsymResult))
	}
	var isExist bool = true
	//比较两个数据是否相等
	for _,v := range excListStr {
		for _,val:= range dsymListStr {
			if string(v) != string(val) {
				isExist = false
				goto Loop
			}
		}
	}
	Loop:
	return isExist
}

//获取命令字符串中的udid
func getUdidList(cmdResult string) []string {
	var resultList []string
	excTmpList1 := strings.Split(cmdResult,"\n")
	for _,v := range excTmpList1 {
		if len(v) > 2 {
			excTmp2 := strings.Split(v," ")
			if cap(excTmp2) > 1 {
				resultList = append(resultList, excTmp2[1])
			}
		}
	}
	return resultList
}

//解析设置环境
func SetEnvironment() {
	err := os.Setenv("DEVELOPER_DIR","/Applications/Xcode.app/Contents/Developer")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("DEVELOPER_DIR %s \n",os.Getenv("DEVELOPER_DIR"))
	}
}

//检测是否存在symbolicatecrash工具，如果不存在，则在xcode拷贝到当前目录
func GetSymbolicatecrashTool() {
	currentDir,err := os.Getwd()
	if err != nil {
		panic(err)
	}

	symbolDir := "/Applications/Xcode.app/Contents/SharedFrameworks/DVTFoundation.framework/Versions/A/Resources/symbolicatecrash"
	exist,err := PathExists(currentDir + "/" + "symbolicatecrash")
	if (err == nil)&& exist == true {
		SetEnvironment()
	} else {
		if err != nil {
			fmt.Printf("err %s \n",err)
		}

		cmd := exec.Command("cp","",symbolDir,currentDir)
		if excResult, err := cmd.Output(); err != nil {
			panic(err)
		} else {
			fmt.Println(string(excResult))
		}

		SetEnvironment()
	}
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//遍历当前目录，解析文件
func ParseCrashLog() {
	var crashList []string
	var dsymName string
	var appName  string

	currentDir,err := os.Getwd()
	if err != nil {
		panic(err)
	}
	files, _ := ioutil.ReadDir(currentDir)
	
	for _, f := range files {
		fileSuffix := path.Ext(f.Name())
		if fileSuffix == ".crash" || fileSuffix == ".ips" {
			crashList = append(crashList, f.Name())
		}
		if fileSuffix == ".dSYM" {
			dsymName = f.Name()
		}
		if fileSuffix == ".app" {
			appName = f.Name()
		}
	}

	//文件检查
	if dsymName == "" {
		fmt.Printf("error %s","请将dsym文件考到当前目录!!!!!")
		return
	}

	if appName == "" {
		fmt.Printf("error %s","请将app文件考到当前目录!!!!")
		return
	}

	//检测dsym文件 和app文件是否匹配
	if CrashFileValid(appName, dsymName) == false {
		fmt.Printf("error %s \n","app和dsym文件不匹配")
		return
	}
	for _,crashFile := range crashList {
		ext := path.Ext(crashFile)

		if ext == ".crash" {
			//cmdStr := "./symbolicatecrash "+ crashFile + " " + dsymName
			cmdStr := "./symbolicatecrash "+ crashFile + " " + dsymName + ">>" + crashFile[0:len(crashFile) - 6]+".log"
			cmd := exec.Command("/bin/bash", "-c", cmdStr)

			//cmd.Stdout = os.Stdout
			//cmd.Run()
			//if cmd.Stdout != nil {
			//	fmt.Println(cmd.Start())
			//}

			if excResult, err := cmd.Output(); err != nil {
				panic(err)
			} else {
				fmt.Println(string(excResult))
			}

			//stdout, err := cmd.StdoutPipe() //指向cmd命令的stdout
			//cmd.Start()
			//content, err := ioutil.ReadAll(stdout)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//fmt.Println(string(content))     //输出命令查看到的内容

		} else if ext == ".ips" {
			//cmdStr := "./symbolicatecrash "+ crashFile + " " + dsymName
			cmdStr := "./symbolicatecrash "+ crashFile + " " + dsymName + ">>" + crashFile[0:len(crashFile) - 4]+".log"
			cmd := exec.Command("/bin/bash", "-c", cmdStr)

			if excResult, err := cmd.Output(); err != nil {
				panic(err)
			} else {
				fmt.Println(string(excResult))
			}

			//stdout, err := cmd.StdoutPipe() //指向cmd命令的stdout
			//cmd.Start()
			//content, err := ioutil.ReadAll(stdout)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//fmt.Println(string(content))     //输出命令查看到的内容
		}
	}
}

func main() {
	ParseCrashLog()
}