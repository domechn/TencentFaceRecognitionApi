package TencentFaceRecognitionApi

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"time"	
	"crypto/sha1"
	"crypto/hmac"	
	"math/rand"
	"encoding/base64"
)

/**
* 通过图片url获取==调用api
* appid  
* mode
* imageUrl 图片路径
* sessionId
* sessionKey
* url 请求路径
*/
func PostByUrl(appid string,mode int,imageUrl ,sessionId , sessionKey , url string) {
	t := time.Now()
	currentUnix := t.Unix()
	resource := rand.NewSource(currentUnix)
	sourceRand  := rand.New(resource)

	//json序列化
	postData := fmt.Sprintf("{\"appid\":\"%s\",\"mode\":%d,\"url\":\"%s\"}",appid,mode,imageUrl)
	
	
	srcStr:=fmt.Sprintf("a=%s&k=%s&e=%d&t=%d&r=%d&u=0&f=",appid,sessionId,currentUnix+2400,currentUnix,sourceRand.Intn(999999999))
	
	
	key := []byte(sessionKey)
	hashHmac := hashHmac(srcStr,key)
	result := bytesCombine(hashHmac,[]byte(srcStr))
	
	encodeString := base64.StdEncoding.EncodeToString(result)
	
	fmt.Println(encodeString)
	
	
	request , _  := http.NewRequest("POST",url,bytes.NewBuffer([]byte(postData)))
	request.Header.Set("Content-Type","application/json;charset=UTF-8")
	request.Header.Set("host","recognition.image.myqcloud.com")
	request.Header.Set("authorization",	encodeString)
	client := &http.Client{}
	resp , _ := client.Do(request)
	respBytes , _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
}

func bytesCombine(pBytes ...[]byte) []byte{
	return bytes.Join(pBytes,[]byte(""))
}

func hashHmac(srcStr string,key []byte) []byte{
	mac := hmac.New(sha1.New,key)
	mac.Write([]byte(srcStr))
	return mac.Sum(nil)
}