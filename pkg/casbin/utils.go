package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2/model"
)

func InitEmbeddedModel() (model.Model, error) {
	data, err := ReadEmbeddedModel()
	if err != nil {
		return nil, fmt.Errorf("read from embedded: %w", err)
	}

	rbacModel, err := model.NewModelFromString(string(data))
	if err != nil {
		return nil, fmt.Errorf("create model from string: %w", err)
	}

	return rbacModel, nil
}

func ReadEmbeddedModel() ([]byte, error) {
	data, err := embeddedModel.ReadFile("rbac_model.conf")
	if err != nil {
		return nil, fmt.Errorf("read model: %w", err)
	}

	return data, nil
}
