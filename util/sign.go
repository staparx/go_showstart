package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type GenerateSignReq struct {
	Path      string `json:"path"`
	Data      string `json:"data"`
	Cusat     string `json:"cusat"`
	Sign      string `json:"sign"`
	Cusit     string `json:"cusit"`
	Cusid     string `json:"cusid"`
	TraceId   string `json:"traceId"`
	Token     string `json:"token"`
	Cterminal string `json:"cterminal"`
}

func GenerateSign(req *GenerateSignReq) string {
	o := fmt.Sprintf(`%s%s%s%s%s%s%s%s%s%s%s`, req.Cusat, req.Sign, req.Cusit, req.Cusid, "wap", req.Token, req.Data, req.Path, "997", req.Cterminal, req.TraceId)
	return Md5Hex(o)
}

func Md5Hex(n string) string {
	// Compute MD5 hash
	hasher := md5.New()
	hasher.Write([]byte(n))
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
