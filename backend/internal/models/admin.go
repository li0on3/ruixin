package models

import (
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username    string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Password    string    `gorm:"type:varchar(255)" json:"-"`
	Email       string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone       string    `gorm:"type:varchar(20)" json:"phone"`
	RealName    string    `gorm:"type:varchar(50)" json:"real_name"`
	Role        string    `gorm:"type:varchar(50)" json:"role"` // super_admin, admin, operator
	Status      int       `json:"status"`                       // 0: 禁用, 1: 正常
	LastLoginAt time.Time `json:"last_login_at"`
	LastLoginIP string    `gorm:"type:varchar(45)" json:"last_login_ip"`
}

type AdminOperationLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	AdminID   uint      `json:"admin_id"`
	AdminName string    `json:"admin_name"`
	Operation string    `json:"operation"`
	Module    string    `json:"module"`
	Details   string    `gorm:"type:text" json:"details"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}
