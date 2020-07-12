# ArchLInux 安装笔记

## 验证启动方式

```
ls /sys/firmware/efi/efivars
```
目录存在以UEFI方式启动，目录不存在系统可能以 BIOS 或 CSM 模式启动。


## 连接到因特网

如果是有线一般自动连上了，尝试 ping baidu.com,如果是wifi可以用wifi-menu命令配置。下面是常用的网络命令指令。
```shell
ip link
ip link set echo0 up
wifi-menu
dhcpcd interface
```

## 更新时间
```
timedatectl set-ntp true
```
可以使用 timedatectl status 检查服务状态。 

##硬盘分区
磁盘若被系统识别到，就会被分配为一个块设备，如 /dev/sda 或者 /dev/nvme0n1。可以使用 lsblk 或者 fdisk 查看。

### 分区示例：
#### BIOS 和 MBR分区
| 挂载点 | 分区      | 分区类型              | 建议大小     |
| ------ | ------    | -----                 | -----        |
| /mnt   | /dev/sdX1 | Linux                 | 剩余空间     |
| [SWAP] | /dev/sdX2 | Linux swap (交换空间) | 大于 512 MiB |
    
#### UEFI with GPT
    
| 挂载点                | 分区                    | 分区类型                | 建议大小     |
|-----------------------|-------------------------|-------------------------|--------------|
| /mnt/boot or /mnt/efi | /dev/sdX1  EFI 系统分区 | 260–512 MiB             |              |
| /mnt                  | /dev/sdX2               | Linux x86-64 根目录 (/) | 剩余空间     |
| [SWAP]                | /dev/sdX3               | Linux swap (交换空间)   | 大于 512 MiB |

### 格式化分区
当分区建立好了，这些分区都需要使用适当的文件系统进行格式化。举个例子，如果根分区在 /dev/sdX1 上并且会使用 ext4 文件系统，运行：
 ```
 mkfs.vfat -F32 /dev/sdx1
 mkfs.ext4 /dev/sdx2
 mkswap /dev/sdx3 #格式化交换分区
 swapon /dev/sdX3 #启用交换分区
 ```
 
 ### 挂载分区
 ```
 mount /dev/sdX2 /mnt  # 根分区
 mkdir -p /mnt/boot/efi #efi分区
 mount /dev/sdx1 /mnt/boot/efi
 ```
 
 ## 安装
 ### 选择镜像
 文件 /etc/pacman.d/mirrorlist 定义了软件包会从哪个镜像源下载，修改文件 /etc/pacman.d/mirrorlist，越前的镜像在下载软件包时有越高的优先权。这个文件接下来还会被 pacstrap 拷贝到新系统里，所以请确保设置正确。
 
 ```
 Server = http://mirrors.aliyun.com/archlinux/$repo/os/$arch
 Server = https://mirrors.cloud.tencent.com/archlinux/$repo/os/$arch
 Server = http://mirrors.163.com/archlinux/$repo/os/$arch
 Server = http://mirrors.tuna.tsinghua.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirrors.bfsu.edu.cn/archlinux/$repo/os/$arch
 Server = https://mirrors.bfsu.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirrors.cqu.edu.cn/archlinux/$repo/os/$arch
 Server = https://mirrors.cqu.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirrors.dgut.edu.cn/archlinux/$repo/os/$arch
 Server = https://mirrors.dgut.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirror.lzu.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirrors.neusoft.edu.cn/archlinux/$repo/os/$arch
 Server = https://mirrors.neusoft.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirror.redrock.team/archlinux/$repo/os/$arch
 Server = https://mirror.redrock.team/archlinux/$repo/os/$arch
 Server = https://mirrors.sjtug.sjtu.edu.cn/archlinux/$repo/os/$arch
 Server = https://mirrors.tuna.tsinghua.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirrors.ustc.edu.cn/archlinux/$repo/os/$arch
 Server = https://mirrors.ustc.edu.cn/archlinux/$repo/os/$arch
 Server = https://mirrors.xjtu.edu.cn/archlinux/$repo/os/$arch
 Server = http://mirrors.zju.edu.cn/archlinux/$repo/os/$arch
```

### 安装系统
```shell
 pacstrap /mnt base linux linux-firmware
```

### 配置系统
```
genfstab -U /mnt >> /mnt/etc/fstab #设置系统启动挂载分区

arch-chroot /mnt   # 切换到新安装的系统

ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime #设置时区
hwclock --systohc #更新时间

vim /etc/locale.gen #修改系统支持的语言，把 en_US.UTF-8 zh_CN.UTF-8前面的"#"去掉
en_US.UTF-8 
zh_CN.UTF-8

echo LANG=en_US.UTF-8 > /etc/locale.conf #将系统 locale 设置为 en_US.UTF-8，系统的 Log 就会用英文显示，
#这样更容易问题的判断和处理。

passwd # 设置root密码
```

### 安装引导程序 GRUB

```
pacman -S dosfstools grub efibootmgr #安装GRUP相关软件包

```

#### BIOS MBR
```
grub-install --target=i386-pc /dev/sdX
```
#### UEFI
```
grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=GRUB
```

### 生成grub.cfg
```
grub-mkconfig -o /boot/grub/grub.cfg
```

### 重启
```
exit
reboot
```


#### 安装 ucode 
```
pacman -S intel-ucode  #intel cpu
pacman -S amd-ucode    #amd cpu
```


####新建用户
```
useradd -m -g users -G wheel -s /bin/bash usernaem
passwd username

visudo # 在 root ALL=(ALL) ALL 下面添加 用户名 ALL=(ALL) ALL
```

#### 网络配置
##### 有线连接
```
systemctl start dhcpcd
systemctl enable dhcpcd
```
##### 无线连接
```
pacman -S netctl iw wpa_supplicant dialog
wifi-menu
```

## 安装桌面环境
####安装显卡驱动
```
pacman -S 驱动包
```
| 显卡       | 驱动包             |
| -          | -                  |
| 通用       | xf86-video-vesa    |
| intel-     | xf86-video-intel   |
| amdgpu     | xf86-video-amdgpu  |
| Geforce7±  | xf86-video-nouveau |
| Geforce6/7 | xf86-video-304xx   |
| ati        | xf86-video-ati     |

#### 安装Xorg
```
pacman -S xorg
```

####安装字体
```
pacman -S ttf-dejavu wqy-microhei
```
####安装桌面环境DWM
```
pacman -S dwm
cp /etc/X11/xinit/xinitrc ~/.xinitrc  # 拷贝Xorg 启动初始化文件
vim  ~/.xinitc #添加exec dwm
startx ### 启动到桌面
```


##常用的软件

acpi
adobe-source-code-pro-fonts
alacritty
alsa-utils
archlinuxcn-keyring
bat
calcurse
calibre
chromium
code
cronie
ctags
docker
efibootmgr
electron-netease-cloud-music
fcitx-configtool
fcitx-googlepinyin
fcitx-lilydjwg-git
fcitx-qt5
fcitx-sogoupinyin
feh
firefox
firefox-i18n-zh-cn
fish
flameshot
foxitreader
fzf
gimp
git
google-chrome
grub
hsetroot
htop
ifuse
intel-ucode
intellij-idea-ultimate-edition
lazygit
lf-bin
libreoffice-still
libreoffice-still-zh-cn
libxft-bgra
linux
linux-firmware
llpp
mpv
neofetch
neovim
nerd-fonts-source-code-pro
nerd-fonts-ubuntu
net-tools
netease-cloud-music
networkmanager
noto-fonts-cjk
noto-fonts-emoji
npm
ntfs-3g
ntp
obs-studio
openssh
os-prober
pamixer
parted
patch
pavucontrol
picom
pkgfile
proxychains-ng
python-pynvim
python2
qq-linux
ranger
redshift
rofi
sudo
sxiv
the_silver_searcher
tmux
trojan
ttf-dejavu
ttf-joypixels
ttf-meslo
typora
ueberzug
unrar
unzip-iconv
vim
virtualbox
vlc
w3m
wget
which
wiznote
wmname
wqy-microhei
xautolock
xclip
xf86-video-intel
xorg-bdftopcf
xorg-docs
xorg-font-util
xorg-fonts-100dpi
xorg-fonts-75dpi
xorg-iceauth
xorg-luit
xorg-mkfontscale
xorg-server
xorg-server-devel
xorg-server-xephyr
xorg-server-xnest
xorg-server-xvfb
xorg-server-xwayland
xorg-sessreg
xorg-smproxy
xorg-x11perf
xorg-xauth
xorg-xbacklight
xorg-xcmsdb
xorg-xcursorgen
xorg-xdpyinfo
xorg-xdriinfo
xorg-xev
xorg-xgamma
xorg-xhost
xorg-xinit
xorg-xinput
xorg-xkbevd
xorg-xkbutils
xorg-xkill
xorg-xlsatoms
xorg-xlsclients
xorg-xmodmap
xorg-xpr
xorg-xprop
xorg-xrandr
xorg-xrdb
xorg-xrefresh
xorg-xset
xorg-xsetroot
xorg-xvinfo
xorg-xwd
xorg-xwininfo
xorg-xwud
xsel
yay
ydcv-rs-git
zathura
zathura-cb
zathura-djvu
zathura-pdf-mupdf
zathura-ps
zsh
