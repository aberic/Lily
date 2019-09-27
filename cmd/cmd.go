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
	"errors"
	"fmt"
	"github.com/aberic/gnomon"
	"github.com/aberic/lily"
	"github.com/aberic/lily/api"
	"github.com/aberic/lily/io"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
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
		if gnomon.String().IsEmpty(confYmlPath) {
			fmt.Println("lily 数据库将使用默认配置策略")
		}
		if daemon {
			fmt.Println("后台运行…")
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
		fmt.Println("lily use")
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
	fmt.Println("daemon", daemon)
	if daemon {
		fmt.Println("确认后台运行…")
		command := exec.Command("./lily", "start", "-p", confYmlPath, "-d false")
		if err := command.Start(); nil != err {
			fmt.Println("start error: ", err.Error())
		}
		fmt.Printf("lily start, [PID] %d running...\n", command.Process.Pid)
		_ = ioutil.WriteFile("lily.lock", []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
		daemon = false
		os.Exit(0)
	} else {
		fmt.Println("lily start")
	}
	fmt.Println("初始化数据库…")
	lily.ObtainLily().Start()
	fmt.Println("启动监听器…")
	rpcListener(conf)
	fmt.Println("完成！")
}

func stop() {
	data, _ := ioutil.ReadFile("lily.lock")
	command := exec.Command("kill", string(data))
	_ = command.Start()
	println("lily stop")
}

func rpcListener(conf *lily.Conf) {
	var (
		listener net.Listener
		err      error
	)

	if listener, err = net.Listen("tcp", strings.Join([]string{":", conf.Port}, "")); nil != err {
		panic(err)
	}
	server := grpc.NewServer()
	api.RegisterLilyAPIServer(server, &io.LilyAPIServer{})
	if err = server.Serve(listener); nil != err {
		panic(err)
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
