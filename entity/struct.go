package entity

import "time"

type Customers struct {
	Id      string
	Name    string
	Address string
	Phone   string
}

type Karyawan struct {
	Id      string
	Name    string
	Address string
	Phone   string
}

type Layanan struct {
	Id    string
	Name  string
	Unit  string
	Price float64
}

type Transaksi struct {
	Id          string
	Customer_id string
	Karyawan_id string
	StartDate   time.Time
	EndDate     time.Time
	Total       float64
}

type DetailTransaksi struct {
	Id           string
	Transaksi_id string
	Layanan_id   string
	Qty          int
}
