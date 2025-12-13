package metadata

import "github.com/shopspring/decimal"

type VisibleImpressionInfo struct {
	Amplitude Amplitude `json:"amplitude"`
}

type Amplitude struct {
	AllTimeQuantitySold            int     `json:"all_time_quantity_sold"`
	BrandName                      string  `json:"brand_name"`
	CategoryL1Name                 string  `json:"category_l1_name"`
	CategoryL2Name                 string  `json:"category_l2_name"`
	DeboostedHighPriceDiffPosition int     `json:"deboosted_high_price_diff_position"`
	DiscountedByRuleIDs            string  `json:"discounted_by_rule_ids"`
	EarliestDeliveryEstimate       int     `json:"earliest_delivery_estimate"`
	IsAuthentic                    int     `json:"is_authentic"`
	IsBestOfferAvailable           bool    `json:"is_best_offer_available"`
	IsFlashDeal                    bool    `json:"is_flash_deal"`
	IsFreeshipXtra                 bool    `json:"is_freeship_xtra"`
	IsGiftAvailable                bool    `json:"is_gift_available"`
	IsHero                         bool    `json:"is_hero"`
	IsHighPricePenalty             bool    `json:"is_high_price_penalty"`
	IsImported                     bool    `json:"is_imported"`
	IsTopBrand                     bool    `json:"is_top_brand"`
	JoinedStrategyRerank           bool    `json:"joined_strategy_rerank"`
	Layout                         string  `json:"layout"`
	MasterProductSKU               string  `json:"master_product_sku"`
	NumberOfReviews                int     `json:"number_of_reviews"`
	OrderRoute                     string  `json:"order_route"`
	Origin                         string  `json:"origin"`
	PartnerRewardsAmount           int     `json:"partner_rewards_amount"`
	Price                          int     `json:"price"`
	PrimaryCategoryName            string  `json:"primary_category_name"`
	ProductRating                  float64 `json:"product_rating"`
	SearchRank                     int     `json:"search_rank"`
	SellerID                       int     `json:"seller_id"`
	SellerProductID                int     `json:"seller_product_id"`
	SellerProductSKU               string  `json:"seller_product_sku"`
	SellerType                     string  `json:"seller_type"`
	StandardDeliveryEstimate       float64 `json:"standard_delivery_estimate"`
	TikiVerified                   int     `json:"tiki_verified"`
	TikinowDeliveryEstimate        int     `json:"tikinow_delivery_estimate"`
	TikiproDeliveryEstimate        int     `json:"tikipro_delivery_estimate"`
	Variant                        bool    `json:"variant"`
}

type Product struct {
	Id              int64           `json:"id"`
	Sku             string          `json:"sku"`
	Name            string          `json:"name"`
	UrlKey          string          `json:"url_key"`
	UrlPath         string          `json:"url_path"`
	BrandName       string          `json:"brand_name"`
	Price           decimal.Decimal `json:"price"`
	Discount        decimal.Decimal `json:"discount"`
	DiscountRate    float64         `json:"discount_rate"`
	RatingAverage   float64         `json:"rating_average"`
	ReviewCount     int32           `json:"review_count"`
	OrderCount      int32           `json:"order_count"`
	FavoriteCount   int32           `json:"favorite_count"`
	ThumbnailUrl    string          `json:"thumbnail_url"`
	FreegiftItems   []Product       `json:"freegift_items"`
	InventoryStatus string          `json:"inventory_status"`
	quantitySold    struct {
		Text  string `json:"text"`
		Value string `json:"string"`
	}
	OriginalPrice         decimal.Decimal       `json:"original_price"`
	Availability          int32                 `json:"availability"`
	PrimaryCategoryPath   string                `json:"primary_category_path"`
	ProductRecoScore      int32                 `json:"product_reco_score"`
	SellerId              int64                 `json:"seller_id"`
	SellerProductId       int64                 `json:"seller_product_id"`
	VisibleImpressionInfo VisibleImpressionInfo `json:"visible_impression_info"`
}
