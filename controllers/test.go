// utils project utils.go
package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

var I_key string = "ebcd1234efgh5678ebcd1234efgh5678" //原始值
var E_key string = "5678dcba1234abcd5678dcba1234abcd" //加密过的
var Iv []byte = []byte{0x58, 0xbd, 0x32, 0x78, 0x0, 0x1, 0x61, 0x24, 0x58, 0xbd, 0x32, 0x78, 0x0, 0x1, 0x61, 0x24}

//解密当前的金额
//返回金额是零表示失败
func Decprice(base64price string) uint64 {
	defer func() uint64 {
		if err := recover(); err != nil {
			return 0
		}
		return 0
	}()
	Lstr := base64price

	for {
		if len(Lstr)%4 == 0 {
			break
		}
		Lstr += "="
	}
	base64.RawStdEncoding.DecodeString(Lstr)
	enc_price, er := base64.URLEncoding.DecodeString(Lstr)

	if er != nil || len(enc_price) < 28 {
		return 0
	}

	iv := enc_price[0:16]
	price := enc_price[16:24]
	sig := enc_price[24:28]

	//fmt.Println(string(kla))
	e_key := []byte(E_key)
	e_key_d := hmac.New(sha1.New, e_key)
	e_key_d.Write(iv)
	e_key_end := e_key_d.Sum(nil)

	la := make([]byte, 8)
	for i := 0; i < 8; i++ {
		la[i] = price[i] ^ e_key_end[i]
	}

	//验证当前的签名
	laa := BytesCombine(la, iv)
	e_keya := []byte(I_key)
	e_keya_d := hmac.New(sha1.New, e_keya)
	e_keya_d.Write(laa)
	e_keya_end := e_keya_d.Sum(nil)

	sign_k := e_keya_end[0:4]
	istrue := bytes.Compare(sign_k, sig) //返回的值是零表示相等

	if istrue == 0 {
		v := BytesToInt64(la)
		return uint64(v)
	}
	return 0
}

func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

//生成签名
func CreatSign(key, url string) string {
	e_keya := []byte(key)
	urls := []byte(url)
	e_keya_d := hmac.New(sha1.New, e_keya)
	e_keya_d.Write(urls)
	e_keya_end := e_keya_d.Sum(nil)
	aa := base64.URLEncoding.EncodeToString(e_keya_end)

	return aa

	signature := base64.RawURLEncoding.EncodeToString(e_keya_end)

	for {
		if len(signature)%4 == 0 {
			break
		}
		signature += "="
	}
	return signature
}

//金额加密
func Encode_money(price int64) string {
	price_byte := Int64ToBytes(price)
	e_key_a := []byte(E_key)
	e_key_d := hmac.New(sha1.New, e_key_a)
	e_key_d.Write(Iv)
	e_key_end := e_key_d.Sum(nil)

	en_price := make([]byte, 8)
	for i := 0; i < 8; i++ {
		en_price[i] = e_key_end[i] ^ price_byte[i]
	}
	data := BytesCombine(price_byte, Iv)
	i_key := []byte(I_key)
	i_keya_d := hmac.New(sha1.New, i_key)
	i_keya_d.Write(data)
	i_keya_end := i_keya_d.Sum(nil)
	sign := i_keya_end[0:4]
	msg := BytesCombine(Iv, en_price, sign)
	enda := base64.RawURLEncoding.EncodeToString(msg)
	fmt.Println(enda)

	return enda //加密完成后的金额字符串
}

//func init() {

//	price_byte := Int64ToBytes(27348492)
//	fmt.Println("-----")
//	fmt.Println(price_byte)
//	fmt.Println("->>>>>>>")

//	e_key := []byte(E_key)
//	e_key_d := hmac.New(sha1.New, e_key)
//	e_key_d.Write(Iv)
//	e_key_end := e_key_d.Sum(nil)

//	en_price := make([]byte, 8)
//	for i := 0; i < 8; i++ {
//		en_price[i] = e_key_end[i] ^ price_byte[i]
//	}

//	fmt.Println(en_price)
//	//加密的数据
//	data := BytesCombine(price_byte, Iv)
//	i_key := []byte(I_key)
//	i_keya_d := hmac.New(sha1.New, i_key)
//	i_keya_d.Write(data)
//	i_keya_end := i_keya_d.Sum(nil)
//	fmt.Println("---------->")
//	sign := i_keya_end[0:4]

//	msg := BytesCombine(Iv, en_price, sign)
//	fmt.Println(msg)
//	enda := base64.RawURLEncoding.EncodeToString(msg)
//	fmt.Println(enda)

//	money := Decprice(enda)
//	fmt.Println(money)
//}
