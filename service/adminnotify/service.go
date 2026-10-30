// Package adminnotifyservice texts every super_admin when a customer submits
// something the shop should react to (contact form, comment, partner
// application, order cancellation). Sending is detached from the request: an
// SMS failure must never fail what the customer did.
package adminnotifyservice

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	"github.com/abolfazlnorzad/steve/entity"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	persian "github.com/yaa110/go-persian-calendar"
)

type Config struct {
	// Enabled turns these admin SMS on or off. The older repair / sell /
	// payment admin SMS are not affected.
	Enabled    bool          `koanf:"enabled"`
	SMSTimeout time.Duration `koanf:"sms_timeout"`
}

type AdminLister interface {
	GetSuperUserBackOffices(ctx context.Context, req backofficeparam.GetSuperUserBackofficeRequest) (backofficeparam.GetSuperUserBackofficeResponse, error)
}

type UserReader interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

type SMSSender interface {
	SendPatternSMS(ctx context.Context, to string, text []string, patternKey string) error
}

// Event is one customer action to report.
type Event struct {
	// Title names what happened, e.g. "پیام تماس با ما".
	Title string
	// UserID, when set, fills Name and Phone from the customer's account
	// (only where they are empty).
	UserID string
	Name   string
	Phone  string
	// Detail is a one-line summary (subject, comment text, order id …);
	// it is flattened and cut short to keep the SMS to one part.
	Detail string
}

type Service struct {
	cfg    Config
	admins AdminLister
	users  UserReader
	sms    SMSSender
}

func New(cfg Config, admins AdminLister, users UserReader, sms SMSSender) Service {
	return Service{cfg: cfg, admins: admins, users: users, sms: sms}
}

const detailMaxRunes = 60

// Notify sends the admin_notify SMS ({0}=title {1}=name {2}=phone {3}=detail
// {4}=Jalali date and time) to every super_admin, in the background.
func (s Service) Notify(e Event) {
	if !s.cfg.Enabled {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.L().Error("adminnotify_panic", "title", e.Title, "panic", r)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), s.timeout())
		defer cancel()

		if e.UserID != "" && (e.Name == "" || e.Phone == "") && s.users != nil {
			if u, err := s.users.GetUserByID(ctx, e.UserID); err == nil {
				if e.Name == "" {
					e.Name = u.Name
				}
				if e.Phone == "" {
					e.Phone = u.PhoneNumber
				}
			}
		}

		resp, err := s.admins.GetSuperUserBackOffices(ctx, backofficeparam.GetSuperUserBackofficeRequest{})
		if err != nil {
			logger.L().Error("adminnotify_admins_lookup_failed", "title", e.Title, "err", err)
			return
		}

		text := []string{e.Title, orDash(e.Name), orDash(e.Phone), orDash(shorten(e.Detail)), nowJalali()}
		for _, admin := range resp.Backoffice {
			if admin.PhoneNumber == "" {
				continue
			}
			if sErr := s.sms.SendPatternSMS(ctx, admin.PhoneNumber, text, mellipayamak.AdminNotifyKey); sErr != nil {
				logger.L().Error("adminnotify_sms_failed", "title", e.Title, "admin_id", admin.ID, "err", sErr)
			}
		}
	}()
}

func (s Service) timeout() time.Duration {
	if s.cfg.SMSTimeout <= 0 {
		return 40 * time.Second
	}
	return s.cfg.SMSTimeout
}

func nowJalali() string {
	now := time.Now()
	if loc, err := time.LoadLocation("Asia/Tehran"); err == nil {
		now = now.In(loc)
	}
	return persian.New(now).Format("yyyy/MM/dd HH:mm")
}

// shorten flattens whitespace (newlines would break the SMS layout) and cuts
// the text to detailMaxRunes characters.
func shorten(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if utf8.RuneCountInString(s) <= detailMaxRunes {
		return s
	}
	r := []rune(s)
	return string(r[:detailMaxRunes-1]) + "…"
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
