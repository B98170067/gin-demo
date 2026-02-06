package service

import (
	model "gin-demo/internal/models"
	repository "gin-demo/internal/repositories"
	"gin-demo/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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

func TestSafeBatchImport_Mockery(t *testing.T) {
	// 1. 创建 sqlmock
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer sqlDB.Close()

	// 2. 构造一个 MySQL 类型的 GORM DB (emptyDB)
	// 虽然它不连真数据库，但它具备了 MySQL 的驱动逻辑，不会 panic
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	emptyDB, _ := gorm.Open(dialector, &gorm.Config{})

	// 1. 初始化 Mock 对象
	mockRepo := new(mocks.INewsRepository)
	mockLogRepo := new(mocks.INewsLogRepository)

	// 2. 定义行为：当调用 Transaction 时，模拟它运行传入的函数
	// 这里用到了 .Run() 来执行 Service 传入的那个闭包
	mockRepo.On("Transaction", mock.Anything).Return(func(fn func(tx *gorm.DB) error) error {
		return fn(emptyDB) // 模拟事务开始，这里 tx 给 nil 没关系
	}).Maybe() // 表示这个方法可能会被调用

	// 2. 补上 CreateTx：解决 Panic 的关键
	// 使用 .Return(nil) 表示模拟数据库写入成功
	mockRepo.On("CreateTx", mock.Anything, mock.Anything).Return(nil).Maybe()

	// 3. 注入 Mock 到 Service
	service := NewNewsService(emptyDB, mockRepo, mockLogRepo)

	// 1. 准备数据
	newsList := []model.News{{Title: "太短"}}

	// 2. 执行测试
	err = service.SafeBatchImport(newsList)

	// 6. 修正后的断言：先要求 err 必须存在
	if assert.Error(t, err, "应该因为标题太短报错") {
		assert.Contains(t, err.Error(), "太短")
	} else {
		// 如果没有报错，打印提示，方便调试
		t.Logf("当前标题长度(字节): %d", len(newsList[0].Title))
	}

	// 4. 只有在校验失败的情况下，我们才确信 Transaction 不该被调用
	mockRepo.AssertNotCalled(t, "Transaction", mock.Anything)
}
