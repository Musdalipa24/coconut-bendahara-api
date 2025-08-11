package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/controller"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
)

func Routes(db *sql.DB, port string) {
	router := httprouter.New()

	// admin
	adminRepo := repository.NewAdminRepo()
	adminService := service.NewAdminService(adminRepo, db)
	adminController := controller.NewAdminController(adminService)

	router.POST("/api/admin/daftar", adminController.SignUp)
	router.POST("/api/admin/login", adminController.SignIn)
	router.GET("/api/admin/:nik", adminController.FindByNik)

	// pemasukan
	pemasukanRepo := repository.NewPemasukanRepo()
	pemasukanService := service.NewPemasukanService(pemasukanRepo, db)
	pemasukanController := controller.NewPemasukanController(pemasukanService)

	router.POST("/api/pemasukan/add", pemasukanController.AddPemasukan)
	router.PUT("/api/pemasukan/update/:id", pemasukanController.UpdatePemasukan)
	router.GET("/api/pemasukan/getall", pemasukanController.GetPemasukan)
	router.GET("/api/pemasukan/get/:id", pemasukanController.GetById)
	router.DELETE("/api/pemasukan/delete/:id", pemasukanController.DeletePemasukan)

	// pengeluaran
	pengeluaranRepo := repository.NewPengeluaranRepo()
	pengeluaranService := service.NewPengeluaranService(pengeluaranRepo, db)
	pengeluaranController := controller.NewPengeluaranController(pengeluaranService)

	router.POST("/api/pengeluaran/add", pengeluaranController.AddPengeluaran)
	router.PUT("/api/pengeluaran/update/:id", pengeluaranController.UpdatePengeluaran)
	router.GET("/api/pengeluaran/getall", pengeluaranController.GetPengeluaran)
	router.GET("/api/pengeluaran/get/:id", pengeluaranController.GetById)
	router.DELETE("/api/pengeluaran/delete/:id", pengeluaranController.DeletePengeluaran)

	// transaksi
	transactionRepo := repository.NewTransactionRepo()
	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionController := controller.NewTransactionController(transactionService)

	router.GET("/api/transaksi/getall", transactionController.GetAllTransaction)
	router.GET("/api/transaksi/getlast", transactionController.GetLastTransaction)

	// laporan keuangan
	laporanKeuanganRepo := repository.NewLaporanKeuanganRepo()
	laporanKeuanganService := service.NewLaporanKeuanganService(laporanKeuanganRepo, db)
	laporanKeuanganController := controller.NewLaporanKeuanganController(laporanKeuanganService)

	router.GET("/api/laporan/getall", laporanKeuanganController.GetAllLaporan)
	router.GET("/api/laporan/saldo", laporanKeuanganController.GetLastBalance)
	router.GET("/api/laporan/pengeluaran", laporanKeuanganController.GetTotalExpenditure)
	router.GET("/api/laporan/pemasukan", laporanKeuanganController.GetTotalIncome)
	router.GET("/api/laporan/range", laporanKeuanganController.GetLaporanByDateRange)

	// iuran
	iuranRepo := repository.NewIuranRepo(db)
	iuranService := service.NewIuranService(iuranRepo)
	iuranController := controller.NewIuranController(iuranService)

	router.POST("/api/iuran/add", iuranController.AddIuran)
	router.PUT("/api/iuran/update/:id", iuranController.UpdateIuran)
	router.GET("/api/iuran/getall", iuranController.GetAllIuran)
	router.GET("/api/iuran/get/:id", iuranController.GetIuranById)
	router.DELETE("/api/iuran/delete/:id", iuranController.DeleteIuran)

	// file uploads
	router.ServeFiles("/api/uploads/*filepath", http.Dir("./uploads/"))

	handler := corsMiddleware(router)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	errServer := server.ListenAndServe()
	if errServer != nil {
		panic(errServer)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
