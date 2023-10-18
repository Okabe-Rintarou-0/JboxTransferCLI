# Jbox Transfer CLI

<img src="https://img.shields.io/badge/go-1.19-blue" alt="">

[JboxTransfer](https://github.com/1357310795/JboxTransfer) 配套 CLI 程序，基于 golang 开发。欢迎提出 issues 和 pr。

# 下载及安装

见 Release。

# 如何使用

+ 使用二维码登录：

    ```shell
    ./jtrans login --qrcode
    ```

+ 登出：

    ```shell
    ./jtrans logout
    ```

+ 同步所有文件：
  
    ```shell
    ./jtrans sync --all
    ```

+ 清除同步任务：

    ```shell
    ./jtrans sync clear
    ```

+ 查看新旧云盘差异：
    
    ```shell
    ./jtrans diff
    ```

+ 下载单个文件

    ```shell
    ./jtrans jbox download FILE_PATH
    ```

+ 上传单个文件

    ```shell
    ./jtrans tbox upload -f LOCAL_PATH -t REMOTE_PATH
    ```

+ 查看文件信息
  
    ```shell
    ./jtrans tbox info FILE_PATH
    ./jtrans jbox info FILE_PATH
    ```

+ 列出文件夹下的文件
  
    ```shell
    ./jtrans tbox ls REMOTE_DIR_PATH
    ./jtrans jbox ls REMOTE_DIR_PATH
    ```

+ 一些比较 geek 的操作

    + 下载文件并计算其 md5 校验和：
        ```shell
        ./jtrans jbox download FILE_PATH --bytes | ./jtrans md5
        ```
    + 下载文件并计算其 crc64 校验和：
        ```shell
        ./jtrans jbox download FILE_PATH --bytes | ./jtrans crc64
        ```
    + 下载文件并上传：
        ```shell
        ./jtrans jbox download FILE_PATH --bytes | ./jtrans tbox upload --bytes -t REMOTE_PATH
        ```

+ 更多操作请见：

    ```shell
    ./jtrans --help
    ```