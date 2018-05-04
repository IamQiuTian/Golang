package main

import (
	"flag"
	"fmt"
	"html/template"
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
	directory *string = flag.String("d", "false", "Directory path")
	file      *string = flag.String("f", "false", "file path")
	port      *string = flag.String("p", "8888", "Listening port")
	auth      *string = flag.String("auth", "nil", "password")
)

var (
	randomstr string
	publicIP  string
	privateIP string
)

func main() {
	flag.Parse()
	// 参数为空或参数都不为空
	if *directory == "false" && *file == "false" || *directory != "false" && *file != "false" {
		flag.Usage()
		return
	}
	// 如果未设置密码就直接生成下载链接
	if *auth == "nil" {
		if *file != "false" {
			file_Server()
		} else {
			directory_Server()
		}
		fmt.Printf("%s", publicIP)
		fmt.Println(privateIP)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
	}

	if *file != "false" {
		file_Server()
	} else {
		directory_Server()
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/check_auth", auth_Check)
	fmt.Printf("Please access %s port at the browser!\n", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))

}

// 登录页
func index(w http.ResponseWriter, r *http.Request) {
	tpl := `
		<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>easydown</title>
			</head>
			<body>
				<form  action="/check_auth" method='post'>
					<label>密码:</label><input type=password id="password" name="password" />
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

	t, _ := template.New("easydown").Parse(tpl)
	t.Execute(w, nil)
}

// 校验密码
func auth_Check(w http.ResponseWriter, r *http.Request) {
	pwd := r.FormValue("password")
	if pwd != *auth {
		w.Write([]byte("password error!"))
		return
	}
	w.Write([]byte(fmt.Sprintf("%s\n%s", publicIP, privateIP)))
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
		fmt.Println("File does not exist")
		return
	}

	// 获取文件路径或IP地址
	randomstr = randomStr(7)
	filename := filepath.Base(*file)

	getPublic(*port, filename)
	getPrivate(*port, filename)

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
		fmt.Println("Directory does not exist")
		return
	}
	filename := "nil"
	randomstr = randomStr(7)

	getPublic(*port, filename)
	getPrivate(*port, filename)

	http.Handle(fmt.Sprintf("/%s/", randomstr), func(prefix string, h http.Handler) http.Handler {
		if prefix == "" {
			return h
		}
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
func getPublic(port, filename string) {
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if filename == "nil" {
		publicIP = fmt.Sprintf("Access link: http://%s:%s/%s/\n", strings.Replace(string(b), "\n", "", -1), port, randomstr)
	} else {
		publicIP = fmt.Sprintf("Download link: http://%s:%s/%s/%s\n", strings.Replace(string(b), "\n", "", -1), port, randomstr, filename)
	}
}

// 获取本地IP
func getPrivate(port, filename string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if filename == "nil" {
					privateIP = fmt.Sprintf("Access link: http://%s:%s/%s/\n", ipnet.IP.String(), port, randomstr)
				} else {
					privateIP = fmt.Sprintf("Download link: http://%s:%s/%s/%s\n", ipnet.IP.String(), port, randomstr, filename)
				}
			}
		}
	}
}



/*
Used:
$ easyDown -f FilePath -p 8080
$ easyDown -d DirectoryPath -p 8080
$ easyDown -d DirectoryPath -p 8080 -auth dkhasnyqw
*/
