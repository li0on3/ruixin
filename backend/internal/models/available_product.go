package models

import "time"

// PublicProductInfo 公开的商品信息（不包含敏感数据）
type PublicProductInfo struct {
	GoodsCode     string                     `json:"goods_code"`
	GoodsName     string                     `json:"goods_name"`
	Category      string                     `json:"category"`
	IsAvailable   bool                       `json:"is_available"`
	StockStatus   string                     `json:"stock_status"` // high, medium, low, out
	PriceRange    string                     `json:"price_range"`
	SKUs          []PublicSKUInfo            `json:"skus"`
	Aliases       []string                   `json:"aliases"`
	SpecsOptions  map[string][]SpecOption    `json:"specs_options,omitempty"` // 所有可用的规格选项
	OrderingTips  string                     `json:"ordering_tips,omitempty"`
}

// PublicSKUInfo 公开的SKU信息
type PublicSKUInfo struct {
	SKUCode      string                 `json:"sku_code"`
	ChineseDesc  string                 `json:"chinese_desc"`
	Specs        map[string]SpecOption  `json:"specs"`
	SpecsCode    string                 `json:"specs_code"`
	IsDefault    bool                   `json:"is_default"`
	IsAvailable  bool                   `json:"is_available"`
}

// SpecOption 规格选项
type SpecOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// AvailableProductsResponse API响应结构
type AvailableProductsResponse struct {
	UpdatedAt time.Time           `json:"updated_at"`
	Products  []PublicProductInfo `json:"products"`
	Notice    string              `json:"notice,omitempty"`
}

// StockStatus 库存状态常量
const (
	StockStatusHigh   = "high"   // 库存充足 (>100)
	StockStatusMedium = "medium" // 库存一般 (20-100)
	StockStatusLow    = "low"    // 库存较少 (<20)
	StockStatusOut    = "out"    // 暂时缺货 (0)
)

// CalculateStockStatus 根据数量计算库存状态
func CalculateStockStatus(count int) string {
	if count > 100 {
		return StockStatusHigh
	} else if count >= 20 {
		return StockStatusMedium
	} else if count > 0 {
		return StockStatusLow
	}
	return StockStatusOut
}