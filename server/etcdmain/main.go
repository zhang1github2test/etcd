// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdmain

import (
	"fmt"
	"os"

	"github.com/coreos/go-systemd/v22/daemon"
	"go.uber.org/zap"
)

func Main(args []string) {
	// 1. 检查系统架构支持
	checkSupportArch()

	if len(args) > 1 { // // 2. 检查命令行参数
		cmd := args[1]
		switch cmd {
		case "gateway", "grpc-proxy": // // 3. 判断是否运行网关或 gRPC 代理模式
			if err := rootCmd.Execute(); err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}

	startEtcdOrProxyV2(args) // // 4. 启动普通 ETCD 节点或代理节点
}

func notifySystemd(lg *zap.Logger) {
	lg.Info("notifying init daemon")
	_, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if err != nil {
		lg.Error("failed to notify systemd for readiness", zap.Error(err))
		return
	}
	lg.Info("successfully notified init daemon")
}
