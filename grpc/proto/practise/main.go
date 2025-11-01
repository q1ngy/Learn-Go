package main

import (
	"fmt"
	"practise/api"

	"github.com/iancoleman/strcase"
	"github.com/mennanov/fieldmask-utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func oneOfDemo() {
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

func optionalDemo() {
	// client
	book := api.Book{
		Title:  "golang",
		Price:  &wrapperspb.Int64Value{Value: 2000},
		Author: &wrapperspb.StringValue{Value: "slim"},
		Note:   proto.String("哈哈哈哈哈"),
	}

	// server
	if book.Price == nil {
	} else {
		fmt.Println("book price: ", book.GetPrice().GetValue())
	}

	if book.Note != nil {
		fmt.Println("book note: ", book.GetNote())
	}
}

func fieldMaskDemo() {
	// client
	paths := []string{"password", "friend.a"}
	req := api.UpdateUserReq{
		Op: "admin",
		User: &api.User{
			Username: "xxx",
			Password: "xxx",
			Friend:   &api.User_Friend{A: "aaa"},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: paths},
	}

	// server
	mask, _ := fieldmask_utils.MaskFromProtoFieldMask(req.UpdateMask, strcase.ToCamel)
	var userDst = make(map[string]interface{})
	fieldmask_utils.StructToMap(mask, req.User, userDst)
	fmt.Printf("user: %v\n", userDst)
}

func main() {
	oneOfDemo()
	optionalDemo()
	fieldMaskDemo()
}
