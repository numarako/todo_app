package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"todo_app/app/models"
	"todo_app/config"
)

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}
	// mustはtemplateをキャッシュ
	// ParseFilesでエラーの場合はpanicを起こすというエラーハンドリングの実装
	templates := template.Must(template.ParseFiles(files...))
	// defineしたレイアウトを読み込むために明示的にlayoutを読み込む必要あり。
	templates.ExecuteTemplate(w, "layout", data)
}

func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	// == に注意
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid session")
		}
	}
	return sess, err
}

// URLの正規表現を変数として使用
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

// 関数を引数にHandlerFunc(関数)を返す関数
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /todos/edit/1
		// URLのうち、validPathとmatchした部分をスライスとして取得
		// "edit”の部分
		q := validPath.FindStringSubmatch(r.URL.Path)
		if q == nil {
			http.NotFound(w, r)
			return
		}
		// "1”の部分 qiのiはint
		qi, err := strconv.Atoi(q[2])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, qi)
	}
}

func StartMainServer() error {
	// 静的ファイルの読み込み
	files := http.FileServer(http.Dir(config.Config.Static))
	// URLパスの/statc/を登録するが、/statc/ディレクトリがないから、/static/を取り除く。
	http.Handle("/static/", http.StripPrefix("/static/", files))
	// ハンドラーを登録(URL)→top関数へ
	http.HandleFunc("/", top)
	// ハンドラーはルーティングのようなもの？
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/save", todoSave)
	// /で終わらせることでURLを部分一致として扱うことができる
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))
	http.HandleFunc("/todos/update/", parseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", parseURL(todoDelete))
	return http.ListenAndServe(":"+config.Config.Port, nil)
}
