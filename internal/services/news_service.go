package service

import (
	"context"
	"fmt"
	model "gin-demo/internal/models"
	repository "gin-demo/internal/repositories"
	"time"
	"unicode/utf8"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type NewsService struct {
	db      *gorm.DB
	repo    repository.INewsRepository
	logRepo repository.INewsLogRepository
}

func NewNewsService(
	db *gorm.DB,
	repo repository.INewsRepository,
	logRepo repository.INewsLogRepository,
) *NewsService {
	return &NewsService{
		db:      db,
		repo:    repo,
		logRepo: logRepo,
	}
}

func (s *NewsService) GetAllNews() ([]model.News, error) {
	return s.repo.FindAll()
}

func (s *NewsService) GetPaged(page, size int, status *int) ([]model.News, int64) {
	return s.repo.FindPaged(page, size, status)
}

func (s *NewsService) CreateNews(news *model.News) error {
	return s.repo.Create(news)
}

func (s *NewsService) CreateWithLog(news *model.News) error {
	return s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateTx(tx, news); err != nil {
			return err // rollback
		}

		log := &model.NewsLog{
			NewsID: news.ID,
			Action: "CREATE",
		}

		if err := s.logRepo.CreateTx(tx, log); err != nil {
			return err // rollback
		}

		return nil // commit
	})
}

func (s *NewsService) GetByID(id uint) (*model.News, error) {
	return s.repo.FindByID(id)
}

func (s *NewsService) Update(news *model.News) error {
	return s.repo.Update(news)
}

func (s *NewsService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// SafeBatchImport 高性能批量导入新闻
// 逻辑：并发校验所有数据 -> 全部通过 -> 开启事务一次性写入
func (s *NewsService) SafeBatchImport(newsList []model.News) error {
	// 主 Context：負責整體的 5 秒超時
	mainCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 校驗用的 Context：僅用於 errgroup 內部的併發控制
	g, checkCtx := errgroup.WithContext(mainCtx)
	for i := range newsList {
		// 闭包陷阱：必须重新定义变量，否则协程内拿到的都是最后一条数据
		news := newsList[i]
		g.Go(func() error {
			// 检查 Context 是否已经因为其他协程报错或超时而取消
			select {
			case <-checkCtx.Done():
				return checkCtx.Err()
			default:
				// 3. 执行并发校验逻辑（如标题长度、重复性检查）
				return s.validateNews(news)
			}
		})
	}

	// 4. 等待所有校验 Goroutine 完成
	if err := g.Wait(); err != nil {
		return fmt.Errorf("批量校验失败或超时: %w", err)
	}

	/*
		// 5. 校验全部通过，开启事务写入数据库
		return s.repo.Transaction(func(tx *gorm.DB) error {
			for _, n := range newsList {
				// 使用 WithContext 让 GORM 也能感知超时
				if err := s.repo.CreateTx(tx.WithContext(ctx), &n); err != nil {
					return err
				}
			}
			return nil
		})
	*/

	// === 5. 分批寫入邏輯開始 ===
	batchSize := 500 // 建議每 500~1000 筆一組
	for i := 0; i < len(newsList); i += batchSize {
		// 計算這一批的結尾索引
		end := i + batchSize
		if end > len(newsList) {
			end = len(newsList)
		}

		chunk := newsList[i:end]

		// 每一批開啟一個獨立的事務
		err := s.repo.Transaction(func(tx *gorm.DB) error {
			for _, n := range chunk {
				// 這裡改用 mainCtx，它不會因為 errgroup 結束而被取消
				if err := mainCtx.Err(); err != nil {
					return err
				}
				if err := s.repo.CreateTx(tx.WithContext(mainCtx), &n); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("分批寫入中斷 (起始索引 %d): %w", i, err)
		}
	}

	return nil
}

// validateNews 模拟具体的校验逻辑
func (s *NewsService) validateNews(n model.News) error {
	// 模拟耗时的 IO 或计算操作（如敏感词过滤）
	time.Sleep(50 * time.Millisecond)

	if n.Title == "" {
		return fmt.Errorf("新闻标题不能为空")
	}

	// 使用 RuneCountInString 计算真实的字符数
	if utf8.RuneCountInString(n.Title) < 5 {
		return fmt.Errorf("标题 '%s' 太短了，至少需要 5 个字", n.Title)
	}
	return nil
}
