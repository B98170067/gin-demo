package service_test

import (
	model "gin-demo/internal/models"
	repository "gin-demo/internal/repositories"
	service "gin-demo/internal/services"
	"gin-demo/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// --- 輔助工具：創建一個不會崩潰的 GORM DB 實例 ---
func setupEmptyDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("創建 sqlmock 失敗: %v", err)
	}
	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("GORM 開啟失敗: %v", err)
	}
	return gormDB, mock
}

// 1. 測試：校驗失敗 (應該攔截在 Transaction 之前)
func TestSafeBatchImport_ValidationError(t *testing.T) {
	gormDB, mockDB := setupEmptyDB(t)

	// 使用真實的 Repository 來測試內部的 validateNews
	newsRepo := repository.NewNewsRepository(gormDB)
	logRepo := repository.NewNewsLogRepository(gormDB)
	svc := service.NewNewsService(gormDB, newsRepo, logRepo)

	newsList := []model.News{
		{Title: "這是一個合格的標題"},
		{Title: "太短"}, // 此條會觸發 RuneCountInString < 5
	}

	err := svc.SafeBatchImport(newsList)

	// 斷言：應該報錯且包含關鍵字
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "太短")
	}

	// 驗證數據庫行為：因為校驗失敗，不應有任何 SQL 執行
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

// 2. 測試：分批寫入成功 (Mockery 模式)
func TestSafeBatchImport_Mockery_BatchSuccess(t *testing.T) {
	emptyDB, _ := setupEmptyDB(t)

	mockRepo := new(mocks.INewsRepository)
	mockLogRepo := new(mocks.INewsLogRepository)
	svc := service.NewNewsService(emptyDB, mockRepo, mockLogRepo)

	// 準備 3 筆合格資料
	newsList := []model.News{
		{Title: "標題長度合格一號"},
		{Title: "標題長度合格二號"},
		{Title: "標題長度合格三號"},
	}

	// 定義行為：
	// 因為 3 筆小於預設的 batchSize (500)，所以 Transaction 只會呼叫 1 次
	mockRepo.On("Transaction", mock.Anything).Return(func(fn func(tx *gorm.DB) error) error {
		return fn(emptyDB) // 執行傳入的閉包
	}).Once()

	// CreateTx 會針對每一筆資料被呼叫
	mockRepo.On("CreateTx", mock.Anything, mock.Anything).Return(nil).Times(3)

	err := svc.SafeBatchImport(newsList)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// 3. 測試：Transaction 失敗時的分批處理
func TestSafeBatchImport_Mockery_TransactionFailure(t *testing.T) {
	emptyDB, _ := setupEmptyDB(t)

	mockRepo := new(mocks.INewsRepository)
	mockLogRepo := new(mocks.INewsLogRepository)
	svc := service.NewNewsService(emptyDB, mockRepo, mockLogRepo)

	newsList := []model.News{
		{Title: "標題長度合格一號"},
	}

	// 模擬 Transaction 返回錯誤（例如資料庫斷線）
	mockRepo.On("Transaction", mock.Anything).Return(func(fn func(tx *gorm.DB) error) error {
		return gorm.ErrInvalidDB // 模擬回傳錯誤
	}).Once()

	err := svc.SafeBatchImport(newsList)

	// 斷言：錯誤訊息應該包含分批寫入中斷的訊息
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "分批寫入中斷")
	}
}

// 4. 測試：空列表處理
func TestSafeBatchImport_EmptyList(t *testing.T) {
	svc := service.NewNewsService(nil, nil, nil)
	err := svc.SafeBatchImport([]model.News{})

	// 根據你的邏輯，空列表 errgroup 會直接 Wait() 結束並 return nil
	assert.NoError(t, err)
}
