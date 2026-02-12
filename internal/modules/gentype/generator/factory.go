package gen

import (
	"fmt"
	"strings"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/gentype/generator/csharp"
	"github.com/khanalsaroj/typegen-server/internal/modules/gentype/generator/java"
	"github.com/khanalsaroj/typegen-server/internal/modules/gentype/generator/typescript"
)

func NewGenerator(req domain.TypeRequest) (Generator, error) {
	lang := strings.ToLower(string(req.TargetLanguage))
	style := strings.ToLower(req.Style)

	switch lang {
	case "java":
		switch style {
		case "dto":
			return &java.Dto{}, nil
		case "record":
			return &java.Record{}, nil
		default:
			return nil, fmt.Errorf("unsupported java type: %s", req.Style)
		}
	case "typescript", "ts":
		switch style {
		case "class", "interface", "type":
			return &typescript.Dto{}, nil
		case "zod":
			return &typescript.Zod{}, nil
		default:
			return nil, fmt.Errorf("unsupported typescript type: %s", req.Style)
		}
	case "csharp":
		switch style {
		case "dto":
			return &csharp.Dto{}, nil
		case "record":
			return &csharp.Record{}, nil
		default:
			return nil, fmt.Errorf("unsupported csharp type: %s", req.Style)
		}
	}
	return nil, fmt.Errorf("unsupported language type: %s", req.Style)
}
