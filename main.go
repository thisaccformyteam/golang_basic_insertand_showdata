// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// // func dbConn() (db *sql.DB) {
// // 	dns := "root:@tcp(127.0.0.1:3379)/test"
// // 	db, err := sql.Open("mysql", dns)
// // 	if err != nil {
// // 		log.Fatal("không kết nối được!", err)
// // 	}
// // 	return db
// // }
// func dbConn() (*sql.DB, error) {
// 	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3379)/test") // Không có password
// 	if err != nil {
// 		return nil, fmt.Errorf("không kết nối được tới cơ sở dữ liệu: %v", err)
// 	}
// 	return db, nil
// }

//	func insertdb(w http.ResponseWriter, r *http.Request) {
//		conn := dbConn()
//		defer conn.Close() // hàm này sẽ thực hiện vào cuối trương trình
//		name := r.FormValue("name")
//		email := r.FormValue("email")
//		if name == "" || email == "" {
//			http.Error(w, "Missing 'name' or 'email' parameter", http.StatusBadRequest)
//			return
//		}
//		_, err := conn.Exec("INSERT INTO users(name, email) VALUES(?, ?)", name, email)
//		if err != nil {
//			fmt.Print("insert false!!")
//		}
//	}
//
//	func main() {
//		http.HandleFunc("/insert", insertdb)
//		http.ListenAndServe(":8000", nil)
//	}
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// Hàm kết nối đến cơ sở dữ liệu
func dbConn() (*sql.DB, error) {
	// dataname:test
	// port:127.0.0.1:3379
	// password:""
	// user:root
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3379)/test") // Không có password
	if err != nil {
		return nil, fmt.Errorf("không kết nối được tới cơ sở dữ liệu: %v", err)
	}
	return db, nil
}

// Hàm insert dữ liệu vào cơ sở dữ liệu
func insertdb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := dbConn()
	if err != nil {
		log.Println("Kết nối cơ sở dữ liệu thất bại:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	name := r.FormValue("name")
	email := r.FormValue("email")
	if name == "" || email == "" {
		http.Error(w, "Missing 'name' or 'email' parameter", http.StatusBadRequest)
		return
	}

	_, err = conn.Exec("INSERT INTO users(name, email) VALUES(?, ?)", name, email)
	if err != nil {
		log.Println("Lỗi khi insert:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Dữ liệu đã được thêm thành công!")
}

// hàm hiển thị dữ liệu
// func showdb(w http.ResponseWriter, r *http.Request) {
// 	conn, err := dbConn()
// 	w.Header().Set("Content-Type", "text/html")
// 	defer conn.Close() //đóng kết nối khi thực hiện xong
// 	if err != nil {
// 		log.Println("ket noi database that bai %v", err)
// 	}
// 	rows, err := conn.Query("SELECT id,name,email FROM users")
// 	if err != nil {
// 		log.Println("some thing when wrong")
// 	}
// 	defer rows.Close() // kết thúc lệnh khi hiện thị xong
// 	// Tạo chuỗi HTML
// 	html := "<table border='1'><tr><th>ID</th><th>Name</th><th>Email</th></tr>"
// 	// Duyệt qua các kết quả
// 	for rows.Next() {
// 		var id int
// 		var name, email string

// 		// Quét các cột trong hàng hiện tại vào biến
// 		if err := rows.Scan(&id, &name, &email); err != nil {
// 			fmt.Errorf("Lỗi khi đọc dữ liệu từ hàng: %v", err)
// 		}

// 		// Tạo hàng HTML từ dữ liệu
// 		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", id, name, email)
// 	}

// 	// Kết thúc bảng
// 	html += "</table>"
// 	if rows.Err() != nil {
// 		fmt.Errorf("Nỗi khi duyệt kết quả")
// 	}
// 	fmt.Fprint(w, html)
// }

func showdb(w http.ResponseWriter, r *http.Request) {
	conn, err := dbConn()
	w.Header().Set("Content-Type", "text/html")
	defer conn.Close() // Đóng kết nối khi thực hiện xong

	if err != nil {
		log.Printf("Kết nối database thất bại: %v", err)
		http.Error(w, "Lỗi kết nối database", http.StatusInternalServerError)
		return
	}

	rows, err := conn.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Println("Lỗi truy vấn dữ liệu")
		http.Error(w, "Lỗi truy vấn dữ liệu", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Kết thúc lệnh khi hiển thị xong

	// Tạo chuỗi HTML
	html := "<table border='1'><tr><th>ID</th><th>Name</th><th>Email</th></tr>"

	// Duyệt qua các kết quả
	for rows.Next() {
		var id int
		var name, email string

		// Quét các cột trong hàng hiện tại vào biến
		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Printf("Lỗi khi đọc dữ liệu từ hàng: %v", err)
			http.Error(w, "Lỗi khi đọc dữ liệu", http.StatusInternalServerError)
			return
		}

		// Tạo hàng HTML từ dữ liệu
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", id, name, email)
	}

	// Kết thúc bảng
	html += "</table>"

	// Kiểm tra lỗi trong quá trình duyệt qua các hàng
	if rows.Err() != nil {
		log.Printf("Lỗi khi duyệt kết quả: %v", rows.Err())
		http.Error(w, "Lỗi duyệt dữ liệu", http.StatusInternalServerError)
		return
	}

	// Trả về HTML cho client
	fmt.Println(html)
	fmt.Fprint(w, html) // Trả chuỗi HTML về cho frontend
}

// Hàm main
func main() {
	http.HandleFunc("/insert", insertdb)
	http.HandleFunc("/show", showdb)
	log.Println("Server đang chạy tại http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
