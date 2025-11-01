package main

import (
	"fmt"
	"oneof/api"
)

func oneofDemo() {
	// client
	req1 := &api.NoticeReaderRequest{
		Msg: "...",
		NoticeWay: &api.NoticeReaderRequest_Email{
			Email: "123@123.com",
		},
	}

	//req2 := &api.NoticeReaderRequest{
	//	Msg: "...",
	//	NoticeWay: &api.NoticeReaderRequest_Phone{
	//		Phone: "123456789",
	//	},
	//}

	// server
	req := req1
	switch v := req.NoticeWay.(type) {
	case *api.NoticeReaderRequest_Email:
		noticeWithEmail(v)
	case *api.NoticeReaderRequest_Phone:
		noticeWithPhone(v)
	}

	//fmt.Print(req1, req2)
}

func noticeWithEmail(in *api.NoticeReaderRequest_Email) {
	fmt.Printf("notice reader by email:%v\n", in.Email)
}

func noticeWithPhone(in *api.NoticeReaderRequest_Phone) {
	fmt.Printf("notice reader by phone:%v\n", in.Phone)
}

func main() {
	oneofDemo()
}
