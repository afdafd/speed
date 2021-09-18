### 快速生成gRpc服务端代里服务类代码

### 配置操作步骤：
- 1. 去github拉取代码：https://github.com/afdafd/speed
- 2. 已经生成好了对应的二进制文件:
     - 1. protoc-gen-speed (对应 MacOS或Linux系统)
     - 2. protoc-gen-speed.exe（对应windows系统）
- 3. 根据自己的系统选择对应的二进制文件；
- 4. 为了能全局使用，建议把文件放到对应的 */bin目录下；
     - 1. Linux或MacOs系统的话，可以放到/usr/local/bin目录下
     - 2. Windows系统的话，可以把文件设置环境变量里
     
### 基本使用：
- 说明：这个小工具是以 *.proto 文件来生成对应的代理服务类的。
- 1. 首先配置好我们自己的 *.proto 相关文件
- 2. 生成对应代理服务类：protoc --speed_out=. *.proto

#### 注：--speed_out可以单独使用，如上实例。也可以同时结合其他指令一并使用：
- 例：通过配置好的 *.proto 文件生成php的客户端、服务端代码，同时生成代理服务类代码：
- protoc --php_out=. --grpc_out=. --plugin=protoc-gen-grpc=/usr/local/bin/grpc_php_plugin --speed_out=. *.proto
> 生成好的代码如下：
```php
 /**
  * Auto Generated Service Classes
  * source: userCenter.proto
  *
  * Class Agent
  * @GrpcService(prefix="Agents/Agent")
  */
 class AgentService
 {
    /**
     * @param \Agents\BaseRequest $request
     * @return \Agents\ResultResponse
     */
     public function getJjAgentInfo(
         \Agents\BaseRequest $request
     ): \Agents\ResultResponse
     {
        try {
           $result = \bean(AgentsController::class)->getJjAgentInfo(\context()->getRequest());
           if ($result['code'] != 1) {
               throw new \Exception($result['message'] ?? '调用gRPC服务端[ getJjAgentInfo ]接口返回结果失败');
           }
 
           $response = new \Agents\ResultResponse;
 
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
               'callMethod' => 'Agents/Agent/getJjAgentInfo'
           ]);
 
           throw new \Exception($e->getMessage(), $e->getCode());
       }
     }
 }
```

####  --speed_out指令和其他指令一起使用时，没有先后顺序，可以随意组合。