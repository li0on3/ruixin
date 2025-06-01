package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusDoing     OrderStatus = 1
	OrderStatusSuccess   OrderStatus = 2
	OrderStatusFailed    OrderStatus = 3
	OrderStatusRefunded  OrderStatus = 4
	OrderStatusCancelled OrderStatus = 5
)

type Order struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 订单基础信息
	OrderNo       string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_no"`
	PutOrderID    string      `gorm:"type:varchar(50);index" json:"put_order_id"`
	DistributorID uint        `json:"distributor_id"`
	CardID        uint        `json:"card_id"`
	CardCode      string      `gorm:"type:varchar(100)" json:"card_code"`
	Status        OrderStatus `json:"status"`

	// 门店信息
	StoreCode    string `gorm:"type:varchar(50)" json:"store_code"`
	StoreName    string `gorm:"type:varchar(200)" json:"store_name"`
	StoreAddress string `gorm:"type:varchar(500)" json:"store_address"`

	// 商品信息
	Goods        OrderGoods `gorm:"type:json" json:"goods"`
	TotalAmount  float64    `json:"total_amount"`
	CostAmount   float64    `json:"cost_amount"`
	ProfitAmount float64    `json:"profit_amount"`
	
	// 瑞幸原始价格（用于对账）
	LuckinPrice     float64 `json:"luckin_price"`      // 瑞幸原始销售价
	LuckinCostPrice float64 `json:"luckin_cost_price"` // 瑞幸原始成本价

	// 取餐信息
	TakeMode    int    `json:"take_mode"` // 1: 到店自取
	TakeCode    string `gorm:"type:varchar(50)" json:"take_code"`
	QRData      string `gorm:"type:text" json:"qr_data"`
	PhoneNumber string `gorm:"type:varchar(20)" json:"phone_number"`

	// 回调信息
	CallbackURL    string     `gorm:"type:varchar(500)" json:"callback_url"`
	CallbackStatus int        `json:"callback_status"`
	CallbackTime   *time.Time `json:"callback_time"`

	// 瑞幸平台响应
	LuckinResponse string `gorm:"type:text" json:"-"`
}

type OrderGoods []OrderGoodsItem

type OrderGoodsItem struct {
	GoodsID       string  `json:"goods_id"`
	GoodsName     string  `json:"goods_name"`
	GoodsImage    string  `json:"goods_image"`
	SKUCode       string  `json:"sku_code"`
	SKUName       string  `json:"sku_name"`
	Quantity      int     `json:"quantity"`
	OriginalPrice float64 `json:"original_price"`
	SalePrice     float64 `json:"sale_price"`
	Specs         []struct {
		SpecsCode string `json:"specs_code"`
		SpecsName string `json:"specs_name"`
		Value     string `json:"value"`
	} `json:"specs"`
}

func (g OrderGoods) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *OrderGoods) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, g)
}
