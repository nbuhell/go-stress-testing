/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"go-stress-testing/heper"
	"go-stress-testing/model"
	"go-stress-testing/server/client"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// http go link
func Http(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < totalNumber; i++ {

		var (
			startTime = time.Now()
			isSucceed = false
			errCode   = model.HttpOk
		)

		resp, err := client.HttpRequest(request.Method, request.Url, request.GetBody(), request.Headers, request.Timeout)
		requestTime := uint64(heper.DiffNano(startTime))
		// resp, err := server.HttpGetResp(request.Url)
		if err != nil {
			errCode = model.RequestErr // 请求错误
		} else {
			// 验证请求是否成功
			errCode, isSucceed = request.VerifyHttp(request, resp)
		}

		requestResults := &model.RequestResults{
			Time:      requestTime,
			IsSucceed: isSucceed,
			ErrCode:   errCode,
		}

		requestResults.SetId(chanId, i)

		ch <- requestResults
	}

	return
}

func HttpSession(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, sessionNumber int, sessionkey string, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	var sk []string = make([]string, sessionNumber)
	var cookie = request.Headers["Cookie"]
	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < totalNumber; i++ {
		var idx int
		if i < uint64(sessionNumber) {
			//idx = *(*int)(unsafe.Pointer(&i))
			idx = int(i)
		} else {
			idx = rand.Intn(sessionNumber)
		}
		var (
			startTime = time.Now()
			isSucceed = false
			errCode   = model.HttpOk
		)
		headers := request.Headers
		if sk[idx] != "" {
			t := cookie
			if t != "" {
				headers["Cookie"] = t + ";" + sk[idx]
			} else {
				headers["Cookie"] = sk[idx]
			}
			//fmt.Print(headers["Cookie"])
		}
		resp, err := client.HttpRequest(request.Method, request.Url, request.GetBody(), headers, request.Timeout)
		requestTime := uint64(heper.DiffNano(startTime))
		// resp, err := server.HttpGetResp(request.Url)
		if err != nil {
			errCode = model.RequestErr // 请求错误
		} else {
			// 验证请求是否成功
			errCode, isSucceed = request.VerifyHttp(request, resp)
		}
		if sk[idx] == "" && resp != nil {
			ck := resp.Header.Get("Set-Cookie")
			arr := strings.Split(ck, ";")
			sk[idx] = arr[0]
			//println(ck, sk);
		}
		requestResults := &model.RequestResults{
			Time:      requestTime,
			IsSucceed: isSucceed,
			ErrCode:   errCode,
		}

		requestResults.SetId(chanId, i)

		ch <- requestResults
	}

	return
}
