package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

func signup(w http.ResponseWriter, r *http.Request) {
	// GETメソッド
	if r.Method == "GET" {
		// セッションの確認
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			// ログインしている場合
			http.Redirect(w, r, "/todos", 302)
		}
		// POSTメソッド
	} else if r.Method == "POST" {
		// ParseFormは入力フォームの解析
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		// ユーザの作成
		if err := user.CreateUser(); err != nil {
			log.Fatalln(err)
		}
		http.Redirect(w, r, "/", 302)

	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// セッションの確認
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		// ログインしている場合
		http.Redirect(w, r, "/todos", 302)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	// 入力されたメールアドレスからユーザを検索
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}
	// パスワードが正しいか検証
	if user.Password == models.Encrypt(r.PostFormValue("password")) {
		// パスワードが正しければセッションを作成する
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}
	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSession()
		if err != nil {
			log.Fatalln(err)
		}
	}
	http.Redirect(w,r,"/login", 302)
}
