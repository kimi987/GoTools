package logintoken

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	mrand "math/rand"
	"time"
)

var (
	ErrWrongCipherSize = errors.New("wrong cipher size")
	ZeroToken          = LoginToken{}
)

const (
	UNENCRYPTED_SIZE = 48
	ENCRYPTED_SIZE   = 64
	TokenExpireTime  = 48 * time.Hour
)

type SecureLoginTokenEncryptor struct {
	key cipher.Block
}

func NewSecureLoginToken(key []byte) (*SecureLoginTokenEncryptor, error) {
	aesKey, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "构建aes key失败")
	}

	return &SecureLoginTokenEncryptor{aesKey}, nil
}

func (s *SecureLoginTokenEncryptor) Encrypt(dst []byte, src []byte) {
	iv := dst[:16]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		binary.LittleEndian.PutUint64(iv, mrand.Uint64())
		binary.LittleEndian.PutUint64(iv[8:], mrand.Uint64())
	}
	stream := cipher.NewCFBEncrypter(s.key, iv)

	stream.XORKeyStream(dst[16:], src[:48])
}

func (s *SecureLoginTokenEncryptor) EncryptToken(t LoginToken) []byte {
	plain := make([]byte, 64)
	t.Marshal(plain[16:])
	s.Encrypt(plain, plain[16:])
	return plain
}

// 传入的数组内容会被修改, 不再是原内容
func (s *SecureLoginTokenEncryptor) Decrypt(b []byte) (LoginToken, error) {
	if len(b) != ENCRYPTED_SIZE {
		return ZeroToken, ErrWrongCipherSize
	}
	iv := b[:16]
	b = b[16:64]

	stream := cipher.NewCFBDecrypter(s.key, iv)

	stream.XORKeyStream(b, b)

	return Unmarshal(b), nil
}

// 总共48个byte
type LoginToken struct {
	UserSelfID         int64
	GameServerID       uint32
	GameServerRandomID uint32 // 重连时游戏服用来判断是不是同一个session
	GameServerAddr     [4]byte
	GameServerPort     uint32
	GenerateTime       time.Time
	SessionKey         uint64
	Reserved           uint64
}

func (t LoginToken) Marshal(result []byte) {
	binary.LittleEndian.PutUint64(result[0:], uint64(t.UserSelfID))
	binary.LittleEndian.PutUint32(result[8:], t.GameServerID)
	binary.LittleEndian.PutUint32(result[12:], t.GameServerRandomID)
	copy(result[16:], t.GameServerAddr[:])
	binary.LittleEndian.PutUint32(result[20:], t.GameServerPort)
	binary.LittleEndian.PutUint64(result[24:], uint64(t.GenerateTime.Unix()))
	binary.LittleEndian.PutUint64(result[32:], t.SessionKey)
	binary.LittleEndian.PutUint64(result[40:], t.Reserved)
}

func Unmarshal(b []byte) LoginToken {
	result := &LoginToken{}
	result.UserSelfID = int64(binary.LittleEndian.Uint64(b[0:]))
	result.GameServerID = binary.LittleEndian.Uint32(b[8:])
	result.GameServerRandomID = binary.LittleEndian.Uint32(b[12:])
	copy(result.GameServerAddr[:], b[16:20])
	result.GameServerPort = binary.LittleEndian.Uint32(b[20:])
	result.GenerateTime = time.Unix(int64(binary.LittleEndian.Uint64(b[24:])), 0)
	result.SessionKey = binary.LittleEndian.Uint64(b[32:])
	result.Reserved = binary.LittleEndian.Uint64(b[40:])

	return *result
}
