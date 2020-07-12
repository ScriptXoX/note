Linux中进行挂起(待机)的命令说明

在Linux中有挂起（也即是待机）的命令为：
rtcwake

关于此命令的很多其它帮助，能够在命令终端进行man:

man rtcwake

关于此命令的简单使用的样例如：

rtcwake -m mem -s 60

表示系统挂起的时候是把当前系统的状态信息等保存到内存中，挂起时间为60秒，即在60秒后会自己主动唤醒；

rtcwake

Linux


