/**
 * @Author: hexing
 * @Description:
 * @File:  service
 * @Version: 1.0.0
 * @Date: 20-4-28 下午6:55
 */

package v1

const (
	CallSuccess = 200
	CallFail    = 201
)

type CallReply struct {
	Id   string `json:"id"`
	Cmd  string `json:"cmd"`
	Data []byte `json:"data"`
}
