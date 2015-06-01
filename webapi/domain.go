package main
//http://regoio.herokuapp.com/
//yum install jwhois
import (
	"log"
	"net/http"
	"regexp"
	"strings"
	//"encoding/json"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/millken/go-whois"
	"github.com/millken/godns"
)

var (
    regex_nameserver = regexp.MustCompile(`(nameserver|Name Server|NS [0-9]+\s*|Hostname\.+):\s*(.*)`)
)

type DomainResource struct {
	Name string
}

type ResponseJson struct {
	Status int `json:"status"`
	Info string `json:"info"`
	Data interface{} `json:"data"`
}

func (this DomainResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/domains").
		Doc("Manage domains").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/ns/{domain}").To(this.DomainNameServer).
		// docs
		Doc("get domain nameserver").
		Operation("DomainNameServer").
		Param(ws.PathParameter("domain", "master domain").DataType("string") ))

	ws.Route(ws.GET("/cname/{domain}").To(this.DomainCname).
		// docs
		Doc("get domain cname").
		Operation("DomainCname").
		Param(ws.PathParameter("domain", "domain").DataType("string") ))

	container.Add(ws)
}

// GET http://localhost:8080/domains/whois/google.com
//
func (this DomainResource) DomainNameServer(request *restful.Request, response *restful.Response) {
	var nsserver []string
	
	rjson := ResponseJson{
		Status : 0,
		Info: "不能获取到ns服务器",
		Data: []string{},
	}
	
	domain := request.PathParameter("domain")
	if len(domain) == 0 {
		response.WriteEntity(rjson)
		return
	}
	log.Printf("whois domain: %s", domain)
	body, err := whois.Whois(domain)
	if err != nil {
		log.Printf("get Whois Error : %s", err.Error())
		response.WriteEntity(rjson)
    } else {
        nameservers := regex_nameserver.FindAllStringSubmatch(body, -1)
		for _, nameserver := range nameservers {
		    nsserver = append(nsserver, strings.ToLower(nameserver[2]))
		}
		if len(nsserver) == 0  {
			options := &godns.LookupOptions{
				DNSServers: godns.ChinaDNSServers, Net: "udp"}	
				nsserver, err = godns.LookupNS(domain, options)
				if err != nil {
					log.Printf("get ns Error : %s", err.Error())
				}else{
					rjson.Status = 1
					rjson.Info = ""
					rjson.Data = nsserver
				}
				response.WriteEntity(rjson)
		}else{
			rjson.Status = 1
			rjson.Info = ""
			rjson.Data = nsserver
		    response.WriteEntity(rjson)
        }
    }	
	
}

func (this DomainResource) DomainCname(request *restful.Request, response *restful.Response) {
	
	rjson := ResponseJson{
		Status : 0,
		Info: "不能获取到cname",
		Data: []string{},
	}
	
	domain := request.PathParameter("domain")
	domain = strings.TrimPrefix(domain, "@.")
	domain = strings.Replace(domain, "*.", "h0d4lx.", 1)
	
	if len(domain) == 0 {
		response.WriteEntity(rjson)
		return
	}
	log.Printf("cname domain: %s", domain)
	options := &godns.LookupOptions{
		DNSServers: godns.ChinaDNSServers, Net: "udp"}	
	body, err := godns.LookupCNAME(domain, options)
	if err != nil {
		log.Printf("get cname Error : %s", err.Error())
		response.WriteEntity(rjson)
    } else {
		rjson.Status = 1
		rjson.Info = ""
		rjson.Data = body
        response.WriteEntity(rjson)
    }	
	
}

func main() {
	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

	wsContainer := restful.NewContainer()
	u := DomainResource{}
	u.Register(wsContainer)

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",
	}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
