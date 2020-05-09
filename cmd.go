/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lily

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/aberic/gnomon"
	"github.com/getwe/figlet4go"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	confYmlPath string // confYmlPath lily配置文件地址
	daemon      bool   // daemon 是否后台启动
	address     string // address lily服务地址
	username    string // username lily服务用户名
	password    string // password lily服务密码
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "检出lily的版本号",
	Long:  `print the version number of lily`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lily Version:v" + Version)
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动lily，会有初始化操作",
	Long:  `start lily service`,
	Args: func(cmd *cobra.Command, args []string) error {
		//fmt.Println("startCmd daemon", daemon)
		if daemon {
			fmt.Println("后台启动...")
		} else {
			fmt.Println("前端启动...")
		}
		if gnomon.StringIsEmpty(confYmlPath) {
			fmt.Println("lily 数据库将使用默认配置策略")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止lily",
	Long:  `stop lily service`,
	Run: func(cmd *cobra.Command, args []string) {
		stop()
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "重新启动lily，如果是首次启动，则会执行初始化操作，如果不是，则尝试加载旧数据",
	Long:  `Restart the lily, and if it is the first time, initialization will be performed, and if it is not, an attempt will be made to load the old data`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lily restart")
	},
}

var connCmd = &cobra.Command{
	Use:   "conn",
	Short: "使用lily指定名称的数据库",
	Long:  `uses a database with the specified name`,
	Args: func(cmd *cobra.Command, args []string) error {
		fmt.Println("address", address)
		fmt.Println("username", username)
		fmt.Println("password", password)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conn()
	},
}

var rootCmd = &cobra.Command{
	Use:   "lily",
	Short: "lily是命令的抬头符",
	Long:  `lily is a cli library db. use lily can operation db, like start or stop.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("command is required , Use lily -h to get more information ")
		}
		switch args[0] {
		default:
			return errors.New("command is required , Use lily -h to get more information ")
		case "conn", "help", "restart", "start", "stop", "version":
			return nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func start() {
	if gnomon.FilePathExists("lock") {
		fmt.Println("lily is start already")
		return
	}
	conf := ObtainConf(confYmlPath)
	//fmt.Println("start daemon", daemon)
	if daemon {
		var (
			command *exec.Cmd
			pid     int
		)
		fmt.Println("启动中...")
		if gnomon.StringIsEmpty(confYmlPath) {
			command = exec.Command("./lily", "start")
		} else {
			command = exec.Command("./lily", "start", "-p", confYmlPath)
		}
		_ = command.Start()
		pid = command.Process.Pid
		fmt.Printf("lily start, [PID] %d running...\n", pid)
		daemon = false
		var (
			running = false
			arr     []string
			err     error
		)
		loadChan := make(chan struct{})
		go loadingFmt(time.Second, loadChan)
		for !running {
			if _, _, arr, err = gnomon.CommandExecSilent("lsof", "-i"); nil != err {
				panic(err)
			}
			for _, str := range arr {
				str = gnomon.StringSingleSpace(str)
				strs := strings.Split(str, " ")
				if strs[0] == "lily" && strs[1] == strconv.Itoa(pid) {
					loadChan <- struct{}{}
					running = true
					fmt.Println()
					fmt.Println("------------------------------------------------------------")
					flag.Parse()
					str := *flag.String("str", "Lily", "input string")
					ascii := figlet4go.NewAsciiRender()
					// most simple Usage
					renderStr, _ := ascii.Render(str)
					fmt.Println(renderStr)
					fmt.Println("------------------------------------------------------------")
					fmt.Println("lily start success")
				}
			}
		}
		_ = ioutil.WriteFile("lily.lock", []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
		os.Exit(0)
	} else {
		fmt.Println("lily start")
	}
	fmt.Println("初始化数据库…")
	ServerStart(conf)
}

func loadingFmt(delay time.Duration, loadChan chan struct{}) {
	s := "."
	sFull := ".     "
	running := false
	go func() {
		<-loadChan
		running = true
	}()
	for {
		if running {
			return
		}
		switch len(s) {
		default:
			s = sFull
		case 1, 2, 3, 4, 5:
			s = strings.Join([]string{s, "."}, "")
		}
		fmt.Printf("\r%s", s)
		if s == sFull {
			s = "."
		}
		time.Sleep(delay)
	}
}

func stop() {
	data, err := ioutil.ReadFile("lily.lock")
	if nil != err {
		panic(errors.New("lily haven not been started or no such file or directory with name lock"))
	}
	_, _, _, err = gnomon.CommandExecTail("kill", string(data))
	if nil != err {
		panic(err)
	}
	_ = os.Remove("lily.lock")
	println("lily stop")
}

// conn 数据库连接
func conn() {
	var (
		sqlContent string
		s          *sql
		err        error
	)
	if gnomon.StringIsEmpty(address) {
		fmt.Println("connection to default gRPC server 'localhost:19877'")
		s = &sql{serverURL: "localhost:19877"}
	} else {
		s = &sql{serverURL: address}
	}
	fmt.Print("lily->: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sqlContent = scanner.Text()
		if gnomon.StringTrimN(sqlContent) == "" {
			fmt.Print("lily->: ")
			continue
		}
		if gnomon.StringTrim(sqlContent) == "exit" {
			fmt.Println("Bye!")
			os.Exit(0)
		}
		//gnomon.Log().Debug("use", gnomon.Log().Field("sqlContent", sqlContent))
		if err = s.analysis(sqlContent); nil != err {
			fmt.Println(err.Error())
		}
		fmt.Print("lily->: ")
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(connCmd)
	startCmd.Flags().StringVarP(&confYmlPath, "path", "p", "", "也许你希望通过指定‘conf.yml’文件来使用自己的配置.")
	startCmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "是否启动后台运行")
	connCmd.Flags().StringVarP(&address, "address", "a", "localhost:19877", "lily服务端地址，默认localhost")
	connCmd.Flags().StringVarP(&username, "username", "u", "", "lily服务端登录用户")
	connCmd.Flags().StringVarP(&password, "password", "p", "", "lily服务端登录密码")
}

// Execute cmd start
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
