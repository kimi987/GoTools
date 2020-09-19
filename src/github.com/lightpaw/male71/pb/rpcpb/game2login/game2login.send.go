package game2login

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func VerifyLoginToken(c client, ctx context.Context, hero_id int64, token string, client_ip string, pf uint32) (*S2CVerifyLoginTokenProto, error) {
	return VerifyLoginTokenProto(c, ctx, &C2SVerifyLoginTokenProto{

		HeroId: hero_id,

		Token: token,

		ClientIp: client_ip,

		Pf: pf,
	})
}

func VerifyLoginTokenProto(c client, ctx context.Context, proto *C2SVerifyLoginTokenProto) (*S2CVerifyLoginTokenProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "verify_login_token proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "verify_login_token", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "verify_login_token fail")
		}

		s2c := &S2CVerifyLoginTokenProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "verify_login_token s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func WriteTlog(c client, ctx context.Context, hero_id int64) (*S2CWriteTlogProto, error) {
	return WriteTlogProto(c, ctx, &C2SWriteTlogProto{

		HeroId: hero_id,
	})
}

func WriteTlogProto(c client, ctx context.Context, proto *C2SWriteTlogProto) (*S2CWriteTlogProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "write_tlog proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "write_tlog", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "write_tlog fail")
		}

		s2c := &S2CWriteTlogProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "write_tlog s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func Push(c client, ctx context.Context, hero_id int64, sid uint32, title string, content string, start_time int64, expire_time int64) (*S2CPushProto, error) {
	return PushProto(c, ctx, &C2SPushProto{

		HeroId: hero_id,

		Sid: sid,

		Title: title,

		Content: content,

		StartTime: start_time,

		ExpireTime: expire_time,
	})
}

func PushProto(c client, ctx context.Context, proto *C2SPushProto) (*S2CPushProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "push proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "push", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "push fail")
		}

		s2c := &S2CPushProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "push s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func PushMulti(c client, ctx context.Context, hero_ids []int64, sid uint32, title string, content string, start_time int64, expire_time int64) (*S2CPushMultiProto, error) {
	return PushMultiProto(c, ctx, &C2SPushMultiProto{

		HeroIds: hero_ids,

		Sid: sid,

		Title: title,

		Content: content,

		StartTime: start_time,

		ExpireTime: expire_time,
	})
}

func PushMultiProto(c client, ctx context.Context, proto *C2SPushMultiProto) (*S2CPushMultiProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "push_multi proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "push_multi", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "push_multi fail")
		}

		s2c := &S2CPushMultiProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "push_multi s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}
