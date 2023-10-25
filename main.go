package main

import (
	"bufio"
	"database/sql"
	"enigma-laundry/entity"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "enigma_laundry"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

var scanner = bufio.NewScanner(os.Stdin)

func main() {
	db := connectDB()
	defer db.Close()

	fmt.Println()

	for {
		fmt.Println()
		fmt.Println("Main Menu:")
		fmt.Println("1. Data Customer")
		fmt.Println("2. Data Karyawan")
		fmt.Println("3. Data Layanan")
		fmt.Println("4. Data Transaksi")
		fmt.Println("5. Keluar")
		fmt.Println()

		var choice int
		fmt.Print("Pilih menu: ")
		// validasi main pilihan tidak valid
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
			continue
		}

		switch choice {
		case 1:
			customerMenu(db)
		case 2:
			karyawanMenu(db)
		case 3:
			layananMenu(db)
		case 4:
			transaksiMenu(db)
		case 5:
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
		}
	}

}

func connectDB() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// menu master customer
func customerMenu(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 80))
	fmt.Println()
	fmt.Println("Customer Menu")
	fmt.Println("1. VIEW Data Customer")
	fmt.Println("2. INSERT Data Customer")
	fmt.Println("3. UPDATE Data Customer")
	fmt.Println("4. DELETE Data Customer")
	fmt.Println("5. Kembali ke Menu Utama")
	fmt.Println()

	var choice int

	fmt.Print("Pilih menu: ")
	_, err := fmt.Scan(&choice)
	// validasi 1 master customer pilihan tidak valid
	if err != nil {
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")

	}

	switch choice {
	case 1:
		viewDataCust(db)
	case 2:
		insertDataCust(db)
	case 3:
		updateDataCust(db)
	case 4:
		deleteDataCust(db)
	case 5:
		return
	default:
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
	}
}

// view customer
func viewDataCust(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 80))
	fmt.Println()

	fmt.Println("VIEW Data Customer")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-5s | %-20s | %-30s | %-20s\n", "ID", "Nama Customer", "Alamat", "Telepon")
	fmt.Println(strings.Repeat("-", 80))

	// ambil func
	customers := viewCustomer()
	for _, customer := range customers {
		fmt.Printf("%-5s | %-20s | %-30s | %-20s\n", customer.Id, customer.Name, customer.Address, customer.Phone)
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
	customerMenu(db)
}

func viewCustomer() []entity.Customers {
	db := connectDB()
	defer db.Close()

	sqlStatement := "SELECT * FROM customers ORDER BY id ASC;"

	rows, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cust := scanCustomer(rows)

	return cust
}

func scanCustomer(rows *sql.Rows) []entity.Customers {
	customers := []entity.Customers{}
	var err error

	for rows.Next() {
		customer := entity.Customers{}
		err = rows.Scan(&customer.Id, &customer.Name, &customer.Address, &customer.Phone)

		if err != nil {
			panic(err)
		}

		customers = append(customers, customer)

	}

	err = rows.Err()
	if err != nil {
		panic(err)

	}

	return customers
}

// insert customer
func insertDataCust(db *sql.DB) {

	fmt.Println("INSERT Data Customer")
	var name string
	var address string
	var phone string
	scanner.Scan()

	fmt.Print("Masukkan nama: ")
	scanner.Scan()
	name = scanner.Text()

	// Validasi 2 master CUSTOMER: Nama harus diisi
	if name == "" {
		fmt.Println("Nama harus diisi.")
		return
	}

	fmt.Print("Masukkan alamat: ")
	scanner.Scan()
	address = scanner.Text()

	fmt.Print("Masukkan nomor telepon: ")
	scanner.Scan()
	phone = scanner.Text()

	// Mengambil ID terakhir
	var lastID string
	err := db.QueryRow("SELECT id FROM customers ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1").Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	// Menghitung ID berikutnya
	nextID := "C1"
	if lastID != "" {
		// Ekstrak angka dari ID terakhir
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			panic(err)
		}
		// Tambahkan 1 ke angka tersebut
		nextNumber := lastNumber + 1
		// Format ID berikutnya
		nextID = fmt.Sprintf("C%d", nextNumber)
	}

	customer := entity.Customers{Id: nextID, Name: name, Address: address, Phone: phone}
	addCustomer(customer)
	fmt.Println()

	fmt.Printf("Data dengan ID %s berhasil dimasukkan.\n", nextID)
	fmt.Println()
	customerMenu(db)
}

func addCustomer(customer entity.Customers) {
	db := connectDB()
	defer db.Close()

	var err error

	sqlStatement := "INSERT INTO customers (id, name, address, phone) VALUES ($1, $2, $3, $4);"

	_, err = db.Exec(sqlStatement, customer.Id, customer.Name, customer.Address, customer.Phone)

	if err != nil {
		panic(err)
	}
}

// update customer
func updateDataCust(db *sql.DB) {

	fmt.Println("Update Data Customer")
	var id string
	var name string
	var address string
	var phone string

	fmt.Print("Masukkan id yang ingin di ubah: ")

	_, err := fmt.Scan(&id)
	id = strings.ToUpper(id)
	if err != nil {
		panic(err)
	}

	// Validasi 3 MASTER CUSTOMER: id harus ada dalam database
	var existingId string
	err = db.QueryRow("SELECT * FROM customers WHERE id = $1", id).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data dengan ID %s yang diberikan tidak ditemukan.\n", id)
			fmt.Println()
			customerMenu(db)
			return
		}
	}
	scanner.Scan()

	fmt.Print("Masukkan nama: ")
	scanner.Scan()
	name = scanner.Text()

	// Validasi: Nama harus diisi
	if name == "" {
		fmt.Println("Nama harus diisi.")
		return
	}

	fmt.Print("Masukkan alamat: ")
	scanner.Scan()
	address = scanner.Text()

	fmt.Print("Masukkan nomor telepon: ")
	scanner.Scan()
	phone = scanner.Text()

	customer := entity.Customers{Id: id, Name: name, Address: address, Phone: phone}
	updateCustomer(customer)
	fmt.Println()
	fmt.Printf("Data dengan ID %s berhasil diubah.\n", id)
	fmt.Println()
	customerMenu(db)
}

func updateCustomer(customer entity.Customers) {
	db := connectDB()
	defer db.Close()
	var err error

	sqlStatement := "UPDATE customers SET name = $2, address = $3, phone = $4 WHERE id = $1;"

	_, err = db.Exec(sqlStatement, customer.Id, customer.Name, customer.Address, customer.Phone)
	if err != nil {
		panic(err)
	}
}

// delete customer
func deleteDataCust(db *sql.DB) {
	fmt.Println("DELETE Data Customer")
	var id string

	fmt.Print("Masukkan ID data yang akan dihapus: ")
	_, err := fmt.Scan(&id)
	id = strings.ToUpper(id)
	if err != nil {
		panic(err)
	}

	// Validasi: id harus ada dalam database
	var existingId string
	err = db.QueryRow("SELECT * FROM customers WHERE id = $1", id).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data dengan ID %s yang diberikan tidak ditemukan.\n", id)
			fmt.Println()
			customerMenu(db)
			return
		}
	}

	// Validasi 4 MASTER KARYAWAN: Cek apakah ada referensi dari tabel lain ke karyawan yang akan dihapus
	var hasReferences bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM transaksi WHERE customer_id = $1);", id).Scan(&hasReferences)

	if err != nil {
		panic(err)
	}

	if hasReferences {
		fmt.Println()
		fmt.Printf("Data customer dengan ID %s tidak dapat dihapus karena memiliki referensi dari tabel lain.\n", id)
		fmt.Println()
		customerMenu(db)
		return
	}

	deleteCustomer(id)
	fmt.Println()
	fmt.Printf("Data dengan ID %s berhasil dihapus.\n", id)
	fmt.Println()
	customerMenu(db)
}

func deleteCustomer(id string) {
	db := connectDB()
	defer db.Close()
	var err error

	sqlStatement := "DELETE FROM customers WHERE id = $1"

	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
}

// menu master karyawan
func karyawanMenu(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 80))
	fmt.Println()
	fmt.Println("Karyawan Menu")
	fmt.Println("1. VIEW Data Karyawan")
	fmt.Println("2. INSERT Data Karyawan")
	fmt.Println("3. UPDATE Data Karyawan")
	fmt.Println("4. DELETE Data Karyawan")
	fmt.Println("5. Kembali ke Menu Utama")
	fmt.Println()

	var choice int
	fmt.Print("Pilih menu: ")
	_, err := fmt.Scan(&choice)
	// validasi 1 master karyawan : pilihan tidak valid
	if err != nil {
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
	}

	switch choice {
	case 1:
		viewDataKaryawan(db)
	case 2:
		insertDataKaryawan(db)
	case 3:
		updateDataKaryawan(db)
	case 4:
		deleteDataKaryawan(db)
	case 5:
		return
	default:
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
	}
}

// view karyawan
func viewDataKaryawan(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 80))
	fmt.Println()
	fmt.Println("VIEW Data Karyawan")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-5s | %-20s | %-30s | %-20s\n", "ID", "Nama Karyawan", "Alamat", "Telepon")
	fmt.Println(strings.Repeat("-", 80))

	// ambil func
	karyawans := viewKaryawan()
	for _, karyawan := range karyawans {
		fmt.Printf("%-5s | %-20s | %-30s | %-20s\n", karyawan.Id, karyawan.Name, karyawan.Address, karyawan.Phone)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
	karyawanMenu(db)
}

func viewKaryawan() []entity.Karyawan {
	db := connectDB()
	defer db.Close()
	sqlStatement := "SELECT * FROM karyawan ORDER BY id ASC;"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	karyawan := scanKaryawan(rows)
	return karyawan
}

func scanKaryawan(rows *sql.Rows) []entity.Karyawan {
	karyawans := []entity.Karyawan{}
	var err error
	for rows.Next() {
		karyawan := entity.Karyawan{}
		err = rows.Scan(&karyawan.Id, &karyawan.Name, &karyawan.Address, &karyawan.Phone)

		if err != nil {
			panic(err)
		}
		karyawans = append(karyawans, karyawan)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return karyawans
}

// insert karyawan
func insertDataKaryawan(db *sql.DB) {

	fmt.Println("INSERT Data Karyawan")
	var name string
	var address string
	var phone string
	scanner.Scan()

	fmt.Print("Masukkan nama: ")
	scanner.Scan()
	name = scanner.Text()

	// Validasi 2 MASTER KARYAWAN : Nama harus diisi
	if name == "" {
		fmt.Println("Nama harus diisi.")
		return
	}

	fmt.Print("Masukkan alamat: ")
	scanner.Scan()
	address = scanner.Text()

	fmt.Print("Masukkan nomor telepon: ")
	scanner.Scan()
	phone = scanner.Text()

	// Mengambil ID terakhir
	var lastID string
	err := db.QueryRow("SELECT id FROM karyawan ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1").Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	// Menghitung ID berikutnya
	nextID := "K1"
	if lastID != "" {
		// Ekstrak angka dari ID terakhir
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			panic(err)
		}
		// Tambahkan 1 ke angka tersebut
		nextNumber := lastNumber + 1
		// Format ID berikutnya
		nextID = fmt.Sprintf("K%d", nextNumber)
	}

	karyawan := entity.Karyawan{Id: nextID, Name: name, Address: address, Phone: phone}
	addKaryawan(karyawan)
	fmt.Println()
	fmt.Printf("Data dengan ID %s berhasil dimasukkan.\n", nextID)
	fmt.Println()
	karyawanMenu(db)
}

func addKaryawan(karyawan entity.Karyawan) {
	db := connectDB()
	defer db.Close()
	var err error
	sqlStatement := "INSERT INTO karyawan (id, name, address, phone) VALUES ($1, $2, $3, $4);"
	_, err = db.Exec(sqlStatement, karyawan.Id, karyawan.Name, karyawan.Address, karyawan.Phone)
	if err != nil {
		panic(err)
	}
}

// update karyawan
func updateDataKaryawan(db *sql.DB) {

	fmt.Println("Update Data Karyawan")

	var id string
	var name string
	var address string
	var phone string

	fmt.Print("Masukkan id yang ingin di ubah: ")
	_, err := fmt.Scan(&id)
	id = strings.ToUpper(id)
	if err != nil {
		panic(err)
	}
	// Validasi 3 MASTER KARYAWAN: id harus ada dalam database
	var existingId string
	err = db.QueryRow("SELECT * FROM karyawan WHERE id = $1", id).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data dengan ID %s yang diberikan tidak ditemukan.\n", id)
			fmt.Println()
			karyawanMenu(db)
			return
		}
	}

	scanner.Scan()
	fmt.Print("Masukkan nama: ")
	scanner.Scan()
	name = scanner.Text()

	// Validasi : Nama harus diisi
	if name == "" {
		fmt.Println("Nama harus diisi.")
		return
	}

	fmt.Print("Masukkan alamat: ")
	scanner.Scan()
	address = scanner.Text()

	fmt.Print("Masukkan nomor telepon: ")
	scanner.Scan()
	phone = scanner.Text()

	karyawan := entity.Karyawan{Id: id, Name: name, Address: address, Phone: phone}
	updateKaryawan(karyawan)
	fmt.Println()
	fmt.Printf("Data dengan ID %s berhasil diubah.\n", id)
	fmt.Println()
	karyawanMenu(db)
}

func updateKaryawan(karyawan entity.Karyawan) {
	db := connectDB()
	defer db.Close()
	var err error
	sqlStatement := "UPDATE karyawan SET name = $2, address = $3, phone = $4 WHERE id = $1;"
	_, err = db.Exec(sqlStatement, karyawan.Id, karyawan.Name, karyawan.Address, karyawan.Phone)
	if err != nil {
		panic(err)
	}
}

// delete karyawan
func deleteDataKaryawan(db *sql.DB) {
	fmt.Println("DELETE Data Karyawan")
	var id string
	fmt.Print("Masukkan ID data yang akan dihapus: ")
	_, err := fmt.Scan(&id)
	id = strings.ToUpper(id)
	if err != nil {
		panic(err)
	}
	// Validasi : id harus ada dalam database
	var existingId string
	err = db.QueryRow("SELECT * FROM karyawan WHERE id = $1", id).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data dengan ID %s yang diberikan tidak ditemukan.\n", id)
			fmt.Println()
			karyawanMenu(db)
			return
		}
	}

	// Validasi 4 MASTER KARYAWAN: Cek apakah ada referensi dari tabel lain ke karyawan yang akan dihapus
	var hasReferences bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM transaksi WHERE karyawan_id = $1);", id).Scan(&hasReferences)

	if err != nil {
		panic(err)
	}

	if hasReferences {
		fmt.Println()
		fmt.Printf("Data karyawan dengan ID %s tidak dapat dihapus karena memiliki referensi dari tabel lain.\n", id)
		fmt.Println()
		karyawanMenu(db)
		return
	}

	deleteKaryawan(id)
	fmt.Println()
	fmt.Printf("Data dengan ID %s berhasil dihapus.\n", id)
	fmt.Println()
	karyawanMenu(db)
}

func deleteKaryawan(id string) {
	db := connectDB()
	defer db.Close()
	var err error
	sqlStatement := "DELETE FROM karyawan WHERE id = $1"
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
}

// menu master layanan
func layananMenu(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 80))
	fmt.Println()
	fmt.Println("Layanan Menu")
	fmt.Println("1. VIEW Data Layanan")
	fmt.Println("2. INSERT Data Layanan")
	fmt.Println("3. UPDATE Data Layanan")
	fmt.Println("4. DELETE Data Layanan")
	fmt.Println("5. Kembali ke Menu Utama")
	fmt.Println()

	var choice int
	fmt.Print("Pilih menu: ")
	_, err := fmt.Scan(&choice)

	// validasi 1 master layanan pilihan tidak valid
	if err != nil {
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
	}

	switch choice {
	case 1:
		viewDataLayanan(db)
	case 2:
		insertDataLayanan(db)
	case 3:
		updateDataLayanan(db)
	case 4:
		deleteDataLayanan(db)
	case 5:
		return
	default:
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
	}
}

// view layanan
func viewDataLayanan(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 80))
	fmt.Println()
	fmt.Println("VIEW Data Layanan")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-5s | %-20s | %-30s | %-20s\n", "ID", "Nama Layanan", "Satuan", "Harga")
	fmt.Println(strings.Repeat("-", 80))

	// ambil func
	layanans := viewLayanan()
	for _, layanan := range layanans {
		fmt.Printf("%-5s | %-20s | %-30s | %-20.2f\n", layanan.Id, layanan.Name, layanan.Unit, layanan.Price)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
	layananMenu(db)
}

func viewLayanan() []entity.Layanan {
	db := connectDB()
	defer db.Close()
	sqlStatement := "SELECT * FROM layanan ORDER BY id ASC;"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	layanan := scanLayanan(rows)
	return layanan
}

func scanLayanan(rows *sql.Rows) []entity.Layanan {
	layanans := []entity.Layanan{}
	var err error
	for rows.Next() {
		layanan := entity.Layanan{}
		err = rows.Scan(&layanan.Id, &layanan.Name, &layanan.Unit, &layanan.Price)

		if err != nil {
			panic(err)
		}
		layanans = append(layanans, layanan)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return layanans
}

// insert layanan
func insertDataLayanan(db *sql.DB) {

	fmt.Println("INSERT Data Layanan")

	var name string
	var unit string
	var price float64
	var err error
	scanner.Scan()

	fmt.Print("Masukkan nama: ")
	scanner.Scan()
	name = scanner.Text()

	// Validasi 2 MASTER LAYANAN : Nama harus diisi
	if name == "" {
		fmt.Println("Nama harus diisi.")
		return
	}

	fmt.Print("Masukkan satuan: ")
	scanner.Scan()
	unit = scanner.Text()

	// Validasi : satuan harus diisi
	if unit == "" {
		fmt.Println("Satuan harus diisi.")
		return
	}

	fmt.Print("Masukkan harga satuan: ")
	scanner.Scan()
	price, err = strconv.ParseFloat(scanner.Text(), 64)

	// validasi 3 MASTER LAYANAN : input harus angka
	if err != nil {
		fmt.Println("Harga satuan tidak valid. Harap masukkan angka.")

	}

	// Validasi 4 MASTER LAYANAN : harga satuan harus diatas 0
	if price <= 0 {
		fmt.Println("harga satuan harus diatas 0")
		return
	}

	// Mengambil ID terakhir
	var lastID string
	err = db.QueryRow("SELECT id FROM layanan ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1").Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	// Menghitung ID berikutnya
	nextID := "L1"
	if lastID != "" {
		// Ekstrak angka dari ID terakhir
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			panic(err)
		}
		// Tambahkan 1 ke angka tersebut
		nextNumber := lastNumber + 1
		// Format ID berikutnya
		nextID = fmt.Sprintf("L%d", nextNumber)
	}

	layanan := entity.Layanan{Id: nextID, Name: name, Unit: unit, Price: price}
	addLayanan(layanan)
	fmt.Println()
	fmt.Printf("Data Layanan dengan ID %s berhasil dimasukkan.\n", nextID)
	fmt.Println()
	layananMenu(db)
}

func addLayanan(layanan entity.Layanan) {
	db := connectDB()
	defer db.Close()
	var err error
	sqlStatement := "INSERT INTO layanan (id, name, unit, price) VALUES ($1, $2, $3, $4);"
	_, err = db.Exec(sqlStatement, layanan.Id, layanan.Name, layanan.Unit, layanan.Price)
	if err != nil {
		panic(err)
	}
}

// update layanan
func updateDataLayanan(db *sql.DB) {

	fmt.Println("Update Data Layanan")
	var id string
	var name string
	var unit string
	var price float64

	fmt.Print("Masukkan id yang ingin di ubah: ")
	_, err := fmt.Scan(&id)
	id = strings.ToUpper(id)
	if err != nil {
		panic(err)
	}

	// Validasi 5 MASTER LAYANAN: id harus ada dalam database
	var existingId string
	err = db.QueryRow("SELECT * FROM layanan WHERE id = $1", id).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data Layanan dengan ID %s yang diberikan tidak ditemukan.\n", id)
			fmt.Println()
			layananMenu(db)
			return
		}
	}
	scanner.Scan()

	fmt.Print("Masukkan nama: ")
	scanner.Scan()
	name = scanner.Text()

	// Validasi : Nama harus diisi
	if name == "" {
		fmt.Println("Nama harus diisi.")
		return
	}

	fmt.Print("Masukkan satuan: ")
	scanner.Scan()
	unit = scanner.Text()

	// Validasi : satuan harus diisi
	if unit == "" {
		fmt.Println("Satuan harus diisi.")
		return
	}

	fmt.Print("Masukkan harga satuan: ")
	scanner.Scan()
	price, err = strconv.ParseFloat(scanner.Text(), 64)

	// validasi : input harus angka
	if err != nil {
		fmt.Println("Harga satuan tidak valid. Harap masukkan angka.")

	}

	// Validasi : harga satuan harus diatas 0
	if price <= 0 {
		fmt.Println("harga satuan harus diatas 0")
		return
	}

	layanan := entity.Layanan{Id: id, Name: name, Unit: unit, Price: price}
	updateLayanan(layanan)
	fmt.Println()
	fmt.Printf("Data Layanan dengan ID %s berhasil diubah.\n", id)
	fmt.Println()
	layananMenu(db)
}

func updateLayanan(layanan entity.Layanan) {
	db := connectDB()
	defer db.Close()
	var err error
	sqlStatement := "UPDATE layanan SET name = $2, unit = $3, price = $4 WHERE id = $1;"
	_, err = db.Exec(sqlStatement, layanan.Id, layanan.Name, layanan.Unit, layanan.Price)
	if err != nil {
		panic(err)
	}
}

// delete layanan
func deleteDataLayanan(db *sql.DB) {
	fmt.Println("DELETE Data Layanan")
	var id string
	fmt.Print("Masukkan ID Layanan yang akan dihapus: ")
	_, err := fmt.Scan(&id)
	id = strings.ToUpper(id)
	if err != nil {
		panic(err)
	}
	// Validasi: id harus ada dalam database
	var existingId string
	err = db.QueryRow("SELECT * FROM layanan WHERE id = $1", id).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data layanan dengan ID %s yang diberikan tidak ditemukan.\n", id)
			fmt.Println()
			layananMenu(db)
			return
		}
	}

	// Validasi 6 MASTER LAYANAN: Cek apakah ada referensi dari tabel lain ke layanan yang akan dihapus
	var hasReferences bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM detailTransaksi WHERE layanan_id = $1);", id).Scan(&hasReferences)
	if err != nil {
		panic(err)
	}

	if hasReferences {
		fmt.Println()
		fmt.Printf("Data layanan dengan ID %s tidak dapat dihapus karena memiliki referensi dari tabel lain.\n", id)
		fmt.Println()
		layananMenu(db)
		return
	}

	deleteLayanan(id)
	fmt.Println()
	fmt.Printf("Data layanan dengan ID %s berhasil dihapus.\n", id)
	fmt.Println()
	layananMenu(db)
}

func deleteLayanan(id string) {
	db := connectDB()
	defer db.Close()
	var err error
	sqlStatement := "DELETE FROM layanan WHERE id = $1"
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
}

// menu tx.transaksi

func transaksiMenu(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 80))
	fmt.Println()
	fmt.Println("Transaksi Menu")
	fmt.Println("1. VIEW Data Transaksi")
	fmt.Println("2. INSERT Data Transaksi")
	
	fmt.Println("3. Kembali ke Menu Utama")
	fmt.Println()

	var choice int
	fmt.Print("Pilih menu: ")
	_, err := fmt.Scan(&choice)

	// validasi 1 tx.transaksi pilihan tidak valid
	if err != nil {
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
	}

	switch choice {
	case 1:
		viewDataTransaksi(db)
	case 2:
		insertDataTransaksi(db)
	case 3:
		return
	default:
		fmt.Println("Pilihan tidak valid, Kembali ke menu utama.")
	}
}

// view transaksi
func viewDataTransaksi(db *sql.DB) {

	fmt.Println()
	fmt.Println(strings.Repeat("#", 120))
	fmt.Println()
	fmt.Println("VIEW Data Transaksi")
	fmt.Println(strings.Repeat("-", 120))
	fmt.Printf("%-5s | %-20s | %-20s | %-20s| %-20s| %-30s\n", "ID", "Nama Customer", "Nama Penerima", "Tanggal Masuk", "Tanggal Keluar", "Total")
	fmt.Println(strings.Repeat("-", 120))

	// ambil func
	transaksi := viewTransaksi()
	for _, trans := range transaksi {
		startTime := trans.StartDate.Format("02-01-2006")
		endTime := trans.EndDate.Format("02-01-2006")
		fmt.Printf("%-5s | %-20s | %-20s | %-20s| %-20s| %-30.2f\n", trans.Id, trans.Customer_id, trans.Karyawan_id, startTime, endTime, trans.Total)
	}

	fmt.Println(strings.Repeat("-", 120))
	fmt.Println()

	fmt.Println("Apakah Anda ingin melihat detail transaksi? (y/n): ")

	var choice string
	_, err := fmt.Scan(&choice)
	if err != nil {
		panic(err)
	}

	// validasi 2 tx.transaksi = ingin melihat detail transaksi
	if strings.ToLower(choice) == "y" {
		viewDetail(db)
	} else if strings.ToLower(choice) == "n" {
		transaksiMenu(db)
	} else {
		fmt.Println("Pilihan tidak valid. Kembali ke menu utama.")
		main()
	}
}

func viewTransaksi() []entity.Transaksi {
	db := connectDB()
	defer db.Close()
	sqlStatement := "SELECT t.id, c.name, k.name, t.startDate, t.endDate, t.total  FROM transaksi AS t JOIN customers AS c ON c.id = t.customer_id JOIN karyawan AS k ON k.id = t.karyawan_id ORDER BY id ASC;"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	transaksi := scanTransaksi(rows)
	return transaksi
}

func scanTransaksi(rows *sql.Rows) []entity.Transaksi {
	transaksies := []entity.Transaksi{}
	var err error
	for rows.Next() {
		transaksi := entity.Transaksi{}
		err = rows.Scan(&transaksi.Id, &transaksi.Customer_id, &transaksi.Karyawan_id, &transaksi.StartDate, &transaksi.EndDate, &transaksi.Total)

		if err != nil {
			panic(err)
		}
		transaksies = append(transaksies, transaksi)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return transaksies
}

// insert transaksi
func insertDataTransaksi(db *sql.DB) {

	fmt.Println("INSERT Data Transaksi")

	var customerId string
	var karyawanId string
	var startDate string
	var endDate string

	scanner.Scan()

	fmt.Print("Masukkan ID customer: ")
	scanner.Scan()
	customerId = scanner.Text()
	customerId = strings.ToUpper(customerId)

	// Validsi 3 tx.Transaksi: id customer harus diisi
	if customerId == "" {
		fmt.Println("ID customer harus diisi.")
		return
	}

	// Validasi 4 tx.transaksi: id harus ada dalam database
	var existingCust string
	err := db.QueryRow("SELECT name FROM customers WHERE id = $1", customerId).Scan(&existingCust)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data Customer dengan ID %s yang diberikan tidak ditemukan.\n", customerId)
			fmt.Println()
			transaksiMenu(db)
			return
		} 
	} else {
		fmt.Println("Nama Customer:",existingCust)
	}

	fmt.Print("Masukkan ID karyawan: ")
	scanner.Scan()
	karyawanId = scanner.Text()
	karyawanId = strings.ToUpper(karyawanId)

	//validasi
	if karyawanId == "" {
		fmt.Println("ID Karyawan harus diisi.")
		return
	}

	//validasi
	var existingKar string
	err = db.QueryRow("SELECT name FROM karyawan WHERE id = $1", karyawanId).Scan(&existingKar)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data Karyawan dengan ID %s yang diberikan tidak ditemukan.\n", karyawanId)
			fmt.Println()
			transaksiMenu(db)
			return
		} 
	} else {
		fmt.Println("Nama Karyawan:",existingKar)
	}

	fmt.Print("Masukan Tanggal Masuk (yyyy-mm-dd): ")
	fmt.Scanln(&startDate)

	fmt.Print("Masukan Tanggal Keluar (yyyy-mm-dd): ")
	fmt.Scanln(&endDate)

	startDateTime, err1 := time.Parse("2006-01-02", startDate)
	endDateTime, err2 := time.Parse("2006-01-02", endDate)

	// validasi 5 tx.transaksi = format harus yyyy-mm-dd
	if err1 != nil || err2 != nil {
		fmt.Println("Format tanggal salah. Pastikan menggunakan format yyyy-mm-dd.")
		return
	}

	// validasi 6 tx.transaksi = Tanggal Masuk tidak boleh mendahului tanggal Keluar
	if endDateTime.Before(startDateTime) {
		fmt.Println("Tanggal Masuk tidak boleh mendahului tanggal Keluar.")
		return
	}

	// Mengambil ID terakhir
	var lastID string
	err = db.QueryRow("SELECT id FROM transaksi ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1").Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	// Menghitung ID berikutnya
	nextID := "T1"
	if lastID != "" {
		// Ekstrak angka dari ID terakhir
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			panic(err)
		}
		// Tambahkan 1 ke angka tersebut
		nextNumber := lastNumber + 1
		// Format ID berikutnya
		nextID = fmt.Sprintf("T%d", nextNumber)
	}

	transaksi := entity.Transaksi{Id: nextID, Customer_id: customerId, Karyawan_id: karyawanId, StartDate: startDateTime, EndDate: endDateTime, Total: 0}
	addTransaksi(transaksi)
	fmt.Println()
	fmt.Printf("Data Transaksi dengan ID %s berhasil dimasukkan.\n", nextID)
	fmt.Println()
	insertDetail(db, nextID)
}

func addTransaksi(transaksi entity.Transaksi) {
	db := connectDB()
	defer db.Close()
	var err error
	sqlStatement := "INSERT INTO transaksi (id, customer_id, karyawan_id, startDate, endDate, total)VALUES ($1, $2, $3, $4, $5, $6);"
	_, err = db.Exec(sqlStatement, transaksi.Id, transaksi.Customer_id, transaksi.Karyawan_id, transaksi.StartDate, transaksi.EndDate, transaksi.Total)
	if err != nil {
		panic(err)
	}
}


// menu tx.detail transaksi

// view detail transaksi
func viewDetail(db *sql.DB) {
	var id string

	fmt.Println("VIEW DETAIL TRANSAKSI")
	fmt.Println("Masukan Id Transaksi: ")
	_, err := fmt.Scan(&id)
	id = strings.ToUpper(id)
	if err != nil {
		panic(err)
	}

	// Validasi 1 tx.detail = ID Transaksi id harus ada dalam database
	var existingId string
	err = db.QueryRow("SELECT * FROM transaksi WHERE id = $1", id).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data dengan ID %s yang diberikan salah.\n", id)
			fmt.Println()
			transaksiMenu(db)
			return
		}
	}

	// Panggil fungsi untuk menampilkan data detail transaksi berdasarkan ID
	viewDetailById(db, id)
}

func viewDetailById(db *sql.DB, transaksiID string) {
	fmt.Println()
	fmt.Println(strings.Repeat("#", 120))
	fmt.Println()
	fmt.Println("VIEW Data Detail Transaksi")
	fmt.Println(strings.Repeat("-", 120))
	fmt.Printf("%-5s | %-20s | %-20s | %-20s | %-20s | %-30s\n", "NO", "Pelayanan", "Jumlah", "Satuan", "Harga", "Total")
	fmt.Println(strings.Repeat("-", 120))

	// Query untuk mendapatkan data detail transaksi berdasarkan ID transaksi
	sqlStatement := "SELECT l.name, d.quantity, l.unit, l.price, SUM(l.price * d.quantity) AS total " +
		"FROM layanan AS l " +
		"JOIN detailTransaksi AS d ON l.id = d.layanan_id " +
		"WHERE d.transaksi_id = $1 " +
		"GROUP BY l.name, d.quantity, l.unit, l.price;"

	rows, err := db.Query(sqlStatement, transaksiID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var no int
	for rows.Next() {
		var name string
		var quantity int
		var price float64
		var unit string
		var total float64

		err := rows.Scan(&name, &quantity, &unit, &price, &total)
		if err != nil {
			panic(err)
		}

		no++
		fmt.Printf("%-5d | %-20s | %-20d | %-20s | %-20.2f | %-30.2f\n", no, name, quantity, unit, price, total)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Repeat("-", 120))
	fmt.Println()
	transaksiMenu(db)
}

// insert detail transaksi
func insertDetail(db *sql.DB, transaksiID string) {

	var layananID string
	var qty int

	fmt.Println("INSERT DETAIL TRANSAKSI")
	fmt.Print("Masukkan ID Layanan: ")
	_, err := fmt.Scan(&layananID)
	layananID = strings.ToUpper(layananID)
	if err != nil {
		panic(err)
	}

	// Validasi
	var existingLayanan string
	err = db.QueryRow("SELECT name FROM layanan WHERE id = $1", layananID).Scan(&existingLayanan)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println()
			fmt.Printf("Data Layanan dengan ID %s yang diberikan tidak ditemukan.\n", layananID)
			fmt.Println()
			insertDetail(db, transaksiID)
			return
		}
	}else {
		fmt.Println("Nama Layanan:", existingLayanan)
	}

	fmt.Print("Masukkan Jumlah: ")
	_, err = fmt.Scan(&qty)

	// validasi 2 tx.detail = jumlah tidak valid,
	if err != nil {
		fmt.Println("jumlah tidak valid, Harap masukkan angka.")
	}

	// validasi 3 tx.detail = Jumlah harus angka positif
	if qty <= 0 {
		fmt.Println("Jumlah harus angka positif.")
		insertDetail(db, transaksiID)
		return
	}

	var lastID string
	err = db.QueryRow("SELECT id FROM detailTransaksi ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1").Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	// Menghitung ID berikutnya
	nextID := "D1"
	if lastID != "" {
		// Ekstrak angka dari ID terakhir
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			panic(err)
		}
		// Tambahkan 1 ke angka tersebut
		nextNumber := lastNumber + 1
		// Format ID berikutnya
		nextID = fmt.Sprintf("D%d", nextNumber)
	}

	detail := entity.DetailTransaksi{Id: nextID, Transaksi_id: transaksiID, Layanan_id: layananID, Qty: qty}
	DetailTransaksi(detail)
	fmt.Println()
	fmt.Printf("Data detail transaksi %s berhasil dimasukkan.\n\n", nextID)

	//  Tambahkan pilihan "y/n" untuk menambahkan transaksi lagi
	fmt.Println("Apakah Anda ingin menambahkan detail transaksi lagi? (y/n): ")

	var choice string
	_, err = fmt.Scan(&choice)
	if err != nil {
		panic(err)
	}

	// validasi 4 tx.detail = mengulang insert di transaksi yang sama
	if strings.ToLower(choice) == "y" {
		insertDetail(db, transaksiID)
	} else if strings.ToLower(choice) == "n" {
		transaksiMenu(db)
	} else {
		fmt.Println("Pilihan tidak valid. Kembali ke menu utama.")
		main()
	}

}

func DetailTransaksi(detailTransaksi entity.DetailTransaksi) {
	db := connectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	insertDetailTransaksi(detailTransaksi, tx)

	total := getTotal(detailTransaksi.Transaksi_id, tx)

	updateTotal(total, detailTransaksi.Transaksi_id, tx)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func validate(err error, tx *sql.Tx) {

	if err != nil {
		tx.Rollback()
		fmt.Println(err, "Transaction Rollback")

	}
}

func insertDetailTransaksi(detailTransaksi entity.DetailTransaksi, tx *sql.Tx) {

	sqlStatement := "INSERT INTO detailTransaksi(id, transaksi_id, layanan_id, quantity) VALUES ($1, $2, $3, $4);"

	_, err := tx.Exec(sqlStatement, detailTransaksi.Id, detailTransaksi.Transaksi_id, detailTransaksi.Layanan_id, detailTransaksi.Qty)

	validate(err, tx)
}

func getTotal(id string, tx *sql.Tx) float64 {

	sqlStatement := "SELECT SUM(l.price * d.quantity) AS sub_total FROM detailTransaksi AS d JOIN layanan AS l ON d.layanan_id = l.id WHERE d.transaksi_id = $1;"

	total := 0.00

	err := tx.QueryRow(sqlStatement, id).Scan(&total)

	validate(err, tx)
	return total
}

func updateTotal(total float64, customer_id string, tx *sql.Tx) {

	sqlStatement := "UPDATE transaksi SET total = $1 WHERE id = $2;"

	_, err := tx.Exec(sqlStatement, total, customer_id)
	validate(err, tx)
}
