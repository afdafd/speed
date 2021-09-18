package php

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const phpBody = `<?php declare(strict_types=1);

{{ $ns := .Namespace -}}
namespace App\Grpc\Services\{{ $ns.Namespace }};

use Swoft\Log\Helper\Log;
use Hzwz\Grpc\Server\Annotation\Mapping\GrpcService;
use App\Http\Controller\{{ $.Service.Name }}Controller;


/**
 * Auto Generated Service Classes
 * source: {{ .File.Name }}
 *
 * Class {{ .Service.Name }}
 * @GrpcService(prefix="{{ $ns.Special }}/{{ $.Service.Name }}") 
 */
class {{ .Service.Name | Service }}
{
{{- range $m := .Service.Method}}
   /**
    * @param \{{ $ns.Namespace }}\{{ name $ns $m.InputType }} $request
    * @return \{{ $ns.Namespace }}\{{ name $ns $m.OutputType }}
    */
    public function {{ $m.Name }}(
        \{{ $ns.Namespace }}\{{ name $ns $m.InputType }} $request
    ): \{{ $ns.Namespace }}\{{ name $ns $m.OutputType }}
    {
       try {
          $result = \bean({{ $.Service.Name }}Controller::class)->{{ $m.Name }}(\context()->getRequest());
          if ($result['code'] != 1) {
              throw new \Exception($result['message'] ?? '调用gRPC服务端[ {{ $m.Name }} ]接口返回结果失败');
          }

          $response = new \{{ $ns.Namespace }}\{{ name $ns $m.OutputType }};

          if (empty($result['data'])) {
              return $response;
          }
       
          foreach ($result['data'] as $filed => $value) {
              $response->{'set' . $filed}($value);
          }
          
          return $response;
      } catch (\Throwable $e) {
          Log::error("gRpcServiceCallError", [
              'errorMsg'   => $e->getMessage(),
              'errorSite'  => $e->getFile() .'|'. $e->getLine(),
              'callMethod' => '{{ $ns.Namespace }}/{{ $.Service.Name }}/{{ $m.Name }}'
          ]);

          throw new \Exception($e->getMessage(), $e->getCode());
      }
    }
{{end -}}
}
`

var tpl1 *template.Template

func init() {
	tpl1 = template.Must(template.New("phpBody").Funcs(template.FuncMap{
		"Service": func(name *string) string {
			return identifier(*name, "Service")
		},
		"name": func(ns *ns, name *string) string {
			return ns.resolve(name)
		},
	}).Parse(phpBody))
}

//生成php文件名
func filename(file *descriptor.FileDescriptorProto, name *string) string {
	ns := namespace(file.Package, "/")
	if file.Options != nil && file.Options.PhpNamespace != nil {
		ns = strings.Replace(*file.Options.PhpNamespace, `\`, `/`, -1)
	}

	return fmt.Sprintf("%s/%s.php", ns, identifier(*name, "Service"))
}

// 生成php的主体
func body(
	req *plugin.CodeGeneratorRequest,
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
) string {
	out := bytes.NewBuffer(nil)

	data := struct {
		Namespace *ns
		File      *descriptor.FileDescriptorProto
		Service   *descriptor.ServiceDescriptorProto
	}{
		Namespace: newNamespace(req, file, service),
		File:      file,
		Service:   service,
	}

	err := tpl1.Execute(out, data)
	if err != nil {
		panic(err)
	}

	return out.String()
}
