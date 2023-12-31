# temp-files

Save temporary files in the server. (保存臨時檔案到服務器)

## 定位

- 保存临时文件
- 文件体积不宜太大
- 保存简短的纯文字笔记

## 输入

- 通过网页上传文件
- 通过网页新建或编辑笔记
- 通过网页删除文件或笔记

## 输出

- 一个简单列表 (包括文件与笔记, 按生成时间排序)
- 每个文件有按钮: 删除、下载
- 每个笔记有按钮: 删除、下载、编辑
- 另有一个上传文件的控件
- 另有一个新建笔记的按钮

## 特点

- 尽可能简单 (功能尽量少, 代码尽量少)
- 不使用数据库

### 不加密

- 为了尽量减少代码，不加密
- 并且，管理员密码保存在浏览器（这样可以减少处理 cookies 的代码）

### 不使用数据库

- 上传后，自动在文件名前附加时间戳
- 下载时，自动删除文件名前的时间戳
- 凡是后缀名为 ".txt" 或 ".md" 的文件都视为笔记

### 新建笔记

- 通过网页新建的笔记一律以 ".txt" 作为扩展名
- 新建笔记时，用户需要填写文件名
  - *待定* 文件名只能使用 0-9, a-z, A-Z, _(下劃線), -(連字號), .(點)

## 目标用户

- 我自己
- 懂编程基础且拥有 VPS 的人

