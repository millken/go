/*
var url='';url= '04620'  +url;url='539' +   url;url= '3610_'  + url;url=   '86136740'   +url;url= 'c0e8884'+ url;url= 'ecf11'+ url;url= '55f7'   +   url;url=  '0'+   url;url='db527f5'+ url;url='68'+ url;url=  '11c'  + url;url=  '02' +   url;url=   'unkey=' +url;url=  '/?yund' +   url;;window.location=url;
*/
/*www.yzsme.gov.cn
<html><body><script>function rgc_(gvfb_){var bedb_,bueb_=new Array(),rkl_="o\x04\x89\x03`\nH\xe5\xe7\x83\xd4\x8c\xe1N\xdbO\xee,E\x1eG\x0bc\xfd\x0eN\xcaO\x10@\xd4\x85\xe5D\xe7B8\x1eB)l\xf9Ic\xf5";for(bedb_=0;bedb_<rkl_.length;bedb_++)bueb_[bedb_]=rkl_.charCodeAt(bedb_);for(bedb_=41;;){if(bedb_<1)break;bueb_[bedb_]=(bueb_[bedb_]+bueb_[bedb_-1])&0xff;bedb_--;}bedb_=4;while(true){if(bedb_>42)break;bueb_[bedb_]=((((((bueb_[bedb_]>>1)|((bueb_[bedb_]<<7)&0xff))-112)&0xff)<<1)&0xff)|(((((bueb_[bedb_]>>1)|((bueb_[bedb_]<<7)&0xff))-112)&0xff)>>7);bedb_++;}bedb_=1;while(true){if(bedb_>41)break;bueb_[bedb_]=(((~((bueb_[bedb_]+213)&0xff))&0xff)+134)&0xff;bedb_++;}rkl_="";for(bedb_=1;bedb_<bueb_.length-1;bedb_++)if(bedb_%8)rkl_+=String.fromCharCode(bueb_[bedb_]^gvfb_);eval(rkl_);}rgc_(74);</script><br><br><br><center><h3><p>访问本页面，您的浏览器需要支持JavaScript</p><p>The browser needs JavaScript to continue</p></h3></center></body></html>
*/
package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

func main() {
	Otto := otto.New()

	Otto.Run(`
		var url='';url= '04620'  +url;url='539' +   url;url= '3610_'  + url;url=   '86136740'   +url;url= 'c0e8884'+ url;url= 'ecf11'+ url;url= '55f7'   +   url;url=  '0'+   url;url='db527f5'+ url;url='68'+ url;url=  '11c'  + url;url=  '02' +   url;url=   'unkey=' +url;url=  '/?yund' +   url;;window.location=url;

		// The value of abc is 4
	`)
	b,_ := Otto.Get("url")
	fmt.Printf("url=%q",b)
	Otto.Run(`
		function rgc_(gvfb_){var bedb_,bueb_=new Array(),rkl_="o\x04\x89\x03`\nH\xe5\xe7\x83\xd4\x8c\xe1N\xdbO\xee,E\x1eG\x0bc\xfd\x0eN\xcaO\x10@\xd4\x85\xe5D\xe7B8\x1eB)l\xf9Ic\xf5";
	`)
	b,_ := Otto.Get("rkl_")
	fmt.Printf("url=%q",b)
}
