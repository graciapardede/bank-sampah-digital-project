package models

import "gorm.io/gorm"

type RedemptionItem struct {
    gorm.Model
    RedemptionID  uint        `json:"redemption_id"`
    Redemption    *Redemption `gorm:"foreignKey:RedemptionID" json:"redemption,omitempty"`
    RewardItemID  uint        `json:"reward_item_id"`
    RewardItem    *RewardItem `gorm:"foreignKey:RewardItemID" json:"reward_item,omitempty"`
    Quantity      int         `json:"quantity"`
}
