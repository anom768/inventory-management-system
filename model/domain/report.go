package domain

import "time"

type Activities struct {
	ID             int       `gorm:"primary_key;column:id;auto_increment" json:"id"`
	ItemID         int       `gorm:"column:item_id" json:"item_id"`
	Action         string    `gorm:"column:action" json:"action"`
	QuantityChange int       `gorm:"column:quantity_change" json:"quantity_change"`
	Timestamp      time.Time `gorm:"column:timestamp" json:"timestamp"`
	PerformedBy    string    `gorm:"column:performed_by" json:"performed_by"`
}

type ReportStock struct {
	TotalItems          int     `json:"total_items"`
	TotalQuantity       int     `json:"total_quantity"`
	TotalInventoryValue float64 `json:"total_inventory_value"`
	Items               []Items `json:"items"`
}
