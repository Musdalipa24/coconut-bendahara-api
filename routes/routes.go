package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/controller"
	"github.com/syrlramadhan/api-bendahara-inovdes/middleware"
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
	router.GET("/api/admin/:username", middleware.VerifyJWT(adminController.FindByUsername))
	router.PUT("/api/admin/update/:username", middleware.VerifyJWT(adminController.UpdateAdmin))

	// pemasukan
	pemasukanRepo := repository.NewPemasukanRepo()
	pemasukanService := service.NewPemasukanService(pemasukanRepo, db)
	pemasukanController := controller.NewPemasukanController(pemasukanService)

	router.POST("/api/pemasukan/add", middleware.VerifyJWT(pemasukanController.AddPemasukan))
	router.PUT("/api/pemasukan/update/:id", middleware.VerifyJWT(pemasukanController.UpdatePemasukan))
	router.GET("/api/pemasukan/getall", middleware.VerifyJWT(pemasukanController.GetPemasukan))
	router.GET("/api/pemasukan/get/:id", middleware.VerifyJWT(pemasukanController.GetById))
	router.DELETE("/api/pemasukan/delete/:id", middleware.VerifyJWT(pemasukanController.DeletePemasukan))

	// pengeluaran
	pengeluaranRepo := repository.NewPengeluaranRepo()
	pengeluaranService := service.NewPengeluaranService(pengeluaranRepo, db)
	pengeluaranController := controller.NewPengeluaranController(pengeluaranService)

	router.POST("/api/pengeluaran/add", middleware.VerifyJWT(pengeluaranController.AddPengeluaran))
	router.PUT("/api/pengeluaran/update/:id", middleware.VerifyJWT(pengeluaranController.UpdatePengeluaran))
	router.GET("/api/pengeluaran/getall", middleware.VerifyJWT(pengeluaranController.GetPengeluaran))
	router.GET("/api/pengeluaran/get/:id", middleware.VerifyJWT(pengeluaranController.GetById))
	router.DELETE("/api/pengeluaran/delete/:id", middleware.VerifyJWT(pengeluaranController.DeletePengeluaran))

	// transaksi
	transactionRepo := repository.NewTransactionRepo()
	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionController := controller.NewTransactionController(transactionService)

	router.GET("/api/transaksi/getall", middleware.VerifyJWT(transactionController.GetAllTransaction))
	router.GET("/api/transaksi/getlast", middleware.VerifyJWT(transactionController.GetLastTransaction))

	// laporan keuangan
	laporanKeuanganRepo := repository.NewLaporanKeuanganRepo()
	laporanKeuanganService := service.NewLaporanKeuanganService(laporanKeuanganRepo, db)
	laporanKeuanganController := controller.NewLaporanKeuanganController(laporanKeuanganService)

	router.GET("/api/laporan/getall", middleware.VerifyJWT(laporanKeuanganController.GetAllLaporan))
	router.GET("/api/laporan/saldo", middleware.VerifyJWT(laporanKeuanganController.GetLastBalance))
	router.GET("/api/laporan/pengeluaran", middleware.VerifyJWT(laporanKeuanganController.GetTotalExpenditure))
	router.GET("/api/laporan/pemasukan", middleware.VerifyJWT(laporanKeuanganController.GetTotalIncome))
	router.GET("/api/laporan/range", middleware.VerifyJWT(laporanKeuanganController.GetLaporanByDateRange))

	iuranRepo := repository.NewIuranRepository()
	iuranService := service.NewIuranService(iuranRepo, db)
	iuranController := controller.NewIuranController(iuranService)

	router.POST("/api/member/add", middleware.VerifyJWT(iuranController.CreateMember))
	router.GET("/api/member/getall", middleware.VerifyJWT(iuranController.GetAllMembers))
	router.GET("/api/member/get/:id", middleware.VerifyJWT(iuranController.GetMemberById))
	router.PUT("/api/member/update/:id_member", middleware.VerifyJWT(iuranController.UpdateIuran))
	router.DELETE("/api/member/delete/:id_member", middleware.VerifyJWT(iuranController.DeleteMember))
	router.PUT("/api/member/reactivate/:id_member", middleware.VerifyJWT(iuranController.ReactivateMember))

	// file uploads
	router.ServeFiles("/api/uploads/*filepath", http.Dir("./uploads/"))
	router.ServeFiles("/uploads/*filepath", http.Dir("./uploads/"))

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
	allowedOrigins := map[string]bool{
		"http://localhost:3000":                    true,
		"https://bendahara-coconut.vercel.app":     true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}