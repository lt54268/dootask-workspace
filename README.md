# dootask-workspace
## 接口使用说明
请求地址：ws://服务器IP:3333/ws
### 一、同步用户ID
 ```
{
  "action": "sync"
}
 ```
### 二、设置创建工作区授权
 ```
{
  "action": "set",
  "data": {
    "user_id": 1,       // 用户ID
    "is_create": true   // false
  }
}
 ```
### 三、创建工作区
 ```
{
  "action": "create",
  "data": {
    "user_id": 1       // 用户ID
  }
}
 ```
### 四、检查已创建的工作区数量
 ```
{
  "action": "check"
}
 ```
### 五、获取已创建工作区的用户ID
 ```
{
  "action": "get"
}
 ```
### 六、流式对话问答
 ```
{
  "action": "stream-chat",
  "data": {
    "message": "哈哈哈",
    "mode": "chat",
    "sessionId": "1",               // 对话ID，每个对话窗口唯一
    "slug": "workspace-for-user-1"
  }
}
 ```
### 七、常规对话问答
 ```
{
  "action": "chat",
  "data": {
    "message": "哈哈哈",
    "mode": "chat",
    "sessionId": "1", 
    "slug": "workspace-for-user-1"
  }
}
 ```
### 八、存储聊天对话
 ```
{
  "action": "back",
  "data": {
    "session_id": 1,
    "slug": "workspace-for-user-3",
    "last_message": "嘻嘻嘻"         // AI回答的最后一条消息
  }
}
```
### 九、获取最后一天聊天记录
```
{
  "action": "send",
  "data": 1         // 对话ID
}
```
