# 数据库准备工作
在正式进入演示工程前，为了保证大家在后面看到的配置均和我的一致，我们利用在 mysql
中添加一个用户`ugozero`，然后创建一个`heygozero`库。

> mysql安装教程请参考[mysql使用说明文档](https://dev.mysql.com/doc/)

# 添加mysql用户
* 登录mysql

    ``` shell script
    $ mysql -h 127.0.0.1 -uroot -p
    Enter password:
    ```
  
* 添加用户`ugozero`,这里就不设置密码了。

    ``` mysql
    mysql> create user 'ugozero'@'127.0.0.1' identified by '';
    Query OK, 0 rows affected (0.02 sec)
    ```
  
* 给用户`ugozero`授权数据库和表权限

    ``` mysql
    mysql> grant all privileges on *.* to 'ugozero'@'127.0.0.1';
    Query OK, 0 rows affected (0.02 sec)
    ```
    
    > 注意：这里为了后续能够利用goctl根据datasource生成model，因此将所有数据库和表的权限都授予了`ugozero`,从安全角度考虑，线上环境请结合实际业务操作，这里仅演示。

* 给用户`ugozero`授权读写权限

    ``` mysql
    mysql> grant all privileges on *.* to 'ugozero'@'127.0.0.1'  WITH GRANT OPTION;
   Query OK, 0 rows affected (0.00 sec)
    ```
    
    > 注意：这里为了后续能够利用goctl根据datasource生成model，因此都授予了`ugozero`读写权限,从安全角度考虑，线上环境请结合实际业务操作，这里仅演示。

* 刷新新用户
    
    ``` mysql
    mysql> FLUSH PRIVILEGES;
    Query OK, 0 rows affected (0.01 sec) 
    ```

# 连接测试

```shell script
$ mysql -h 127.0.0.1 -uugozero -p
  Enter password:
  Welcome to the MySQL monitor.  Commands end with ; or \g.
  Your MySQL connection id is 191
  Server version: 8.0.21 MySQL Community Server - GPL
  
  Copyright (c) 2000, 2020, Oracle and/or its affiliates. All rights reserved.
  
  Oracle is a registered trademark of Oracle Corporation and/or its
  affiliates. Other names may be trademarks of their respective
  owners.
  
  Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
  
  mysql>
```

# 创建DB

``` mysql
mysql> create database if not exists heygozero character set utf8mb4 collate utf8mb4_general_ci;
Query OK, 1 row affected (0.01 sec)
```

# End

上一篇 [《创建工程》](./project-create.md)

下一篇 [《常见FAQ集合》](./faq.md)

# 猜你想

* [《目录说明》](../index.md)