package trade

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	alCreateMethod = "alipay.trade.precreate"
	alQueryMethod  = "alipay.trade.query"
	// alCharset      = "utf-8"
	// alSignType     = "RSA"
	nameCharset   = "charset"
	valueCharset  = "utf-8"
	nameSignType  = "sign_type"
	nameAppID     = "app_id"
	nameSign      = "sign"
	nameTimestamp = "timestamp"
)

type signer interface {
	sign(data string) ([]byte, error)
	verify(data []byte, sign string) error
}

type rsaSigner struct {
	key    *rsa.PrivateKey //用于发送数据的私钥
	pubKey *rsa.PublicKey  //阿里账号的公钥
}

type rsa2Signer struct {
	key    *rsa.PrivateKey //用于发送数据的私钥
	pubKey *rsa.PublicKey  //阿里账号的公钥
}

func (signer rsaSigner) sign(data string) ([]byte, error) {
	hashed := sha1.Sum([]byte(data))
	return rsa.SignPKCS1v15(rand.Reader, signer.key, crypto.SHA1, hashed[:])
}

func (signer rsaSigner) verify(data []byte, sign string) error {
	if sign == "" {
		return errors.New(string(data))
	}
	hashed := sha1.Sum(data)
	s, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(signer.pubKey, crypto.SHA1, hashed[:], s)
}

func (signer rsa2Signer) sign(data string) ([]byte, error) {
	// h := sha256.New()
	// h.Write([]byte(data))
	// fmt.Printf("%x", h.Sum(nil))
	hashed := sha256.Sum256([]byte(data))
	return rsa.SignPKCS1v15(rand.Reader, signer.key, crypto.SHA256, hashed[:])
}

func (signer rsa2Signer) verify(data []byte, sign string) error {
	if sign == "" {
		return errors.New(string(data))
	}
	hashed := sha256.Sum256(data)
	s, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(signer.pubKey, crypto.SHA256, hashed[:], s)
}

type (
	//preCreateContent 阿里创建预付的内容
	preCreateContent struct {
		OutTradeNo  string `json:"out_trade_no"`
		TotalAmount string `json:"total_amount"`
		Subject     string `json:"subject"`
		// TimeoutExpress string `json:"timeout_express"`
	}
	//preCreateRequest 阿里创建预付的包装数据
	preCreateRequest struct {
		// AppID      string `json:"app_id"`
		Method string `json:"method"`
		// Charset    string `json:"charset"`
		// SignType   string `json:"sign_type"`
		// Sign       string `json:"sign"`
		// Timestamp  string `json:"timestamp"`
		NotifyURL  string `json:"notify_url"`
		BizContent string `json:"biz_content"`
	}
	//PreCreateResult 阿里创建预付的结果
	preCreateResult struct {
		AlipayTradePrecreateResponse json.RawMessage `json:"alipay_trade_precreate_response"`
		Sign                         string          `json:"sign"`
	}
)

type (
	//queryContent 查询交易的内容
	queryContent struct {
		OutTradeNo string `json:"out_trade_no"`
	}
	//queryRequest 阿里查询预付的包装数据
	queryRequest struct {
		// AppID      string `json:"app_id"`
		Method string `json:"method"`
		// Charset    string `json:"charset"`
		// SignType   string `json:"sign_type"`
		// Sign       string `json:"sign"`
		// Timestamp  string `json:"timestamp"`
		BizContent string `json:"biz_content"`
	}
	//queryResult 查询交易的结果
	queryResult struct {
		AlipayTradeQueryResponse json.RawMessage `json:"alipay_trade_query_response"`
		Sign                     string          `json:"sign"`
	}
)

//request 向阿里服务器发起请求
func (client *AlipayClient) request(data interface{}, result interface{}) error {
	// client.sign(data)
	// query := alMarshal(data)
	query := client.getSignedData(data)
	url := fmt.Sprintf("%s?%s", client.gatewayURL, query)
	// fmt.Printf("request: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("response: %s\n", body)
	return json.Unmarshal(body, result)
}

//sign 对要发送的数据按阿里规则签名
func (client *AlipayClient) getSignedData(data interface{}) string {
	t := reflect.TypeOf(data).Elem()
	v := reflect.ValueOf(data).Elem()
	pairs := make([]string, 0, 64)
	values := url.Values{}
	for i := 0; i < t.NumField(); i++ {
		var name = t.Field(i).Tag.Get("json")
		var value = v.Field(i).Interface()
		if v, ok := value.(string); ok && v != "" {
			pairs = append(pairs, name+"="+v)
			values.Set(name, v)
		}
	}
	pairs = append(pairs, nameCharset+"="+valueCharset)
	values.Set(nameCharset, valueCharset)
	signType := string(client.signType)
	pairs = append(pairs, nameSignType+"="+signType)
	values.Set(nameSignType, signType)
	pairs = append(pairs, nameAppID+"="+client.appID)
	values.Set(nameAppID, client.appID)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	pairs = append(pairs, nameTimestamp+"="+timestamp)
	values.Set(nameTimestamp, timestamp)
	sort.Strings(pairs)
	var str = strings.Join(pairs, "&")
	// fmt.Println(len(pairs), cap(pairs), str)
	buf, err := client.signer.sign(str)
	if err != nil {
		panic(err)
	}
	sign := base64.StdEncoding.EncodeToString(buf)
	values.Set(nameSign, sign)
	return values.Encode()
}

// //alMarshal 把数据序列化成url格式
// func alMarshal(data interface{}) string {
// 	t := reflect.TypeOf(data)
// 	v := reflect.ValueOf(data)
// 	if v.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 		v = v.Elem()
// 	}
// 	d := url.Values{}
// 	for i := 0; i < t.NumField(); i++ {
// 		var name = t.Field(i).Tag.Get("json")
// 		var value = v.Field(i).Interface()
// 		if v, ok := value.(string); ok && v != "" {
// 			d.Set(name, v)
// 		}
// 	}
// 	return d.Encode()
// }

// //sign 对要发送的数据按阿里规则签名
// func (client *AlipayClient) sign(data interface{}) error {
// 	t := reflect.TypeOf(data).Elem()
// 	v := reflect.ValueOf(data).Elem()
// 	pairs := make([]string, 0, 64)
// 	for i := 0; i < t.NumField(); i++ {
// 		var name = t.Field(i).Tag.Get("json")
// 		var value = v.Field(i).Interface()
// 		if v, ok := value.(string); ok && v != "" && name != "sign" {
// 			pairs = append(pairs, name+"="+v)
// 		}
// 	}
// 	sort.Strings(pairs)
// 	var str = strings.Join(pairs, "&")
// 	//	fmt.Println(len(pairs), cap(pairs), str)
// 	hashed := sha1.Sum([]byte(str))
// 	s, err := rsa.SignPKCS1v15(rand.Reader, client.key, crypto.SHA1, hashed[:])
// 	if err != nil {
// 		return err
// 	}
// 	v.FieldByName("Sign").SetString(base64.StdEncoding.EncodeToString(s))
// 	return nil
// }

// //verify 验证应答的签名
// func (client *AlipayClient) verify(data []byte, sign string) error {
// 	if sign == "" {
// 		return errors.New(string(data))
// 	}
// 	hashed := sha1.Sum(data)
// 	s, err := base64.StdEncoding.DecodeString(sign)
// 	if err != nil {
// 		return err
// 	}
// 	return rsa.VerifyPKCS1v15(client.pubKey, crypto.SHA1, hashed[:], s)
// }
