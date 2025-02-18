package casbin

import (
	"embed"
	"fmt"
	"github.com/tecnologer/tempura/pkg/models/notificationtype"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/tecnologer/tempura/pkg/models/notification"
	"gorm.io/gorm"
)

//go:embed rbac_model.conf
var embeddedModel embed.FS

type Enforcer struct {
	*casbin.Enforcer
}

type Request struct {
	ChatID   int64
	Type     notificationtype.Type
	Value    float64
	LastSent time.Time
}

func (r *Request) PrepareToEnforce() []any {
	requestType := r.Type.String()
	if r.Type == notificationtype.All {
		requestType = "*"
	}

	return []any{strconv.FormatInt(r.ChatID, 10), requestType, r.Value, r.Value, r.LastSent}
}

func NewEnforcer(adapter persist.Adapter) (*Enforcer, error) {
	rbacModel, err := InitEmbeddedModel()
	if err != nil {
		return nil, fmt.Errorf("create model: %w", err)
	}

	params := []any{rbacModel}

	if adapter != nil {
		params = append(params, adapter)
	}

	enforcer, err := casbin.NewEnforcer(params...)
	if err != nil {
		return nil, fmt.Errorf("create enforcer: %w", err)
	}

	enforcer.AddFunction("Since", SinceFunc)
	enforcer.AddFunction("InRange", InRangeFunc)

	return &Enforcer{Enforcer: enforcer}, nil
}

func NewEnforcerWithDb(db *gorm.DB) (*Enforcer, error) {
	gormAdapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("create gorm adapter: %w", err)
	}

	enforcer, err := NewEnforcer(gormAdapter)
	if err != nil {
		return nil, fmt.Errorf("enforcer gorm adapter: %w", err)
	}

	return enforcer, nil
}

func (e *Enforcer) Enforce(rules ...*Request) (bool, error) {
	var (
		isAuthorized bool
		err          error
	)

	for _, rule := range rules {
		isAuthorized, err = e.Enforcer.Enforce(rule.PrepareToEnforce()...)
		if err != nil {
			return false, fmt.Errorf("enforce: %w", err)
		}

		if isAuthorized {
			return true, nil
		}
	}

	return isAuthorized, nil
}
