package Job

import (
	"github.com/nutrun/lentil"
	"Utils"
	"Config"
	"net/http"
	"strings"
	"encoding/base64"
	"io/ioutil"
	"time"
	"errors"
	"fmt"
)

type Jobsms struct {
}

func (j *Jobsms) Run() {
	for {	
		Utils.LogInfo("Jobsms delay %d Second", Config.GetLoopTime().Sms)
		time.Sleep(time.Second * Config.GetLoopTime().Sms)

	    beanstalkd, err := lentil.Dial(Config.GetBeanstalk().Server)
	    if err != nil {
	        Utils.LogPanicErr(err)
	    }else{
	    	err = beanstalkd.Use(Config.GetBeanstalk().MobileQueue)
	    }	

	    if err != nil {
	        Utils.LogPanicErr(err)
	    }else{
	    	for i := 0; i < 10; i++ {
				job,err := beanstalkd.PeekReady()
			    if err != nil {
			        //Utils.LogPanicErr(err)
			        break
			    }else{
			    	body := strings.SplitN(string(job.Body), "\t", 2)
			    	if len(body) == 2 {
			    		r,err := base64.StdEncoding.DecodeString(body[1])
			    		if err == nil {
			    			fmt.Printf("Job id: %d  \nmobile: %s \nbody: %s", job.Id, body[0], string(r))
			    			e := SendSms( body[0], string(r))
			    			if e != nil {
			    				Utils.LogPanicErr(e)
			    			}
			    		}    	
			    	}					
			    	beanstalkd.Delete(job.Id)
			    }	    		
	    	}			
		}

	}	
}

func SendSms(mobile, body string)(err error) {
	url := fmt.Sprintf(Config.GetUrl().SmsApi, mobile, body)
	resp, err := http.Get(url)
	if err != nil {
		return  err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return  err
	}else{
		if string(res) == "0" {
			return  nil
		}else{
			 return errors.New(string(res))
		}
	}
	return  nil
}