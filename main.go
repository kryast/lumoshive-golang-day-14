package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Pesanan struct {
	ID         int      `json:"id"`
	Items      []string `json:"items"`
	TotalHarga int      `json:"total_harga"`
	Status     string   `json:"status"`
}

type MenuItem struct {
	Nama  string
	Harga int
}

var DaftarPesanan []Pesanan
var IdSekarang int

var menu = []MenuItem{
	{"Nasi Lele", 15000},
	{"Bakso", 12000},
	{"Ayam Goreng", 20000},
	{"Air Mineral", 5000},
	{"Es Teh", 8000},
}

func main() {
	Menu()
}

func Menu() {
	var input int
	fmt.Println("=== Menu Utama ===")
	fmt.Println("1. Tambah Pesanan")
	fmt.Println("2. Edit Pesanan")
	fmt.Println("3. Checkout")
	fmt.Println("4. History Pesanan")
	fmt.Println("99. Exit")

	fmt.Print("Pilih menu: ")
	fmt.Scan(&input)

	switch input {
	case 1:
		ClearScreen()
		TambahPesanan()
	case 2:
		ClearScreen()
		EditPesanan()
	case 3:
		ClearScreen()
		Checkout()
	case 4:
		ClearScreen()
		HistoryPesanan()
	case 99:
		ClearScreen()
	default:
		fmt.Println("Error, pilih menu yang valid.")
		Menu()
	}
}

func TambahPesanan() {
	ClearScreen()
	var pesanan Pesanan
	IdSekarang++
	pesanan.ID = IdSekarang

	fmt.Println("=== Daftar Menu ===")
	for i, item := range menu {
		fmt.Printf("%d. %s - Rp %d\n", i+1, item.Nama, item.Harga)
	}
	fmt.Println("0. Kembali")

	for {
		var input int
		fmt.Print("Masukkan nomor menu: ")
		fmt.Scan(&input)

		if input == 0 {
			break
		}

		if input > 0 && input <= len(menu) {
			item := menu[input-1]
			pesanan.Items = append(pesanan.Items, item.Nama)
			pesanan.TotalHarga += item.Harga
			fmt.Printf("'%s' berhasil ditambahkan ke pesanan.\n", item.Nama)
		} else {
			fmt.Println("Tidak Tersedia. Silakan coba lagi.")
		}
	}

	if len(pesanan.Items) > 0 {
		pesanan.Status = "Di proses"
		DaftarPesanan = append(DaftarPesanan, pesanan)
		ClearScreen()
		fmt.Println("Pesanan berhasil ditambahkan.")
	} else {
		ClearScreen()
		fmt.Println("Tidak ada menu yang dipesan.")
	}

	Menu()
}

func EditPesanan() {
	ClearScreen()
	fmt.Println("=== Daftar Pesanan yang Bisa Diedit ===")
	for _, item := range DaftarPesanan {
		if item.Status != "Selesai" {
			fmt.Printf("ID: %d | daftar yang dipesan : %s | Total Harga: Rp %d\n", item.ID, item.Items[0:], item.TotalHarga)
		}
	}

	if len(DaftarPesanan) == 0 {
		fmt.Println("Tidak ada pesanan yang dapat diedit.")
		Menu()
		return
	}
	fmt.Println("0. Kembali")

	var id int
	fmt.Print("\nMasukkan ID pesanan yang ingin diubah: ")
	fmt.Scan(&id)

	for i := range DaftarPesanan {
		pesanan := &DaftarPesanan[i]
		if pesanan.ID == id && pesanan.Status != "Selesai" {
			ClearScreen()
			fmt.Println("=== Daftar Menu yang Dipesan ===")
			for j, NamaItem := range pesanan.Items {
				fmt.Printf("%d. %s\n", j+1, NamaItem)
			}

			fmt.Println("\n=== Daftar Menu ===")
			for j, item := range menu {
				fmt.Printf("%d. %s - Rp %d\n", j+1, item.Nama, item.Harga)
			}

			var index int
			fmt.Print("Masukkan nomor menu yang ingin diubah: ")
			fmt.Scan(&index)

			if index > 0 && index <= len(pesanan.Items) {
				hapusItem := pesanan.Items[index-1]
				pesanan.Items = append(pesanan.Items[:index-1], pesanan.Items[index:]...)
				pesanan.TotalHarga -= HargaItem(hapusItem)

				var indexBaru int
				fmt.Print("Masukkan nomor menu item baru yang ingin ditambahkan: ")
				fmt.Scan(&indexBaru)

				if indexBaru > 0 && indexBaru <= len(menu) {
					item := menu[indexBaru-1]
					pesanan.Items = append(pesanan.Items, item.Nama)
					pesanan.TotalHarga += item.Harga

					ClearScreen()
					fmt.Println("Pesanan berhasil diubah!")
				} else {
					ClearScreen()
					fmt.Println("Tidak Tersedia.")
				}
			} else {
				ClearScreen()
				fmt.Println("Tidak Tersedia.")
			}
			Menu()
			return
		}
	}

	ClearScreen()
	Menu()
}

func HargaItem(NamaItem string) int {
	for _, item := range menu {
		if item.Nama == NamaItem {
			return item.Harga
		}
	}
	return 0
}

func Checkout() {
	ClearScreen()
	fmt.Println("=== Daftar Pesanan untuk Checkout ===")
	jumlahPesanan := 0
	for _, pesanan := range DaftarPesanan {
		if pesanan.Status != "Selesai" {
			fmt.Printf("ID: %d | daftar yang dipesan : %s | Total Harga: Rp %d\n", pesanan.ID, pesanan.Items[0:], pesanan.TotalHarga)
			jumlahPesanan++
		}
	}

	if jumlahPesanan == 0 {
		fmt.Println("Tidak ada pesanan yang tersedia untuk checkout.")
		Menu()
		return
	}
	fmt.Println("0. Kembali")

	var id int
	fmt.Print("\nMasukkan ID pesanan untuk checkout: ")
	fmt.Scan(&id)

	for i, pesanan := range DaftarPesanan {
		if pesanan.ID == id && pesanan.Status != "Selesai" {
			fmt.Println("\n=== Menu yang Dipesan ===")
			for _, item := range pesanan.Items {
				fmt.Printf(" - %s\n", item)
			}

			DaftarPesanan[i].Status = "Di antar"
			fmt.Println("Pesanan telah dibayar.")
			time.Sleep(2 * time.Second)

			DaftarPesanan[i].Status = "Selesai"
			fmt.Println("Pesanan telah diterima.")
			Menu()
			return
		}
	}
	if id == 0 {
		ClearScreen()
		Menu()
	}

	ClearScreen()
	fmt.Println("ID pesanan tidak ditemukan atau sudah selesai.")
	Menu()
}

func HistoryPesanan() {
	ClearScreen()
	fmt.Println("=== History Pesanan ===")

	if len(DaftarPesanan) == 0 {
		fmt.Println("History kosong.")
		Menu()
		return
	}

	for _, pesanan := range DaftarPesanan {
		fmt.Printf("ID: %d | daftar yang dipesan : %s | Total Harga: Rp %d | Status : %s\n", pesanan.ID, pesanan.Items[0:], pesanan.TotalHarga, pesanan.Status)
	}

	fmt.Println("0. Kembali")
	var input int
	fmt.Print("Masukkan angka: ")
	fmt.Scan(&input)

	if input == 0 {
		ClearScreen()
		Menu()
	}
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
