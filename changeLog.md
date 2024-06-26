### v2022_09_20

1.重写即时通讯架构，将map全面加锁

### v2022_09_21

1. 获取详细的招聘信息(GetRecruitInfo)增加redis命中模式
2. 简历投递新增是否存在的检测

### v2022_09_22

1. 优化redis命中的检测机制

### v2022_09_30

1. 修改其他需求里的隐藏逻辑
2. 优化boss个人发布界面的隐藏访问逻辑

### v2022_10_02

1. 关闭了控制台的debug输出
2. 修改了收藏小标签不显示的bug
3. 修改了查看公司发布全部不准确的bug

### v2022_10_08

1. 修复了管理员删公司内个人身份时的事务bug

### v2022_10_09

1. 优化了标签的修改的逻辑bug，用事务囊括了修改标签下的发布内容的标签（job_label）

### v2022_10_12

1. 全局修改图片在数据库的保存格式，摆脱域名限制（调用format里的函数实现图片格式化输出）
2. 优化诸多关于labels表格的数据库查询（向articles冗余字段，减少inner join）

### v2022_10_13

1. 添加置顶排序逻辑
2. 优化需求页面查找数据的逻辑，添加索引
3. 更新creatorSql文件夹
4. 优化置顶排序逻辑，自动累加

### v2022_10_15

1. 添加数据库及图片数据自动备份，手动备份
2. 优化照片存储路径逻辑
3. 添加各类数量统计

### v2022_10_15

1.优化收藏的排序逻辑
2.优化用户查看收藏的数据库执行逻辑，全部命中索引
3.增加静态文件夹的中间件检查

### v2022_10_18

1. 增加带审核公司数据导出xls的功能

### v2022_10_19

1.增加记录用户创建时间
2.优化公司查询及显示逻辑
3.优化模糊查询逻辑
4.增加导出公司全部发布信息xls的功能

### v2022_10_22

1.改变websocket链接策略

### v2022_10_23

1.优化管理员端媒体存储链接策略

### v2022_10_30

1.修改会员风采bug及添加字段

### v2022_11_01

1.优化照片删除逻辑
2.优化部分修改逻辑
3.减少无用的日志输出

### v2022_11_10

1.增加首页全部信息置顶的逻辑

### v2022_11_11

1.增加首页宣传栏（焦点）的redis永久缓存方式
2.优化删除文件的返回逻辑

### v2022_11_12

1.优化msgobj的字段更新逻辑
2.放弃trainers表的批量已读更新，仅设msgobjs的已读更新
3.降低删除照片接口的权级

### v2022_11_16

1.封装图片保存的函数
2.新增图片压缩机制
3.新增图片裁剪机制
4.部分接口改为事务

### v2022_1224

1.修改备份删除失败bug
2.完善备份过期检测

### v2022_0119

1.修改备份逻辑