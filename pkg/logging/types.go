// Copyright (C) 2025 ANSYS, Inc. and/or its affiliates.
// SPDX-License-Identifier: MIT
//
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package logging

import (
	"sync"

	"go.uber.org/zap"
)

// ContextKey defines the supported context keys.
type ContextKey string

const (
	InstructionGuid ContextKey = "instructionGuid"
	WorkflowId      ContextKey = "workflowId"
	WorkflowRunId   ContextKey = "workflowRunId"
	UserId          ContextKey = "userId"
	AdapterType     ContextKey = "adapterType"
	WatchFolderPath ContextKey = "watchFolderPath"
	WatchFilePath   ContextKey = "watchFilePath"
	ReaderGuid      ContextKey = "readerGuid"
	ClientGuid      ContextKey = "clientGuid"
	Action          ContextKey = "action"
	Rest_Call_Id    ContextKey = "restCallId"
	Rest_Call       ContextKey = "restCall"
	UserMail        ContextKey = "userMail"
)

// Initialize the global logger variable.
var Log loggerWrapper

// Initialize config variables
var ERROR_FILE_LOCATION string
var LOG_LEVEL string
var LOCAL_LOGS bool
var LOCAL_LOGS_LOCATION string
var DATADOG_LOGS bool
var DATADOG_SOURCE string
var DATADOG_STAGE string
var DATADOG_VERSION string
var DATADOG_SERVICE_NAME string
var DATADOG_API_KEY string
var DATADOG_LOGS_URL string
var DATADOG_METRICS bool
var DATADOG_METRICS_URL string

// Config represents the configuration for the logging package.
type Config struct {
	ErrorFileLocation string
	LogLevel          string
	LocalLogs         bool
	LocalLogsLocation string
	DatadogLogs       bool
	DatadogSource     string
	DatadogStage      string
	DatadogVersion    string
	DatadogService    string
	DatadogAPIKey     string
	DatadogLogsURL    string
	DatadogMetrics    bool
	DatadogMetricsURL string
}

// ContextMap represents a context for managing key-value pairs with specific context keys. It allows setting, retrieving,
// and copying context data associated with various keys.
type ContextMap struct {
	data sync.Map
}

// loggerWrapper represents a wrapper for the zap.Logger to provide custom logging functionality.
type loggerWrapper struct {
	lw *zap.Logger
}

// Point represents a data point in a time series metric.
type Point struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// Resource represents a named resource associated with a metric.
type Resource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Metric represents a time series metric.
type Metric struct {
	Metric    string     `json:"metric"`
	Type      int        `json:"type"`
	Points    []Point    `json:"points"`
	Resources []Resource `json:"resources"`
}

// Metrics represents a collection of metrics.
type Metrics struct {
	Series []Metric `json:"series"`
}
