// Copyright 2023 PingCAP, Inc.
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

package dispatcher

import (
	"context"

	"github.com/pingcap/tidb/disttask/framework/proto"
	"github.com/pingcap/tidb/util/syncutil"
	"golang.org/x/exp/maps"
)

// TaskFlowHandle is used to control the process operations for each global task.
type TaskFlowHandle interface {
	ProcessNormalFlow(ctx context.Context, h TaskHandle, gTask *proto.Task) (subtaskMetas [][]byte, err error)
	ProcessErrFlow(ctx context.Context, h TaskHandle, gTask *proto.Task, receiveErr [][]byte) (subtaskMeta []byte, err error)
	IsRetryableErr(err error) bool
}

var taskFlowHandleMap struct {
	syncutil.RWMutex
	handleMap map[string]TaskFlowHandle
}

// RegisterTaskFlowHandle is used to register the global task handle.
func RegisterTaskFlowHandle(taskType string, dispatcherHandle TaskFlowHandle) {
	taskFlowHandleMap.Lock()
	taskFlowHandleMap.handleMap[taskType] = dispatcherHandle
	taskFlowHandleMap.Unlock()
}

// ClearTaskFlowHandle is only used in test
func ClearTaskFlowHandle() {
	taskFlowHandleMap.Lock()
	maps.Clear(taskFlowHandleMap.handleMap)
	taskFlowHandleMap.Unlock()
}

// GetTaskFlowHandle is used to get the global task handle.
func GetTaskFlowHandle(taskType string) TaskFlowHandle {
	taskFlowHandleMap.Lock()
	defer taskFlowHandleMap.Unlock()
	return taskFlowHandleMap.handleMap[taskType]
}

func init() {
	taskFlowHandleMap.handleMap = make(map[string]TaskFlowHandle)
}
