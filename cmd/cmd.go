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

package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/aberic/gnomon"
	"github.com/aberic/lily"
	"github.com/aberic/lily/api"
	"github.com/getwe/figlet4go"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	confYmlPath string
	daemon      bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "检出lily的版本号",
	Long:  `print the version number of lily`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lily Version:v" + lily.Version)
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动lily，会有初始化操作",
	Long:  `start lily service`,
	Args: func(cmd *cobra.Command, args []string) error {
		fmt.Println("startCmd daemon", daemon)
		if daemon {
			fmt.Println("后台运行…")
		} else {
			fmt.Println("前端启动…")
		}
		if gnomon.String().IsEmpty(confYmlPath) {
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

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "使用lily指定名称的数据库",
	Long:  `uses a database with the specified name`,
	Run: func(cmd *cobra.Command, args []string) {
		use()
	},
}

var rootCmd = &cobra.Command{
	Use:   "lily",
	Short: "lily是命令的抬头符",
	Long:  `lily is a cli library db. use lily can operation db, like start or stop.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here
		if len(args) < 1 {
			return errors.New("command is required , Use lily -h to get more information ")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func start() {
	conf := lily.ObtainConf(confYmlPath)
	fmt.Println("start daemon", daemon)
	if daemon {
		var (
			command *exec.Cmd
			pid     int
		)
		fmt.Println("确认后台运行…")
		fmt.Println("初始化数据库…")
		lily.ObtainLily().Start()
		fmt.Println("启动监听器…")
		if gnomon.String().IsEmpty(confYmlPath) {
			command = exec.Command("./lily", "start")
		} else {
			command = exec.Command("./lily", "start", "-p", confYmlPath)
		}
		_ = command.Start()
		pid = command.Process.Pid
		fmt.Printf("lily start, [PID] %d running...\n", pid)
		_ = ioutil.WriteFile("lily.lock", []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
		daemon = false
		var (
			running = false
			arr     []string
			err     error
		)
		for !running {
			if _, _, arr, err = gnomon.Command().ExecCommandSilent("lsof", "-i"); nil != err {
				panic(err)
			}
			for _, str := range arr {
				str = gnomon.String().SingleSpace(str)
				strs := strings.Split(str, " ")
				if strs[0] == "lily" && strs[1] == strconv.Itoa(pid) {
					running = true
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
		os.Exit(0)
	} else {
		fmt.Println("lily start")
	}
	rpcListener(conf)
}

func stop() {
	data, err := ioutil.ReadFile("lily.lock")
	if nil != err {
		panic(errors.New("lily haven not been started or no such file or directory with name lily.lock"))
	}
	_, _, _, err = gnomon.Command().ExecCommandTail("kill", string(data))
	if nil != err {
		panic(err)
	}
	println("lily stop")
}

func rpcListener(conf *lily.Conf) {
	var (
		listener net.Listener
		err      error
	)

	fmt.Println(strings.Join([]string{"Listen announces on the local network address with port: ", conf.Port}, ""))
	if listener, err = net.Listen("tcp", strings.Join([]string{":", conf.Port}, "")); nil != err {
		panic(err)
	}
	fmt.Println("creates a gRPC server")
	server := grpc.NewServer()
	fmt.Println("register gRPC listener")
	api.RegisterLilyAPIServer(server, &lily.APIServer{})
	fmt.Println("OFF")
	if err = server.Serve(listener); nil != err {
		panic(err)
	}
}

func use() {
	var (
		sql string
		err error
	)
	fmt.Print("lily->: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sql = scanner.Text()
		if gnomon.String().TrimN(sql) == "" {
			fmt.Print("lily->: ")
			continue
		}
		if sql == "exit" {
			fmt.Println("Bye!")
			os.Exit(0)
		}
		gnomon.Log().Debug("use", gnomon.Log().Field("sql", sql))
		if err = analysis(sql); nil != err {
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
	rootCmd.AddCommand(useCmd)
	startCmd.Flags().StringVarP(&confYmlPath, "path", "p", "", "也许你希望通过指定‘conf.yml’文件来使用自己的配置.")
	startCmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "是否启动后台运行")
}

// Execute cmd start
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		//zap.S().Debug(err)
		os.Exit(1)
	}
}
