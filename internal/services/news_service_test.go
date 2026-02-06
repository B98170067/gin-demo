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
	// 1. 创建 sqlmock 模拟数据库连接
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建 mock 失败: %v", err)
	}

	// 2. 将 sqlmock 注入 GORM
	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true, // ⭐ 很重要，避免 GORM 查版本
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm open 失败: %v", err)
	}

	// 3. 初始化 Service
	newsRepo := repository.NewNewsRepository(gormDB)
	// 假设 NewsLogRepository 类似
	service := NewNewsService(gormDB, newsRepo, nil)

	// 4. 准备测试数据：第二条标题太短，应该触发校验失败
	newsList := []model.News{
		{Title: "这是一个合格的标题"},
		{Title: "太短"},
	}

	// 5. 执行测试
	err = service.SafeBatchImport(newsList)

	// 6. 断言结果
	if err == nil {
		t.Error("预期应该发生校验错误，但实际上返回了 nil")
	}

	// 7. 验证数据库是否真的没有执行写入（事务不应该开启）
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("数据库预期行为不匹配: %v", err)
	}
}
