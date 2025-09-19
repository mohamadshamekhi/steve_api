package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RepairStatus string

const (
	StatusNew             RepairStatus = "new"              // جدید (ثبت شده توسط کاربر)
	StatusPendingReview   RepairStatus = "pending_review"   // در انتظار بررسی کارشناس
	StatusPendingCustomer RepairStatus = "pending_customer" // منتظر پاسخ مشتری (مثلاً برای تایید هزینه)
	StatusInProgress      RepairStatus = "in_progress"      // در حال تعمیر
	StatusCompleted       RepairStatus = "completed"        // تعمیر کامل شده
	StatusShippedBack     RepairStatus = "shipped_back"     // برای مشتری ارسال شده
	StatusCanceled        RepairStatus = "canceled"         // لغو شده (توسط کاربر یا ادمین)
)

type Repair struct {
	ID          string          `bson:"_id"`
	UserID      string          `bson:"user_id"`
	BrandID     string          `bson:"brand_id"`
	ProductID   string          `bson:"product_id"`
	Description string          `bson:"description"`
	Status      RepairStatus    `bson:"status"`
	Galleries   []RepairGallery `bson:"galleries"`
	Price       uint64          `bson:"price"`
	CreatedAt   time.Time       `bson:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at"`
}

type RepairGallery struct {
	MediaID string `bson:"media_id" json:"media_id"`
	Url     string `bson:"url" json:"url"`
}

func (r *Repair) BeforeCreate() {
	r.ID = ulid.Make().String()
	r.CreatedAt = time.Now()
	r.UpdatedAt = r.CreatedAt
}
