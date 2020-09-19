package login2game

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func CompleteQuestionnaire(c client, ctx context.Context, hero_id int64, qnid string) (*S2CCompleteQuestionnaireProto, error) {
	return CompleteQuestionnaireProto(c, ctx, &C2SCompleteQuestionnaireProto{

		HeroId: hero_id,

		Qnid: qnid,
	})
}

func CompleteQuestionnaireProto(c client, ctx context.Context, proto *C2SCompleteQuestionnaireProto) (*S2CCompleteQuestionnaireProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "complete_questionnaire proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "complete_questionnaire", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "complete_questionnaire fail")
		}

		s2c := &S2CCompleteQuestionnaireProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "complete_questionnaire s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func BuyProduct(c client, ctx context.Context, order_id string, order_amount uint64, order_time int64, pid uint32, sid uint32, hero_id int64, product_id uint64) (*S2CBuyProductProto, error) {
	return BuyProductProto(c, ctx, &C2SBuyProductProto{

		OrderId: order_id,

		OrderAmount: order_amount,

		OrderTime: order_time,

		Pid: pid,

		Sid: sid,

		HeroId: hero_id,

		ProductId: product_id,
	})
}

func BuyProductProto(c client, ctx context.Context, proto *C2SBuyProductProto) (*S2CBuyProductProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "buy_product proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "buy_product", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "buy_product fail")
		}

		s2c := &S2CBuyProductProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "buy_product s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}
