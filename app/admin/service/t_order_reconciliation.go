package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/go-admin-team/go-admin-core/storage"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/global"

	adminLogger "github.com/go-admin-team/go-admin-core/logger"
	"github.com/mitchellh/mapstructure"
)

type OrderReconciliation struct {
	service.Service
	OrderDB *gorm.DB
}

func (e *OrderReconciliation) InitDB() {
	if err := e.InitOrderDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized.")
}

// GetPage 获取OrderReconciliation列表
func (e *OrderReconciliation) GetPage(c *dto.OrderReconciliationGetPageReq, p *actions.DataPermission, list *[]models.OrderReconciliation, count *int64) error {
	var err error
	var data models.OrderReconciliation

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("OrderReconciliationService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取OrderReconciliation对象
func (e *OrderReconciliation) Get(d *dto.OrderReconciliationGetReq, p *actions.DataPermission, model *models.OrderReconciliation) error {
	var data models.OrderReconciliation

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetOrderReconciliation error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建OrderReconciliation对象
func (e *OrderReconciliation) Insert(c *dto.OrderReconciliationInsertReq) error {
	var err error
	var data models.OrderReconciliation
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("OrderReconciliationService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改OrderReconciliation对象
func (e *OrderReconciliation) Update(c *dto.OrderReconciliationUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.OrderReconciliation{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrderReconciliationService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除OrderReconciliation
func (e *OrderReconciliation) Remove(d *dto.OrderReconciliationDeleteReq, p *actions.DataPermission) error {
	var data models.OrderReconciliation

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveOrderReconciliation error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *OrderReconciliation) MergeData() (logZfbNotify []*models.LogZfbNotify, err error) {
	err = e.OrderDB.Find(&logZfbNotify).Error

	for _, l := range logZfbNotify {
		l.ParmMap, err = l.ParseParm()
		if err != nil {
			continue
		}
		if msg, ok := l.ParmMap["msg"]; ok && msg.(string) == "退款成功" {
			if data, ok1 := l.ParmMap["data"]; ok1 {
				l.ParmMap = data.(map[string]interface{})
			}

		}
		//l.ParmMap中map key为gmt_create,gmt_payment,notify_time转成time.Time类型
		for k, v := range l.ParmMap {
			if l.Type == "挂号退款" || l.Type == "门诊缴费退款" || l.Type == "住院预缴退费" || l.Type == "挂号医保退款" {
				if k == "gmt_refund_pay" {
					t, err := time.Parse("2006-01-02 15:04:05", v.(string))
					if err != nil {
						continue
					}
					l.ParmMap[k] = t
				}
			} else {
				if k == "gmt_create" || k == "gmt_payment" || k == "notify_time" {
					t, err := time.Parse("2006-01-02 15:04:05", v.(string))
					if err != nil {
						continue
					}
					l.ParmMap[k] = t
				}
			}

		}
	}

	//把logZfbNotify中的ParmMap中的数据批量并发插入到AlipayNotify表中
	var wg sync.WaitGroup
	wg.Add(len(logZfbNotify))
	for _, l := range logZfbNotify {
		go func(l *models.LogZfbNotify) {
			defer wg.Done()
			tx := e.OrderDB.Begin()
			//l.Type == "挂号医保退款"  暂不支持
			if l.Type == "挂号退款" || l.Type == "门诊缴费退款" || l.Type == "住院预缴退费" || l.Type == "挂号医保退款" {
				var AlipayRefundNotify models.AlipayRefundNotify
				//l.ParmMap["refund_detail_item_list"]转成json字符串
				if details, ok := l.ParmMap["refund_detail_item_list"]; ok {
					l.ParmMap["refund_details"], _ = json.Marshal(details)
					//删除refund_detail_item_list
					delete(l.ParmMap, "refund_detail_item_list")
				}
				// 将map中的数据存储到AlipayNotify结构体中
				err := mapstructure.WeakDecode(l.ParmMap, &AlipayRefundNotify)
				if err != nil {
					return
				}
				if AlipayRefundNotify.OutTradeNo == "" {
					return
				}
				AlipayRefundNotify.TypeName = l.Type
				err = tx.Create(&AlipayRefundNotify).Error
				if err != nil {
					tx.Rollback()
				}
			} else {
				var AlipayNotify models.AlipayNotify

				// 将map中的数据存储到AlipayNotify结构体中
				err := mapstructure.WeakDecode(l.ParmMap, &AlipayNotify)
				if err != nil {
					return
				}
				if AlipayNotify.OutTradeNo == "" {
					return
				}
				AlipayNotify.TypeName = l.Type
				err = tx.Create(&AlipayNotify).Error
				if err != nil {
					tx.Rollback()
				}
			}

			tx.Commit()
		}(l)
	}
	wg.Wait()

	return logZfbNotify, err
}

func (e *OrderReconciliation) MergeDataWx(c *gin.Context) error {
	l := make(map[string]interface{})
	l["pageIndex"] = 2
	l["pageSize"] = 30
	q := sdk.Runtime.GetMemoryQueue(c.Request.Host)
	message, err := sdk.Runtime.GetStreamMessage("", global.RhOrder, l)
	if err != nil {
		adminLogger.Errorf("GetStreamMessage error, %s", err.Error())
		//报错错误，不中断请求
	} else {
		err = q.Append(message)
		if err != nil {
			adminLogger.Errorf("Append message error, %s", err.Error())
		}
	}
	return err
}

func (e *OrderReconciliation) MergeDataWxQueue() (logWxNotify []*models.LogWxNotify, err error) {

	err = e.OrderDB.Find(&logWxNotify).Error

	for _, l := range logWxNotify {
		l.ParmMap, err = l.ParseParm()
		if err != nil {
			continue
		}
		for k, v := range l.ParmMap {
			if k == "success_time" {
				t, err := time.Parse(time.RFC3339, v.(string))
				if err != nil {
					continue
				}
				l.ParmMap[k] = t
			}
		}
	}

	//把logWxNotify中的ParmMap中的数据批量并发插入到WechatNotify表中
	// var wg sync.WaitGroup
	// wg.Add(len(logWxNotify))
	var WechatRefundNotifys []models.WxRefundNotify
	var WechatNotifys []models.WxNotify
	tx := e.OrderDB.Session(&gorm.Session{})
	refoundSize := 0
	paySize := 0
	payIds := make([]int, 0)
	for _, l := range logWxNotify {
		// go func(l *models.LogWxNotify) {
		//defer wg.Done()
		if outRefundNo, ok := l.ParmMap["out_refund_no"]; ok {
			outRefundNo = outRefundNo.(string)
			var WechatRefundNotify models.WxRefundNotify
			if outRefundNo == "" {
				refoundSize++
				continue
				//adminLogger.Errorf("rollback err, %s", err.Error())
			}
			amount := l.ParmMap["amount"].(map[string]interface{})
			WechatRefundNotify.Total = int(amount["total"].(float64))
			WechatRefundNotify.Refund = int(amount["refund"].(float64))
			WechatRefundNotify.PayerTotal = int(amount["payer_total"].(float64))
			WechatRefundNotify.PayerRefund = int(amount["payer_refund"].(float64))

			//l.ParmMap["refund_detail_item_list"]转成json字符串
			// if details, ok := l.ParmMap["refund_detail_item_list"]; ok {
			// 	l.ParmMap["refund_details"], _ = json.Marshal(details)
			// 	//删除refund_detail_item_list
			// 	delete(l.ParmMap, "refund_detail_item_list")
			// }
			// 将map中的数据存储到WechatNotify结构体中
			err := mapstructure.WeakDecode(l.ParmMap, &WechatRefundNotify)
			if err != nil {
				refoundSize++
				adminLogger.Errorf("WeakDecode err, %s", err.Error())
				continue
			}
			typeName := l.Type
			if l.Type == "微信回调" {
				typeName = "微信退款回调"
			}
			WechatRefundNotify.TypeName = typeName
			WechatRefundNotifys = append(WechatRefundNotifys, WechatRefundNotify)
		} else {
			var WechatNotify models.WxNotify

			// 将map中的数据存储到WechatNotify结构体中
			err := mapstructure.WeakDecode(l.ParmMap, &WechatNotify)
			if err != nil {
				paySize++
				adminLogger.Errorf("WeakDecode err, %s", err.Error())
				continue
			}
			if WechatNotify.OutTradeNo == "" {
				paySize++
				payIds = append(payIds, l.Id)
				//adminLogger.Errorf("OutTradeNo err, %s", err.Error())
				continue
			}

			amount := l.ParmMap["amount"].(map[string]interface{})
			WechatNotify.Total = int(amount["total"].(float64))
			WechatNotify.Currency = amount["currency"].(string)
			WechatNotify.PayerTotal = int(amount["payer_total"].(float64))
			WechatNotify.PayerCurrency = amount["payer_currency"].(string)

			WechatNotify.TypeName = l.Type
			WechatNotifys = append(WechatNotifys, WechatNotify)
		}
		//}(l)
	}
	//循环每次插入100条数据
	batchSize := 100
	var wg sync.WaitGroup
	errCh := make(chan error, 2) // Channel to receive errors

	// Batch creation for WechatRefundNotifys
	for i := 0; i < len(WechatRefundNotifys); i += batchSize {
		end := i + batchSize
		if end > len(WechatRefundNotifys) {
			end = len(WechatRefundNotifys)
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			err := tx.CreateInBatches(WechatRefundNotifys[start:end], len(WechatRefundNotifys[start:end])).Error
			if err != nil {
				errCh <- err
			}
		}(i, end)
		time.Sleep(500 * time.Millisecond)
	}

	// Batch creation for WechatNotifys
	for i := 0; i < len(WechatNotifys); i += batchSize {
		end := i + batchSize
		if end > len(WechatNotifys) {
			end = len(WechatNotifys)
		}
		// if i > 12000 {
		// 	adminLogger.Errorf("12000")
		// }
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			err := tx.CreateInBatches(WechatNotifys[start:end], len(WechatNotifys[start:end])).Error
			if err != nil {
				errCh <- err
			}
		}(i, end)
		time.Sleep(500 * time.Millisecond)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check for errors
	select {
	case err := <-errCh:
		adminLogger.Errorf("Batch creation failed, %s", err.Error())
		tx.Rollback()
	default:
		// No errors, commit the transaction
		adminLogger.Info("refoundSize: %d, paySize: %d, payIds: %v", refoundSize, paySize, payIds)
		tx.Commit()
	}

	return logWxNotify, err
}

func (e *OrderReconciliation) ScanWx() (logWxNotify []*models.LogWxScanNotify, err error) {
	if err = e.OrderDB.Find(&logWxNotify).Error; err != nil {
		return nil, err
	}

	for _, l := range logWxNotify {
		l.ResultMap, err = l.ParseForResult()
		if err != nil {
			continue
		}
		l.ParmMap, err = l.ParseParm()
		if err != nil {
			continue
		}
		l.DataMap, err = l.ParseForData()
		if err != nil {
			continue
		}
		if returnCode, ok := l.ResultMap["return_code"]; ok && returnCode.(string) == "SUCCESS" {
			if resultCode, ok1 := l.ResultMap["result_code"]; ok1 && resultCode.(string) == "SUCCESS" {
				l.ResultSuccessMap = l.ResultMap
				for k, v := range l.ResultSuccessMap {
					if k == "time_end" {
						t, err := time.Parse("20060102150405", v.(string))
						if err != nil {
							fmt.Println("Failed to parse time:", err)
							continue
						}
						l.ResultSuccessMap[k] = t
					}
				}
				for _, _ = range l.ParmMap {
					l.ParmMap["time"] = l.Time
				}
				for _, _ = range l.DataMap {
					l.DataMap["time"] = l.Time
				}
			}
		}
	}

	var wg sync.WaitGroup
	tx := e.OrderDB.Session(&gorm.Session{})
	tx.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	wg.Add(len(logWxNotify))
	for _, l := range logWxNotify {
		go func(l *models.LogWxScanNotify) {
			defer wg.Done()
			if l.ResultSuccessMap != nil {
				var result models.WxScanResult
				if err := mapstructure.WeakDecode(l.ResultSuccessMap, &result); err != nil {
					fmt.Printf("Decode ResultSuccessMap failed: %v\n", err)
					return
				}
				if err := tx.Create(&result).Error; err != nil {
					fmt.Printf("Create models.WxScanResult failed: %v\n", err)
					return
				}
				var parm models.WxScanParm
				if err := mapstructure.WeakDecode(l.ParmMap, &parm); err != nil {
					fmt.Printf("Decode WxScanParm failed: %v\n", err)
					return
				}
				parm.AuthCode = l.ParmMap["AuthCode"].(string)
				parm.OrderID = l.ParmMap["OrderId"].(string)
				parm.PayType = int8(l.ParmMap["PayType"].(float64))
				parm.UserID = l.ParmMap["UserID"].(string)
				if err := tx.Create(&parm).Error; err != nil {
					fmt.Printf("Create models.WxScanParm failed: %v\n", err)
					return
				}
				var data models.WxScan
				if err := mapstructure.WeakDecode(l.DataMap, &data); err != nil {
					fmt.Printf("Decode WxScanData failed: %v\n", err)
					return
				}
				if err := tx.Create(&data).Error; err != nil {
					fmt.Printf("Create models.WxScanData failed: %v\n", err)
					return
				}
			}
		}(l)
		time.Sleep(100 * time.Millisecond) // 避免 goroutine 并行执行导致数据插入顺序不一致
	}
	wg.Wait()

	return logWxNotify, nil
}

func (e *OrderReconciliation) InitOrderDB() error {
	dsn := "root:lg87516@tcp(127.0.0.1:3306)/rhhy_recon?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}
	//db.AutoMigrate(&models.LogZfbNotify{})
	e.OrderDB = db
	return nil
}

func SaveRhOrder(message storage.Messager) (err error) {
	//准备db
	// db := sdk.Runtime.GetDbByKey(message.GetPrefix())
	// if db == nil {
	// 	err = errors.New("db not exist")
	// 	log.Errorf("host[%s]'s %s", message.GetPrefix(), err.Error())
	// 	// Log writing to the database ignores error
	// 	return nil
	// }
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		adminLogger.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var l dto.OrderReconciliationGetPageReq
	err = json.Unmarshal(rb, &l)
	if err != nil {
		adminLogger.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	s := OrderReconciliation{}
	s.InitDB()
	_, err = s.ScanWx() //微信扫码支付订单
	if err != nil {
		adminLogger.Errorf("test error, %s", err.Error())
		return err
	}
	// _, err = s.MergeDataWxQueue() //微信支付回调订单
	// if err != nil {
	// 	adminLogger.Errorf("test error, %s", err.Error())
	// 	return err
	// }

	return nil
}
