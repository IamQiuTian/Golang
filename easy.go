package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	download  *bool   = flag.Bool("download", false, "download file")
	update    *bool   = flag.Bool("update", false, "update file")
	directory *string = flag.String("d", "false", "Directory path")
	file      *string = flag.String("f", "false", "file path")
	port      *string = flag.String("p", "8888", "Listening port")
	pwd       *string = flag.String("pwd", "nil", "password")
)

var (
	randomstr   = randomStr(28)
	randomstrup = randomStr(28)
	publicIP    = getPublic()
	privateIP   = getPrivate()
)

var filename string

func main() {
	flag.Parse()
	// 参数为空或参数都不为空
	if *download == false && *update == false || *download != false && *update != false {
		flag.Usage()
		return
	}
	if *update == true && *directory == "false" {
		flag.Usage()
		return
	}
	if *directory == "false" && *file == "false" || *directory != "false" && *file != "false" {
		flag.Usage()
		return
	}

	if *download {
		fileDown()
	} else {
		updateFile()
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

// 文件上传控制
func updateFile() {
	ok, _ := pathExist(*directory)
	if !ok {
		log.Fatal("Directory does not exist")
	}

	http.HandleFunc(fmt.Sprintf("/%s", randomstrup), func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"%s  %s  %s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
		)

		fileup, fileupinfo, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("upload error!"))
			return
		}

		fileupname := fileupinfo.Filename
		filewn, err := os.OpenFile(fmt.Sprintf("%s/%s", *directory, fileupname), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			w.Write([]byte("upload error!"))
			return
		}

		defer func() {
			if err := fileup.Close(); err != nil {
				log.Fatal("Close: ", err.Error())
				return
			}
			if err := filewn.Close(); err != nil {
				log.Fatal("Close: ", err.Error())
				return
			}
		}()

		_, err = io.Copy(filewn, fileup)
		if err != nil {
			w.Write([]byte("file update error"))
		}
		w.Write([]byte("file update success"))
	})

	if *pwd != "nil" {
		http.HandleFunc("/", authIndex)
		http.HandleFunc(fmt.Sprintf("/%s", randomStr), updateIndex)
		http.HandleFunc("/check_auth", auth_Check)

	} else {
		http.HandleFunc(fmt.Sprintf("/"), updateIndex)
	}

	fmt.Printf("Update link: http://%s:%s/\n", publicIP, *port)
	fmt.Printf("Update link: http://%s:%s/\n\n", privateIP, *port)
}

// 文件上传页面展示
func updateIndex(w http.ResponseWriter, r *http.Request) {
	tpl := `
    <html>
    <head>
        <title>upload file</title>
    </head>
       <style>
           .update {
               position: relative;
               display: inline-block;
               background: #D0EEFF;
               border: 1px solid #99D3F5;
               border-radius: 4px;
               padding: 4px 12px;
               overflow: hidden;
               color: #1E88C7;
               text-decoration: none;
               text-indent: 0;
               line-height: 20px;
           }
           .update input {
               position: absolute;
               font-size: 100px;
               right: 0;
               top: 0;
               opacity: 0;
           }
          .update:hover {
               background: #AADFFD;
               border-color: #78C3F3;
               color: #004974;
               text-decoration: none;
          }
       </style>

        <body>
             <form id="uploadForm" method="POST" enctype="multipart/form-data" action="/{{ . }}">
             <input type="FILE" id="file" name="file" class="update"/>
             <input type="SUBMIT" value="upload"  class="update">
           </form>
        </body>
    </html>
   `
	log.Printf(
		"%s  %s  %s",
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
	)
	t, _ := template.New("fileupdate").Parse(tpl)
	t.Execute(w, randomstrup)
}

// 文件下载
func fileDown() {
	// 如果未设置密码就直接生成下载链接
	if *pwd == "nil" {
		if *file != "false" {
			file_Server()
		} else {
			directory_Server()
		}

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
	}

	if *file != "false" {
		file_Server()
	} else {
		directory_Server()
	}

	http.HandleFunc("/", authIndex)
	http.HandleFunc("/check_auth", auth_Check)
	fmt.Printf("Please access %s at the browser!\n\n", "0.0.0.0:"+*port)
}

// 密码验证页面
func authIndex(w http.ResponseWriter, r *http.Request) {
	tpl := `
		<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>file download</title>
			</head>

            <style>
            .file {
                color:#333;
                line-height:normal;
                font-family:"Microsoft YaHei",Tahoma,Verdana,SimSun;
                font-style:normal;
                font-variant:normal;
                font-size-adjust:none;
                font-stretch:normal;
                font-weight:normal;
                margin:auto;
                padding-left:700px;
                font-size:15px;
                outline-width:medium;
                outline-style:none;
                outline-color:invert;
                border-top-left-radius:3px;
                border-top-right-radius:3px;
                border-bottom-left-radius:3px;
                border-bottom-right-radius:3px;
                text-shadow:0px 1px 2px #fff;
                background-attachment:scroll;
                background-repeat:repeat-x;
                background-position-x:left;
                background-position-y:top;
                background-size:auto;
                background-origin:padding-box;
                background-clip:border-box;
                background-color:rgb(255,255,255);
                border-top-color:#ccc;
                border-right-color:#ccc;
                border-bottom-color:#ccc;
                border-left-color:#ccc;
                border-top-width:1px;
                border-right-width:1px;
                border-bottom-width:1px;
                border-left-width:1px;
                border-top-style:solid;
                border-right-style:solid;
                border-bottom-style:solid;
                border-left-style:solid;
            }
            .file:focus {
                border: 1px solid #fafafa;
                -webkit-box-shadow: 0px 0px 6px #007eff;
                -moz-box-shadow: 0px 0px 5px #007eff;
                box-shadow: 0px 0px 5px #007eff;   
    
            }
            </style>

			<body>

				<form  action="/check_auth" method='post' class="file">
					<label>密码:</label><input type=password id="password" name="password"/>
					<input type="submit" id="sub" value="提交">
				</form>
			</body>
			<script>
				$('#sub').on('click',function(){
					var password = $('#password').val();
                    $.ajax({
						type: 'post',
						url: '/auth_check',
						contentType: 'application/x-www-form-urlencoded',
						data: {"password": password},
					});
   			 })
			</script>
		</html>
	`
	log.Printf(
		"%s  %s  %s",
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
	)

	t, _ := template.New("filedown").Parse(tpl)
	t.Execute(w, nil)
}

// 校验密码
func auth_Check(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	if password != *pwd {
		w.Write([]byte("password error!"))
		return
	}

	if *update != false {
		http.Redirect(w, r, fmt.Sprintf("/%s", randomstr), 302)
	}

	if *file != "false" {
		http.Redirect(w, r, fmt.Sprintf("/%s/%s", randomstr, filename), 302)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/%s/", randomstr), 302)
	}

}

// 生成随机字符串
func randomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 文件下载
func file_Server() {
	ok, _ := pathExist(*file)
	if !ok {
		log.Fatal("File does not exist")
	}

	filename = filepath.Base(*file)

	if *pwd == "nil" {
		randomstr = "down"
	}

	fmt.Printf("Download link: http://%s:%s/%s/%s\n", publicIP, *port, randomstr, filename)
	fmt.Printf("Download link: http://%s:%s/%s/%s\n\n", privateIP, *port, randomstr, filename)

	// 打印访问日志
	http.HandleFunc(fmt.Sprintf("/%s/%s", randomstr, filename), func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"%s  %s  %s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
		)
		http.ServeFile(w, r, *file)
	})
}

// 目录游览
func directory_Server() {
	ok, _ := pathExist(*directory)
	if !ok {
		log.Fatal("Directory does not exist")
	}

	if *pwd == "nil" {
		randomstr = "down"
	} else {
		randomstr = randomStr(7)
	}

	fmt.Printf("Access link: http://%s:%s/%s/\n", publicIP, *port, randomstr)
	fmt.Printf("Access link: http://%s:%s/%s/\n\n", privateIP, *port, randomstr)

	http.Handle(fmt.Sprintf("/%s/", randomstr), func(prefix string, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf(
				"%s  %s  %s",
				r.RemoteAddr,
				r.Method,
				r.RequestURI,
			)
			if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
				r2 := new(http.Request)
				*r2 = *r
				r2.URL = new(url.URL)
				*r2.URL = *r.URL
				r2.URL.Path = p
				h.ServeHTTP(w, r2)
			} else {
				http.NotFound(w, r)
			}
		})
	}(fmt.Sprintf("/%s/", randomstr), http.FileServer(http.Dir(*directory))))
}

// 判断文件或目录是否存在
func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 获取公网IP
func getPublic() string {
	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
		return "0.0.0.0"
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "0.0.0.0"
	}

	return strings.Replace(string(b), "\n", "", -1)
}

// 获取本地IP
func getPrivate() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}

/*

easy -download -f demo.txt -p 80  //生成文件下载链接
easy -download -f demo.txt -p 80 -pwd 123456 //生成文件下载链接,需要密码校验
easy -download -d /data -p 80  //生成目录下载链接
easy -download -d /data -p 80 -pwd 123456 //生成目录下载链接,需要密码校验
easy -update -d /data -p 80  //生成上传链接，文件上传后存在/data目录下
easy -update -d /data  -p 80 -pwd 123456 //生成上传链接,需要密码校验

*/
