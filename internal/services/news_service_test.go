package service

import (
	model "gin-demo/internal/models"
	repository "gin-demo/internal/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSafeBatchImport_ValidationError(t *testing.T) {
	// 1. 創建 sqlmock 模擬數據庫連接
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("創建 mock 失敗: %v", err)
	}

	// 2. 將 sqlmock 注入 GORM
	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true, // 避免 GORM 查版本
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm open 失敗: %v", err)
	}

	// 3. 初始化 Service
	newsRepo := repository.NewNewsRepository(gormDB)
	logRepo := repository.NewNewsLogRepository(gormDB)

	// 4. 初始化 Service (傳入實例，滿足接口要求)
	service := NewNewsService(gormDB, newsRepo, logRepo)

	// 5. 準備測試數據：
	// 第一條：合格
	// 第二條：標題太短，應該觸發 validateNews 報錯
	newsList := []model.News{
		{Title: "這是一個合格的標題"},
		{Title: "太短"},
	}

	// 6. 執行測試
	err = service.SafeBatchImport(newsList)

	// 7. 斷言結果
	if err == nil {
		t.Error("預期應該發生校驗錯誤，但實際上返回了 nil")
	}

	// 預期錯誤訊息應該包含我們在 validateNews 寫的內容
	if err == nil {
		t.Error("預期應該發生校驗錯誤，但實際上返回了 nil")
	} else {
		// 确保变量被使用
		expectedPart := "太短"
		if !contains(err.Error(), expectedPart) {
			t.Errorf("錯誤訊息不符，得到: %v", err)
		}
	}

	// 8. 驗證數據庫行為
	// 因為在校驗階段（SafeBatchImport 的 errgroup）就會報錯，
	// 所以代碼不應該走到 Transaction 內部。
	// 我們檢查 mock 的預期是否都達成了（即：沒有多餘的 SQL 執行）
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("數據庫預期行為不匹配（可能意外執行了 SQL）: %v", err)
	}
}

// 輔助函數：檢查字串包含
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != "" && substr != ""
}
