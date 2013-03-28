package Utils

import (
	//"errors"
	"fmt"
)

var err map[int]string
type TaskErr struct{
	ErrCode int
	ErrMsg string
}
func (e TaskErr) Error() string {
    return fmt.Sprintf("%v:%v", e.ErrCode, e.ErrMsg)
}

func init() {
	err = map[int]string{
		/******************工具错误***************/
		100:"data_stream_parse_err", //数据流解析错误
		101:"decode_json_err",	//解码json错误
		102:"encode_json_err",	//加码json错误
		103:"parse_config_err",	//解析配置文件错误
		104:"config_notfound",	//配置文件不存在
		/*******************task错误***********************/
		1000:"del_task_nofound",	//要删除的任务没找到
		1001:"add_expired_task_nofound",//要增加到期时间的任务没找到
	}
}


func LogErr(errCode interface{}) error {
	switch value := errCode.(type) {
	case int:

		errMsg, exist := err[value]
		if exist == true {
			LogInfo("err=%s\n",fmt.Sprintf("errCode=%d,errMsg=%s", value, errMsg))
			return TaskErr{value, errMsg}
		} else {
			LogInfo("err=%s\n","unknow err")
			return TaskErr{1000, "unknow err"}
		}
	case error:
		LogInfo("err=%s\n",errCode)
		return errCode.(error)
	}
	return nil

}
func LogPanicErr(err interface{}) {
	LogInfo("\npanic=%s\n",err)
}

