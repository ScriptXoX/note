package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

// 目标坐标 目标颜色
var targetX, targetY, targetColor = 605, 135, []string{"a", "9"}

// 战斗状态坐标 战斗状态颜色
var combatStatusX, combatStatusY, combatColor = 418, 923, []string{"f"}

// 血量坐标  血量颜色
var hpX, hpY, hpColor = 796, 906, []string{"9", "a", "b", "c", "d", "e", "f"}

// 是否在水中颜色坐标 和 颜色
var waterX, waterY, waterColor = 1184, 929, []string{"a5"}

// 定义E R F 三个武器技能释放时间
var indexE, indexR, indexF int64 = time.Now().Unix() - 6, time.Now().Unix() - 9, time.Now().Unix() - 14

func main() {

	// for i := 0; i < len(targetColor); i++ {
	// 	fmt.Println(targetColor[i])
	// }

}

func gw2bot() {

	//等待2秒 进入激活窗口
	robotgo.Sleep(2)

	for i := 0; ; i++ {

		//1.自动锁怪（利用游戏机制 选择最近的怪,如果没找到怪，转向60度，继续向前走）
		fmt.Print("开始寻怪")
		getTarget()

		//2.尝试攻击 （按1技能和机甲6技能 如果进入战斗中开始打怪，如果没有进入战斗 调正方向到目标先前走一段距离 再次尝试超过3次 如果还是没有进入战斗状态则放弃向后走）
		if !tryAttack() {
			fmt.Println("没有打到怪，继续寻怪")
			continue
		}
		//3.攻击 每次释放技能检查目标如果血条消失停止攻击继续所怪  如果攻击超过50次 则放弃向后走
		fmt.Println("开始打怪")
		doAttack()
		//4..在尝试攻击和自动寻怪的时候 检测是否走到水里去了如果是，向后走20s

	}

}

// 自动寻怪
func getTarget() {
	//处理在水中
	processInWater()

	fmt.Println("自动寻怪")
	robotgo.KeyPress("m")
	robotgo.MilliSleep(300)

	if isGetTarget() {
		fmt.Println("找到目标怪")
	} else {
		fmt.Println("没有找到目标怪 换方向 向前走")

		//换方向
		robotgo.KeyDown("alt")
		robotgo.KeyDown("a")
		robotgo.MilliSleep(800)
		robotgo.KeyUp("a")
		robotgo.KeyUp("alt")

		robotgo.MilliSleep(100)

		//跳着向前走
		runWithJump(6)

	}

}

// 尝试攻击
func tryAttack() bool {

	processInWater()

	robotgo.KeyPress("6")
	robotgo.MilliSleep(200)
	robotgo.KeyPress("q")
	robotgo.MilliSleep(200)

	fmt.Println("调正视角")
	robotgo.MilliSleep(600)
	robotgo.KeyDown("alt")
	robotgo.KeyPress("a")
	robotgo.KeyUp("alt")
	robotgo.MilliSleep(200)
	robotgo.KeyDown("alt")
	robotgo.KeyPress("d")
	robotgo.KeyUp("alt")
	robotgo.MilliSleep(200)

	fmt.Println("检测是否在战斗状态")

	return isCombat()

}

// 进行攻击
func doAttack() {
	for i := 0; i < 50; i++ {
		if !isGetTarget() {
			fmt.Println("目标消失,停止打怪")
			return
		}

		processHP()

		robotgo.Sleep(1)

		nowSec := time.Now().Unix()
		fmt.Println("E", nowSec, indexE)
		if nowSec-indexE >= 6 {
			fmt.Println("放E 技能")
			robotgo.KeyPress("e")
			robotgo.MilliSleep(500)

			indexE = time.Now().Unix()
		}

		if nowSec-indexR >= 9 {
			fmt.Println("放R技能")
			robotgo.KeyPress("r")
			robotgo.MilliSleep(500)
			indexR = time.Now().Unix()

		}

		if nowSec-indexF >= 14 {
			fmt.Println("放F 技能")
			robotgo.KeyPress("f")
			robotgo.MilliSleep(500)
			indexF = time.Now().Unix()

		}

	}

}

// 处理在水中
func processInWater() {
	color := robotgo.GetPixelColor(waterX, waterY, 0)
	color = color[0:2]
	fmt.Println(color)
	if colorMatch(color, waterColor) {
		fmt.Println("走到水中去了，掉头")

		robotgo.KeyDown("alt")
		robotgo.MilliSleep(200)
		robotgo.KeyPress("w")
		robotgo.KeyUp("alt")

		runWithJump(10)

	}
}

// 血量处理
func processHP() {
	heathcolor := robotgo.GetPixelColor(hpX, hpY, 0)
	heathcolor = heathcolor[0:1]
	if !colorMatch(heathcolor, hpColor) {

		fmt.Println("翻滚 加血 转身 往回走")

		//翻滚
		robotgo.KeyPress("v")
		robotgo.MilliSleep(500)

		//加血

		robotgo.KeyPress("t")
		robotgo.MilliSleep(800)

		robotgo.KeyDown("alt")
		robotgo.KeyPress("w")
		robotgo.KeyUp("alt")

		runWithJump(10)
	}
}

func isCombat() bool {

	combatc := ""

	combatcolor := robotgo.GetPixelColor(combatStatusX, combatStatusY, 0)
	combatc = combatcolor[0:1]
	fmt.Println("战斗状态颜色", combatcolor, "-----", combatc)
	if combatc != "f" {
		fmt.Println("战斗中")
		return true
	} else {
		fmt.Println("脱战")
		return false
	}

}

func colorMatch(color string, srcColors []string) bool {
	for i := 0; i < len(srcColors); i++ {
		if color == srcColors[i] {
			return true
		}
	}
	return false
}

func isGetTarget() bool {
	color := robotgo.GetPixelColor(targetX, targetY, 0)
	fmt.Print("目标颜色:", color)
	color = color[0:1]
	fmt.Println("--------", color)
	return colorMatch(color, targetColor)
}

func runWithJump(sec int) {
	robotgo.KeyDown("w")
	robotgo.MilliSleep(200)
	for i := 0; i < sec; i++ {
		if i%2 == 0 {
			robotgo.KeyPress("space")
			robotgo.Sleep(1)
		}
	}
	robotgo.KeyUp("w")
}
