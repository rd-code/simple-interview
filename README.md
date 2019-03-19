# 说明
协议参考 https://redis.io/topics/protocol

用例参考 parse/parse_test.go文件

进行了初步测试，复杂测试还没有进行

所有的输入统一按照io.Reader进行输入

实现了 simpleString,Errors,Integers,Bulk Strings,Arrays