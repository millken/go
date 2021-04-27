#!/bin/sh
if [ -n "$1" ]; then
	server=$1
else
	server=""
fi
if [ -n "$2" ]; then
	start_date=`date  --date="$2 -1 day" +%Y%m%d`
else
	start_date=""
fi
if [ -n "$3" ]; then
	end_date=$3
else
	end_date=""
fi

while  [[ $start_date != $end_date ]]
do
start_date=`date  --date="$start_date 1 day" +%Y%m%d`
y=`date  --date="$start_date" +%Y`
m=`date  --date="$start_date" +%m`
d=`date  --date="$start_date" +%d`

curl -XDELETE "http://$server:9200/nginx-$y-$m-$d/"
curl -XPOST "http://$server:9200/nginx-$y-$m-$d" -d '
{
    "mappings" : {
      "nginx.access" : {
        "properties" : {
          "EnvVersion" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "Hostname" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "Logger" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "Payload" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "Pid" : {
            "type" : "long"
          },
          "Severity" : {
            "type" : "long"
          },
          "Timestamp" : {
            "type" : "date",
            "format" : "dateOptionalTime"
          },
          "Type" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "Uuid" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "body_bytes_sent" : {
            "type" : "long",
            "index": "not_analyzed"
          },
          "country" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "http_referer" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "http_user_agent" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "isp" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "province" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "remote_addr" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "user_agent_browser" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "user_agent_os" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "user_agent_version" : {
            "type" : "long",
            "index": "not_analyzed"
          },
          "server_addr" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "request" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "http_host" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "method" : {
            "type" : "string",
            "index": "not_analyzed"
          },
          "path" : {
            "type" : "string",
            "index": "not_analyzed"
          },                   
          "status" : {
            "type" : "long",
            "index": "not_analyzed"
          }
        }
      }
    }
}
' -o /dev/null
done
