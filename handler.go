// Copyright 2020 Ye Zi Jie. All Rights Reserved.
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
//
// Author: fish
// Email: fishinlove@163.com
// Created at 2020/03/06 13:36:28

package logit

import (
    "io"
)

// Handler is an interface representation of log handler.
// Every log will be handled by handler, and you can customize your own handler
// to handle logs in your way. The return value is meaningful, false means
// next handler will not be used, only true will go on handling process.
// Notice that if one handler returns false, then all handlers after it
// will not be used anymore.
type Handler interface {
    Handle(log *Log) bool
}

type DefaultHandler struct {
    writer  io.Writer
    encoder Encoder
}

func NewDefaultHandler(writer io.Writer, encoder Encoder) Handler {
    return &DefaultHandler{
        writer:  writer,
        encoder: encoder,
    }
}

func (dh *DefaultHandler) Handle(log *Log) bool {
    dh.writer.Write(dh.encoder.Encode(log))
    return true
}