package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// Hàm kết nối đến cơ sở dữ liệu Mysql trong xampp
func dbConn() (*sql.DB, error) {
	// dataname:test
	// port:127.0.0.1:3379 - port mặc định là 3306
	// password:""
	// user:root
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3379)/test") // Không có password
	if err != nil {
		return nil, fmt.Errorf("không kết nối được tới cơ sở dữ liệu: %v", err)
	}
	return db, nil
}

// biến toàn cục để kết nối database
var conn, err = dbConn()

// Hàm insert dữ liệu vào cơ sở dữ liệu
func insertdb(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

// Hàm để thêm CORS headers-hàm này được tạo ra để nó chấp nhân mọi đầu vào theo cái chính sách cros gì gì đó của html
// nếu xóa hàm này đi thì XIN CHÚC MỪNG fetch sẽ không lấy được dữ liệu
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// hàm này được dùng để hiển thị dữ liệu
func showdb(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	//sau khi enablecors thì xóa dòng code bên dưới đi được
	w.Header().Set("Content-Type", "text/paint")
	//hàm này đồng thời cũng quyết định vị trí sẽ xuất hiện của các bảng
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
		html += fmt.Sprintf(`<tr><td><a href="edit.html?id=%d">%d</a></td><td>%s</td><td>%s</td></tr>`, id, id, name, email)
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

// tìm kiếm theo tên của bản
// có thể sử thành tìm kiếm tùy ý nếu muốn
func findbyName(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	// Lấy giá trị của biến name từ query string
	//biến name được lấy từ request của trang ví dụ "http//localhost:8000/findbyname?name="hi hi" thì biến name ở đây bằng "hi hi"
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Thiếu tên trong request", http.StatusBadRequest)
		return
	}
	// truy vấn dữ liệu theo tên
	rowss, err := conn.Query("select id,name,email from users where name like ?", "%"+name+"%")
	if err != nil {
		fmt.Fprint(w, "truy vấn bị lỗi ")
		fmt.Print("truy vấn tìm tên bị lỗi")
	}
	defer rowss.Close()
	var requestText string = "<table border='1'><tr><th>ID</th><th>Name</th><th>Email</th></tr>"
	for rowss.Next() {
		var id int
		var name, email string
		// Quét các cột trong hàng hiện tại vào biến
		if err := rowss.Scan(&id, &name, &email); err != nil {
			log.Printf("Lỗi khi đọc dữ liệu từ hàng: %v", err)
			http.Error(w, "Lỗi khi đọc dữ liệu", http.StatusInternalServerError)
			return
		}
		requestText += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", id, name, email)
	}
	requestText += "</table>"
	fmt.Fprint(w, requestText)
}

func editUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	// Lấy ID của người dùng từ query string
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Thiếu ID người dùng", http.StatusBadRequest)
		return
	}
	// Truy vấn lấy thông tin người dùng theo ID
	var name, email string
	err = conn.QueryRow("SELECT name, email FROM users WHERE id = ?", id).Scan(&name, &email)
	if err != nil {
		http.Error(w, "Không tìm thấy người dùng", http.StatusNotFound)
		return
	}

	// Hiển thị form HTML với các giá trị đã có trong database
	html := `
			<form action="http://localhost:8000/updateUser" method="POST">
				<input type="hidden" name="id" value="` + id + `">
				<label>Tên:</label><br>
				<input type="text" name="name" value="` + name + `"><br>
				<label>Email:</label><br>
				<input type="email" name="email" value="` + email + `"><br>
				<button type="submit">Cập nhật</button>
			</form>`
	// Trả về trang HTML
	w.Header().Set("Content-Type", "text/paint")
	fmt.Fprint(w, html)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != http.MethodPost {
		http.Error(w, "Phương thức không hợp lệ", http.StatusMethodNotAllowed)
		return
	}

	// Lấy dữ liệu từ form
	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")

	// Kiểm tra các giá trị có hợp lệ không
	if id == "" && name == "" && email == "" {
		http.Error(w, "Thiếu thông tin", http.StatusBadRequest)
		return
	}
	// Thực hiện câu truy vấn UPDATE
	_, err = conn.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", name, email, id)
	if err != nil {
		http.Error(w, "Cập nhật dữ liệu thất bại", http.StatusInternalServerError)
		return
	}

	// Thông báo cập nhật thành công
	fmt.Fprintf(w, "Cập nhật thành công người dùng có ID: %s", id)
}

// Hàm main
func main() {

	defer conn.Close()
	http.HandleFunc("/insert", insertdb)
	http.HandleFunc("/show", showdb)
	http.HandleFunc("/options", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
	})
	http.HandleFunc("/findbyname", findbyName)
	http.HandleFunc("/edit", editUser)
	http.HandleFunc("/updateUser", updateUser)
	log.Println("Server đang chạy tại http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
