package model


type EaNotification struct {
	ID              *uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	UserID          uint16  `gorm:"column:user_id"`
	NotifyMsg       string  `gorm:"column:notify_msg"`
	Status          string  `gorm:"column:status"`
	RecordTimestamp string  `gorm:"column:record_timestamp"`
}



func (EaNotification) TableName() string {
	return "ea_notification_table"
}
