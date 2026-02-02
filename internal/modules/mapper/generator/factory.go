package generator

import (
	"fmt"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/mapper/generator/java"
	"strings"
)

func NewGenerator(req domain.MapperRequest) (Mapper, error) {
	tgtTyp := strings.ToLower(string(req.TargetType))
	switch tgtTyp {
	case "mybatis-xml":
		return &java.Xml{}, nil
	case "mybatis-annotation":
		return &java.XmlAnnotation{}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", req.TargetType)
	}
}
