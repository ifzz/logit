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
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/03/03 14:58:21

package files

import (
	"errors"
	"os"
	"sync"
	"time"
)

// DurationRollingFile is a time sensitive file.
//
//  file := NewDurationRollingFile(time.Second, func(now time.Time) string {
//      return "D:/" + now.Format(formatOfTime) + ".txt"
//  })
//  defer file.Close()
//  file.Write([]byte("Hello!"))
//
// You can use it like using os.File!
type DurationRollingFile struct {

	// file points the writer which will be used this moment.
	file *os.File

	// directory is the target storing all created files.
	directory string

	// lastTime is the created time of current file above.
	lastTime time.Time

	// duration is the core field of this struct.
	// Every times currentTime - lastTime >= duration, the file will
	// roll to an entire new file for writing. This field should be always
	// larger than minDuration for some safe considerations. See minDuration.
	duration time.Duration

	// nameGenerator is for generating the name of every created file.
	// You can customize your format of filename by implementing this function.
	// Default is DefaultNameGenerator().
	nameGenerator NameGenerator

	// mu is a lock for safe concurrency.
	mu *sync.Mutex
}

const (
	// minDuration prevents io system from creating file too fast.
	// Default is one second.
	minDuration = 1 * time.Second
)

// NewDurationRollingFile creates a new duration rolling file.
// duration is how long did it roll to next file.
// nextFilename is a function for generating next file name.
// Every times rolling to next file will call nextFilename first.
// now is the created time of next file. Notice that duration's min value
// is one second. See minDuration.
func NewDurationRollingFile(directory string, duration time.Duration) *DurationRollingFile {

	// 防止时间间隔太小导致滚动文件时 IO 的疯狂蠕动
	if duration < minDuration {
		panic(errors.New("Duration is smaller than " + minDuration.String() + "\n"))
	}

	return &DurationRollingFile{
		directory:     directory,
		duration:      duration,
		nameGenerator: DefaultNameGenerator(),
		mu:            &sync.Mutex{},
	}
}

// rollingToNextFile will roll to next file for drf.
func (drf *DurationRollingFile) rollingToNextFile(now time.Time) {

	// 如果创建新文件发生错误，就继续使用当前的文件，等到下一次时间间隔再重试
	newFile, err := CreateFileOf(drf.nameGenerator.NextName(drf.directory, now))
	if err != nil {
		return
	}

	// 关闭当前使用的文件，初始化新文件
	drf.file.Close()
	drf.file = newFile
	drf.lastTime = now
}

// ensureFileIsCorrect ensures drf is writing to a correct file this moment.
func (drf *DurationRollingFile) ensureFileIsCorrect() {
	now := time.Now()
	if drf.file == nil || now.Sub(drf.lastTime) >= drf.duration {
		drf.rollingToNextFile(now)
	}
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
func (drf *DurationRollingFile) Write(p []byte) (n int, err error) {
	drf.mu.Lock()
	defer drf.mu.Unlock()

	// 确保当前文件对于当前时间点来说是正确的
	drf.ensureFileIsCorrect()
	return drf.file.Write(p)
}

// Close releases any resources using just moment.
// It returns error when closing.
func (drf *DurationRollingFile) Close() error {
	drf.mu.Lock()
	defer drf.mu.Unlock()
	return drf.file.Close()
}

// SetNameGenerator replaces drf.nameGenerator to newNameGenerator.
func (drf *DurationRollingFile) SetNameGenerator(newNameGenerator NameGenerator) {
	drf.mu.Lock()
	defer drf.mu.Unlock()
	drf.nameGenerator = newNameGenerator
}
