package flywheel

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/aluttik/go-crossplane"
)

// parsed from test/nginx.conf https://www.nginx.com/resources/wiki/start/topics/examples/full/
var parsedExample = `{"status":"ok","errors":[],"config":[{"file":"../../test/nginx.conf","status":"ok","errors":[],"parsed":[{"directive":"user","line":1,"args":["www","www"]},{"directive":"worker_processes","line":2,"args":["5"]},{"directive":"error_log","line":3,"args":["logs/error.log"]},{"directive":"pid","line":4,"args":["logs/nginx.pid"]},{"directive":"worker_rlimit_nofile","line":5,"args":["8192"]},{"directive":"events","line":7,"args":[],"block":[{"directive":"worker_connections","line":8,"args":["4096"]}]},{"directive":"http","line":11,"args":[],"block":[{"directive":"include","line":12,"args":["conf/mime.types"],"includes":[1]},{"directive":"include","line":13,"args":["proxy.conf"],"includes":[2]},{"directive":"include","line":14,"args":["fastcgi.conf"],"includes":[3]},{"directive":"index","line":15,"args":["index.html","index.htm","index.php"]},{"directive":"default_type","line":17,"args":["application/octet-stream"]},{"directive":"log_format","line":18,"args":["main","$remote_addr - $remote_user [$time_local]  $status ","\"$request\" $body_bytes_sent \"$http_referer\" ","\"$http_user_agent\" \"$http_x_forwarded_for\""]},{"directive":"access_log","line":21,"args":["logs/access.log","main"]},{"directive":"sendfile","line":22,"args":["on"]},{"directive":"tcp_nopush","line":23,"args":["on"]},{"directive":"server_names_hash_bucket_size","line":24,"args":["128"]},{"directive":"server","line":26,"args":[],"block":[{"directive":"listen","line":27,"args":["80"]},{"directive":"server_name","line":28,"args":["domain1.com","www.domain1.com"]},{"directive":"access_log","line":29,"args":["logs/domain1.access.log","main"]},{"directive":"root","line":30,"args":["html"]},{"directive":"location","line":32,"args":["~","\\.php$"],"block":[{"directive":"fastcgi_pass","line":33,"args":["127.0.0.1:1025"]}]}]},{"directive":"server","line":37,"args":[],"block":[{"directive":"listen","line":38,"args":["80"]},{"directive":"server_name","line":39,"args":["domain2.com","www.domain2.com"]},{"directive":"access_log","line":40,"args":["logs/domain2.access.log","main"]},{"directive":"location","line":43,"args":["~","^/(images|javascript|js|css|flash|media|static)/"],"block":[{"directive":"root","line":44,"args":["/var/www/virtual/big.server.com/htdocs"]},{"directive":"expires","line":45,"args":["30d"]}]},{"directive":"location","line":49,"args":["/"],"block":[{"directive":"proxy_pass","line":50,"args":["http://127.0.0.1:8080"]}]}]},{"directive":"upstream","line":54,"args":["big_server_com"],"block":[{"directive":"server","line":55,"args":["127.0.0.3:8000","weight=5"]},{"directive":"server","line":56,"args":["127.0.0.3:8001","weight=5"]},{"directive":"server","line":57,"args":["192.168.0.1:8000"]},{"directive":"server","line":58,"args":["192.168.0.1:8001"]}]},{"directive":"server","line":61,"args":[],"block":[{"directive":"listen","line":62,"args":["80"]},{"directive":"server_name","line":63,"args":["big.server.com"]},{"directive":"access_log","line":64,"args":["logs/big.server.access.log","main"]},{"directive":"location","line":66,"args":["/"],"block":[{"directive":"proxy_pass","line":67,"args":["http://big_server_com"]}]}]}]}]},{"file":"../../test/conf/mime.types","status":"ok","errors":[],"parsed":[{"directive":"types","line":1,"args":[],"block":[{"directive":"text/html","line":2,"args":["html","htm","shtml"]},{"directive":"text/css","line":3,"args":["css"]},{"directive":"text/xml","line":4,"args":["xml","rss"]},{"directive":"image/gif","line":5,"args":["gif"]},{"directive":"image/jpeg","line":6,"args":["jpeg","jpg"]},{"directive":"application/x-javascript","line":7,"args":["js"]},{"directive":"text/plain","line":8,"args":["txt"]},{"directive":"text/x-component","line":9,"args":["htc"]},{"directive":"text/mathml","line":10,"args":["mml"]},{"directive":"image/png","line":11,"args":["png"]},{"directive":"image/x-icon","line":12,"args":["ico"]},{"directive":"image/x-jng","line":13,"args":["jng"]},{"directive":"image/vnd.wap.wbmp","line":14,"args":["wbmp"]},{"directive":"application/java-archive","line":15,"args":["jar","war","ear"]},{"directive":"application/mac-binhex40","line":16,"args":["hqx"]},{"directive":"application/pdf","line":17,"args":["pdf"]},{"directive":"application/x-cocoa","line":18,"args":["cco"]},{"directive":"application/x-java-archive-diff","line":19,"args":["jardiff"]},{"directive":"application/x-java-jnlp-file","line":20,"args":["jnlp"]},{"directive":"application/x-makeself","line":21,"args":["run"]},{"directive":"application/x-perl","line":22,"args":["pl","pm"]},{"directive":"application/x-pilot","line":23,"args":["prc","pdb"]},{"directive":"application/x-rar-compressed","line":24,"args":["rar"]},{"directive":"application/x-redhat-package-manager","line":25,"args":["rpm"]},{"directive":"application/x-sea","line":26,"args":["sea"]},{"directive":"application/x-shockwave-flash","line":27,"args":["swf"]},{"directive":"application/x-stuffit","line":28,"args":["sit"]},{"directive":"application/x-tcl","line":29,"args":["tcl","tk"]},{"directive":"application/x-x509-ca-cert","line":30,"args":["der","pem","crt"]},{"directive":"application/x-xpinstall","line":31,"args":["xpi"]},{"directive":"application/zip","line":32,"args":["zip"]},{"directive":"application/octet-stream","line":33,"args":["deb"]},{"directive":"application/octet-stream","line":34,"args":["bin","exe","dll"]},{"directive":"application/octet-stream","line":35,"args":["dmg"]},{"directive":"application/octet-stream","line":36,"args":["eot"]},{"directive":"application/octet-stream","line":37,"args":["iso","img"]},{"directive":"application/octet-stream","line":38,"args":["msi","msp","msm"]},{"directive":"audio/mpeg","line":39,"args":["mp3"]},{"directive":"audio/x-realaudio","line":40,"args":["ra"]},{"directive":"video/mpeg","line":41,"args":["mpeg","mpg"]},{"directive":"video/quicktime","line":42,"args":["mov"]},{"directive":"video/x-flv","line":43,"args":["flv"]},{"directive":"video/x-msvideo","line":44,"args":["avi"]},{"directive":"video/x-ms-wmv","line":45,"args":["wmv"]},{"directive":"video/x-ms-asf","line":46,"args":["asx","asf"]},{"directive":"video/x-mng","line":47,"args":["mng"]}]}]},{"file":"../../test/proxy.conf","status":"ok","errors":[],"parsed":[{"directive":"proxy_redirect","line":1,"args":["off"]},{"directive":"proxy_set_header","line":2,"args":["Host","$host"]},{"directive":"proxy_set_header","line":3,"args":["X-Real-IP","$remote_addr"]},{"directive":"proxy_set_header","line":4,"args":["X-Forwarded-For","$proxy_add_x_forwarded_for"]},{"directive":"client_max_body_size","line":5,"args":["10m"]},{"directive":"client_body_buffer_size","line":6,"args":["128k"]},{"directive":"proxy_connect_timeout","line":7,"args":["90"]},{"directive":"proxy_send_timeout","line":8,"args":["90"]},{"directive":"proxy_read_timeout","line":9,"args":["90"]},{"directive":"proxy_buffers","line":10,"args":["32","4k"]}]},{"file":"../../test/fastcgi.conf","status":"ok","errors":[],"parsed":[{"directive":"fastcgi_param","line":1,"args":["SCRIPT_FILENAME","$document_root$fastcgi_script_name"]},{"directive":"fastcgi_param","line":2,"args":["QUERY_STRING","$query_string"]},{"directive":"fastcgi_param","line":3,"args":["REQUEST_METHOD","$request_method"]},{"directive":"fastcgi_param","line":4,"args":["CONTENT_TYPE","$content_type"]},{"directive":"fastcgi_param","line":5,"args":["CONTENT_LENGTH","$content_length"]},{"directive":"fastcgi_param","line":6,"args":["SCRIPT_NAME","$fastcgi_script_name"]},{"directive":"fastcgi_param","line":7,"args":["REQUEST_URI","$request_uri"]},{"directive":"fastcgi_param","line":8,"args":["DOCUMENT_URI","$document_uri"]},{"directive":"fastcgi_param","line":9,"args":["DOCUMENT_ROOT","$document_root"]},{"directive":"fastcgi_param","line":10,"args":["SERVER_PROTOCOL","$server_protocol"]},{"directive":"fastcgi_param","line":11,"args":["GATEWAY_INTERFACE","CGI/1.1"]},{"directive":"fastcgi_param","line":12,"args":["SERVER_SOFTWARE","nginx/$nginx_version"]},{"directive":"fastcgi_param","line":13,"args":["REMOTE_ADDR","$remote_addr"]},{"directive":"fastcgi_param","line":14,"args":["REMOTE_PORT","$remote_port"]},{"directive":"fastcgi_param","line":15,"args":["SERVER_ADDR","$server_addr"]},{"directive":"fastcgi_param","line":16,"args":["SERVER_PORT","$server_port"]},{"directive":"fastcgi_param","line":17,"args":["SERVER_NAME","$server_name"]},{"directive":"fastcgi_index","line":19,"args":["index.php"]},{"directive":"fastcgi_param","line":21,"args":["REDIRECT_STATUS","200"]}]}]}`

type dummyProvider struct{}

func (d dummyProvider) Override(_ context.Context, directive, path string) ([]string, error) {
	return []string{"dummy" + directive}, nil
}

func (d dummyProvider) Close() error {
	return nil
}

func TestWritePayloadTmp(t *testing.T) {
	var payload crossplane.Payload
	if err := json.Unmarshal([]byte(parsedExample), &payload); err != nil {
		t.Fatalf("failed to unmarshal test data: %v", err)
	}
	err := OverridePayload(context.Background(), &payload, dummyProvider{})
	if err != nil {
		t.Fatalf("failed to override payload: %v", err)
	}

	uFiles, err := WritePayloadTmp(&payload, &crossplane.BuildOptions{})
	if err != nil {
		t.Fatalf("failed to write files to tmp: %v", err)
	}

	var baseConfig string
	for _, uFile := range uFiles {
		if strings.HasPrefix(filepath.Base(uFile.OGName), "nginx") {
			baseConfig = uFile.Name()
			break
		}
	}

	// don't follow includes as we flatten the structure in WritePayloadTmp
	testPayload, err := crossplane.Parse(baseConfig, &crossplane.ParseOptions{
		SingleFile:             true,
		SkipDirectiveArgsCheck: true,
	})
	if err != nil {
		t.Fatalf("failed to parse tmp files: %v", err)
	}
	if testPayload.Status != "ok" {
		t.Logf("read payload: %+v", testPayload)
		t.Fatalf("failed to read payload")
	}

	for _, c := range testPayload.Config {
		for _, d := range c.Parsed {
			if len(d.Args) > 0 {
				expected := "dummy" + d.Directive
				if d.Args[0] != expected {
					t.Logf("failed on directive: %s (%v)", d.Directive, d)
					t.Errorf("expected '%v' got '%v'", expected, d.Args[0])
				}
			}
		}
	}
}

func TestOverrideDirective(t *testing.T) {
	directive := crossplane.Directive{
		Directive: "foo",
		Line:      1,
		Args:      []string{"hi", "mom"},
	}
	overrideDirective(context.Background(), &directive, dummyProvider{}, "")

	if !reflect.DeepEqual(directive.Args, []string{"dummyfoo"}) {
		t.Errorf("failed to modify args")
	}
}
