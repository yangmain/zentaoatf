<?php
<<<TC
CaseId:   200
Title:    登录失败账号锁定策略
Steps:    @开头的为含验证点的步骤
   step2000         连续输入3次错误的密码
   @step2010        第4次尝试登录
   group2100        不连续输入3次错误的密码
      step2101      输入2次错误的密码
      step2102      输入1次正确的密码
      step2103      再输入1次错误的密码
      @step2104     再输入1次正确的密码

expects:
# 
/* @step2010期望结果, 可以有多行 */

# 
/* @step2104期望结果, 可以有多行 */

readme:
- 脚本输出日志和expects章节中，#号标注的验证点需保持一致对应
- 脚本中/* */注释的为编写代码的位置，//符号注释的为说明文字
- 参考样例https://github.com/easysoft/zentaoatf/tree/master/xdoc/sample

TC;

/* 此处编写操作步骤代码 */

echo "#\n";  // @step2010: 系统提示账号被锁定
/* 输出验证点实际结果 */

echo "#\n";  // @step2104: 登录成功，账号未被锁定
/* 输出验证点实际结果 */

?>
