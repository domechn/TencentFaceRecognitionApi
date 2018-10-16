package TencentFaceRecognitionApi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type TencentAPI struct {
	appid     string
	mode      int
	secretID  string
	secretKey string
	url       string
}

func NewTencentAPI(appid string, mode int, secretID, secretKey, url string) *TencentAPI {
	return &TencentAPI{
		appid:     appid,
		mode:      mode,
		secretID:  secretID,
		secretKey: secretKey,
		url:       url,
	}
}

/**
* 通过图片url获取==调用api
* appid
* mode
* imageUrl 图片路径
* sessionId
* sessionKey
* url 请求路径
 */
func (tAPI *TencentAPI) PostByUrl(imageUrl string) string {
	//json序列化
	postData := fmt.Sprintf("{\"appid\":\"%s\",\"mode\":%d,\"url\":\"%s\"}", tAPI.appid, tAPI.mode, imageUrl)

	encodeString := tAPI.sign()
	request, err := http.NewRequest("POST", tAPI.url, bytes.NewBuffer([]byte(postData)))
	if err != nil {
		log.Panic(err)
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("host", "recognition.image.myqcloud.com")
	request.Header.Set("authorization", encodeString)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return string(respBytes)
}

func (tAPI *TencentAPI) postByFile(filePath string) string {
	fileByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}
	//json序列化
	fmt.Println(fileByte)
	postData := fmt.Sprintf("{\"appid\":\"%s\",\"mode\":%d,\"image\":\"%b\"}", tAPI.appid, tAPI.mode, fileByte)
	fmt.Println(postData)

	encodeString := tAPI.sign()

	request, err := http.NewRequest("POST", tAPI.url, bytes.NewBuffer([]byte(postData)))
	if err != nil {
		log.Panic(err)
	}
	request.Header.Set("Content-Type", "multipart/form-data")
	request.Header.Set("host", "recognition.image.myqcloud.com")
	request.Header.Set("authorization", encodeString)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(encodeString)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return string(respBytes)
}

func bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func hashHmac(srcStr string, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(srcStr))
	return mac.Sum(nil)
}

func (tAPI *TencentAPI) sign() string {
	t := time.Now()
	currentUnix := t.Unix()
	resource := rand.NewSource(currentUnix)
	sourceRand := rand.New(resource)

	srcStr := fmt.Sprintf("a=%s&k=%s&e=%d&t=%d&r=%d&u=0&f=", tAPI.appid, tAPI.secretID, currentUnix+10, currentUnix, sourceRand.Intn(999999999))

	key := []byte(tAPI.secretKey)
	hashHmac := hashHmac(srcStr, key)
	result := bytesCombine(hashHmac, []byte(srcStr))

	encodeString := base64.StdEncoding.EncodeToString(result)
	return encodeString
}
