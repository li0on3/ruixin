package httpclient

// CheckGoodsCardResponse 检查优惠卡响应
type CheckGoodsCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   struct {
		Status  int    `json:"status"`
		OrderNo string `json:"orderNo"`
	} `json:"data"`
}

// CityByCardRequest 获取城市列表请求
type CityByCardRequest struct {
	ProductID int    `json:"productId"`
	Card      string `json:"card"`
}

// CityByCardResponse 获取城市列表响应
type CityByCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   []City `json:"data"`
}

type City struct {
	CityID     int    `json:"cityId"`
	CityName   string `json:"cityName"`
	CityPinyin string `json:"cityPinyin"`
}

// StoreByCardRequest 搜索门店请求
type StoreByCardRequest struct {
	Card      string `json:"card"`
	ProductID int    `json:"productId"`
	CityID    int    `json:"cityId"`
	Keywords  string `json:"keywords"`
}

// StoreByCardResponse 搜索门店响应
type StoreByCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   []Store `json:"data"`
}

type Store struct {
	StoreCode      string `json:"storeCode"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	City           string `json:"city"`
	CityID         int    `json:"cityId"`
	Lat            string `json:"lat"`
	Lng            string `json:"lng"`
	StartTime      string `json:"starttime"`
	EndTime        string `json:"endtime"`
	DistanceKM     string `json:"distancekm"`
	OfficialStatus bool   `json:"officialStatus"`
	CloseNote      string `json:"closeNote"`
}

// MenuByCardRequest 获取菜单请求
type MenuByCardRequest struct {
	ProductID       int    `json:"productId"`
	OrderType       int    `json:"orderType"`
	StoreCode       string `json:"storeCode"`
	UpDiscountRate  string `json:"upDiscountRate"`
	Card            string `json:"card"`
}

// MenuByCardResponse 获取菜单响应
type MenuByCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   struct {
		StoreCode      string `json:"storeCode"`
		StoreName      string `json:"storeName"`
		StoreAddress   string `json:"storeAddress"`
		StartTime      string `json:"starttime"`
		EndTime        string `json:"endtime"`
		OfficialStatus bool   `json:"officialSataus"`
		Notice         string `json:"notice"`
		Menu           []MenuCategory `json:"menu"`
	} `json:"data"`
}

type MenuCategory struct {
	CategoryCode  string `json:"categoryCode"`
	CategoryName  string `json:"categoryName"`
	CategoryImage string `json:"categoryImage"`
	GoodsList     []Good `json:"goodsList"`
}

type Good struct {
	GoodsCode   string `json:"goodsCode"`
	GoodsName   string `json:"goodsName"`
	GoodsImage  string `json:"goodsImage"`
	CostPrice   string `json:"costPrice"`
	LinePrice   string `json:"linePrice"`
	SaleStatus  int    `json:"saleStatus"`
	IsChoice    int    `json:"isChoice"`
}

// GoodsByCardRequest 获取商品详情请求
type GoodsByCardRequest struct {
	ProductID      int    `json:"productId"`
	OrderType      int    `json:"orderType"`
	StoreCode      string `json:"storeCode"`
	LinkID         string `json:"linkId"`
	UpDiscountRate string `json:"upDiscountRate"`
	Card           string `json:"card"`
}

// GoodsByCardResponse 获取商品详情响应
type GoodsByCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   struct {
		GoodsCode  string      `json:"goodsCode"`
		GoodsName  string      `json:"goodsName"`
		GoodsImage string      `json:"goodsImage"`
		CostPrice  string      `json:"costPrice"`
		LinePrice  string      `json:"linePrice"`
		GoodsDesc  string      `json:"goodsDesc"`
		GoodsSkus  []GoodsSku  `json:"goodsSkus"`
		GoodsSpecs []GoodsSpec `json:"goodsSpecs"`
	} `json:"data"`
}

type GoodsSku struct {
	SkuCode     string `json:"skuCode"`
	SkuName     string `json:"skuName"`
	SkuShowName string `json:"skuShowName"`
	CostPrice   string `json:"costPrice"`
	LinePrice   string `json:"linePrice"`
}

type GoodsSpec struct {
	SpecsCode       string     `json:"specsCode"`
	SpecsName       string     `json:"specsName"`
	SpecsType       int        `json:"specsType"`
	SpecsChoices    int        `json:"specsChoices"`
	SpecsMaxChoices int        `json:"specsMaxChoices"`
	SpecsMinChoices int        `json:"specsMinChoices"`
	SpecsItems      []SpecItem `json:"specsItems"`
}

type SpecItem struct {
	Code      string `json:"code"`
	AddPrice  string `json:"addPrice"`
	Amount    string `json:"amount"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	IsDefault int    `json:"isDefault"`
	CostPrice string `json:"costPrice"`
	LinePrice string `json:"linePrice"`
}

// CheckByCardRequest 检查订单请求
type CheckByCardRequest struct {
	StoreCode      string `json:"storeCode"`
	ProductID      int    `json:"productId"`
	Goods          string `json:"goods"`
	OrderNo        string `json:"orderNo"`
	TakeMode       int    `json:"takeMode"`
	UpDiscountRate string `json:"upDiscountRate"`
	Card           string `json:"card"`
}

// CheckByCardResponse 检查订单响应
type CheckByCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   struct {
		Goods                   []CheckGoods `json:"goods"`
		CityID                  int          `json:"cityId"`
		City                    string       `json:"city"`
		StoreCode               string       `json:"storeCode"`
		Name                    string       `json:"name"`
		Address                 string       `json:"address"`
		StartTime               string       `json:"starttime"`
		EndTime                 string       `json:"endtime"`
		OfficialStatus          string       `json:"officialSataus"`
		DeliveryPrice           string       `json:"deliveryPrice"`
		TotalOriginalPrice      string       `json:"totalOriginalPrice"`
		TotalSalePrice          string       `json:"totalSalePrice"`
		TotalSavePrice          string       `json:"totalSavePrice"`
		TotalPlatformCostPrice  string       `json:"totalPlatformCostPrice"`
		OrderNo                 string       `json:"orderNo"`
		Introduce               string       `json:"introduce"`
	} `json:"data"`
}

type CheckGoods struct {
	GoodsID            string `json:"goodsId"`
	GoodsName          string `json:"goodsName"`
	GoodsImg           string `json:"goodsImg"`
	Num                int    `json:"num"`
	TotalOriginalPrice string `json:"totalOriginalPrice"`
	TotalSalePrice     string `json:"totalSalePrice"`
	PutDetailJSON      string `json:"putDetailJson"`
	GoodsOutID         string `json:"goodsOutId"`
	SalePrice          string `json:"salePrice"`
	PlatformCostPrice  string `json:"platformCostPrice"`
}

// OrderByCardRequest 下单请求
type OrderByCardRequest struct {
	ProductID      int    `json:"productId"`
	OrderNo        string `json:"orderNo"`
	PhoneNo        string `json:"phoneNo"`
	TakeMode       string `json:"takeMode"`
	CallbackURL    string `json:"callBackUrl"`
	UpDiscountRate string `json:"upDiscountRate"`
	Card           string `json:"card"`
}

// OrderByCardResponse 下单响应
type OrderByCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   struct {
		Biz                    string        `json:"biz"`
		PutOrderID             string        `json:"putOrderId"`
		OutID                  string        `json:"outId"`
		Status                 string        `json:"status"`
		CallbackTag            string        `json:"callBackTag"`
		TackInfo               string        `json:"tackinfo"`
		TackCode               string        `json:"tackCode"`
		SalePrice              string        `json:"salePrice"`
		TotalPlatformCostPrice string        `json:"totalPlatformCostPrice"`
		RefundAmount           string        `json:"refundAmount"`
		StoreOrderID           string        `json:"storeOrderId"`
		Note                   string        `json:"note"`
		AddTime                string        `json:"addtime"`
		StoreCode              string        `json:"storeCode"`
		StoreName              string        `json:"storeName"`
		StoreAddress           string        `json:"storeAddress"`
		CourierLocationH5      string        `json:"courierLocationH5"`
		DeliveryPrice          string        `json:"deliveryPrice"`
		RefundAll              []interface{} `json:"refundAll"`
		RefundAdd              []interface{} `json:"refundAdd"`
		Detail                 []OrderDetail `json:"detail"`
	} `json:"data"`
}

type OrderDetail struct {
	Biz                string `json:"biz"`
	GoodsID            string `json:"goodsId"`
	GoodsName          string `json:"goodsName"`
	Num                int    `json:"num"`
	TotalOriginalPrice string `json:"totalOriginalPrice"`
	TotalSalePrice     string `json:"totalSalePrice"`
	OriginalPrice      string `json:"originalPrice"`
	SalePrice          string `json:"salePrice"`
	RefundNum          int    `json:"refundNum"`
	RefundAmount       string `json:"refundAmount"`
	Memo               string `json:"memo"`
	PlatformCostPrice  string `json:"platformCostPrice"`
}

// QueryByCardRequest 查询订单请求
type QueryByCardRequest struct {
	Brand   string `json:"brand"`
	Card    string `json:"card"`
	OrderNo string `json:"orderNo"`
}

// QueryByCardResponse 查询订单响应
type QueryByCardResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Data   struct {
		Biz                    string         `json:"biz"`
		PutOrderID             string         `json:"putOrderId"`
		OutID                  string         `json:"outId"`
		Status                 string         `json:"status"`
		TakeCode               string         `json:"takeCode"`
		QRData                 string         `json:"qrdata"`
		AddTime                string         `json:"addtime"`
		StoreCode              string         `json:"storeCode"`
		StoreName              string         `json:"storeName"`
		SalePrice              string         `json:"salePrice"`
		CostPrice              string         `json:"costPrice"`
		RefundAmount           string         `json:"refundAmount"`
		TotalPlatformCostPrice string         `json:"totalPlatformCostPrice"`
		RefundAll              []interface{}  `json:"refundAll"`
		RefundAdd              []interface{}  `json:"refundAdd"`
		Goods                  []QueryGoods   `json:"goods"`
		TakeInfoList           []TakeInfo     `json:"takeInfoList"`
	} `json:"data"`
}

type QueryGoods struct {
	Biz                string `json:"biz"`
	GoodsID            string `json:"goodsId"`
	GoodsName          string `json:"goodsName"`
	Num                int    `json:"num"`
	TotalOriginalPrice string `json:"totalOriginalPrice"`
	TotalSalePrice     string `json:"totalSalePrice"`
	OriginalPrice      string `json:"originalPrice"`
	SalePrice          string `json:"salePrice"`
	RefundNum          int    `json:"refundNum"`
	RefundAmount       string `json:"refundAmount"`
	Memo               string `json:"memo"`
	PlatformCostPrice  string `json:"platformCostPrice"`
}

type TakeInfo struct {
	TakeInfo        string `json:"takeinfo"`
	QRData          string `json:"qrdata"`
	LockerCode      string `json:"lockerCode"`
	OrderPhone      string `json:"orderPhone"`
	GoodsID         string `json:"goodsId"`
	PutOrderDetailID int   `json:"putOrderDetailId"`
	DetailJSON      string `json:"detailJson"`
}