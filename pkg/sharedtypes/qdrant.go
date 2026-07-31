package sharedtypes

import (
	"encoding/json"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

const (
	qdrantQueryContext        = "context"
	qdrantQueryDiscover       = "discover"
	qdrantQueryFormula        = "formula"
	qdrantQueryFusion         = "fusion"
	qdrantQueryNearest        = "nearest"
	qdrantQueryNearestWithMmr = "nearestWithMmr"
	qdrantQueryOrderBy        = "orderBy"
	qdrantQueryRecommend      = "recommend"
	qdrantQueryRrf            = "rrf"
	qdrantQuerySample         = "sample"
)

type QdrantQuery struct {
	*qdrant.Query
}

func (q QdrantQuery) MarshalJSON() ([]byte, error) {
	var variant string
	switch q.Variant.(type) {
	case *qdrant.Query_Context:
		variant = qdrantQueryContext
	case *qdrant.Query_Discover:
		variant = qdrantQueryDiscover
	case *qdrant.Query_Formula:
		variant = qdrantQueryFormula
	case *qdrant.Query_Fusion:
		variant = qdrantQueryFusion
	case *qdrant.Query_Nearest:
		variant = qdrantQueryNearest
	case *qdrant.Query_NearestWithMmr:
		variant = qdrantQueryNearestWithMmr
	case *qdrant.Query_OrderBy:
		variant = qdrantQueryOrderBy
	case *qdrant.Query_Recommend:
		variant = qdrantQueryRecommend
	case *qdrant.Query_Rrf:
		variant = qdrantQueryRrf
	case *qdrant.Query_Sample:
		variant = qdrantQuerySample
	default:
		return nil, fmt.Errorf("unexpected qdrant.isQuery_Variant: %#v", q.Variant)
	}

	innerBytes, err := json.Marshal(q.Variant)
	if err != nil {
		return nil, err
	}
	var inner json.RawMessage
	if err = json.Unmarshal(innerBytes, &inner); err != nil {
		return nil, err
	}
	return json.Marshal(enumTaggedJson{variant, inner})

}

func (q QdrantQuery) UnmarshalJSON(data []byte) error {
	var intermediate enumTaggedJson
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}
	switch intermediate.Variant {
	case qdrantQueryContext:
		var variant *qdrant.Query_Context
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryDiscover:
		var variant *qdrant.Query_Discover
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryFormula:
		var variant *qdrant.Query_Formula
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryFusion:
		var variant *qdrant.Query_Fusion
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryNearest:
		var variant *qdrant.Query_Nearest
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryNearestWithMmr:
		var variant *qdrant.Query_NearestWithMmr
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryOrderBy:
		var variant *qdrant.Query_OrderBy
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryRecommend:
		var variant *qdrant.Query_Recommend
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQueryRrf:
		var variant *qdrant.Query_Rrf
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case qdrantQuerySample:
		var variant *qdrant.Query_Sample
		if err := json.Unmarshal(intermediate.Inner, &variant); err != nil {
			return err
		}
		q.Variant = variant
	case "":
		// on initial launch of the workflow a variable of this type will have a nil variant
	default:
		return fmt.Errorf("unknow query variant %q", intermediate.Variant)
	}
	return nil
}

type QdrantCondition struct {
	*qdrant.Condition
}

type QdrantFilter struct {
	*qdrant.Filter
}

type QdrantScoredPoint struct {
	*qdrant.ScoredPoint
}

type enumTaggedJson struct {
	Variant string          `json:"variant"`
	Inner   json.RawMessage `json:"inner"`
}
