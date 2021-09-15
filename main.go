package main

import (
	"customPro/autoServiceClasses/protoc-gen-speed/php"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	req, err := readRequest(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err = writeResponse(os.Stdout, php.Generate(req)); err != nil {
		panic(err)
	}
}

//读取请求
func readRequest(input io.Reader) (*plugin.CodeGeneratorRequest, error) {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		return nil, err
	}

	return req, nil
}

//将结果序列化
//最后输出
func writeResponse(out io.Writer, resp *plugin.CodeGeneratorResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = out.Write(data)
	return err
}

