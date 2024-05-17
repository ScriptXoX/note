package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
)

// 目标坐标 目标颜色
var targetX, targetY, targetColor = 610, 148, []string{"a", "9"}

var target2X, target2Y = 649, 172

// 战斗状态坐标 战斗状态颜色
var combatStatusX, combatStatusY, combatColor = 421, 919, []string{"f"}

// 血量坐标  血量颜色
var hpX, hpY, hpColor = 798, 904, []string{"9", "a", "b", "c", "d", "e", "f"}

// 是否在水中颜色坐标 和 颜色
var waterX, waterY, waterColor = 1184, 929, []string{"a5", "b9"}

// 定义E R F 三个武器技能释放时间
var indexE, indexR, indexF int64 = time.Now().Unix() - 6, time.Now().Unix() - 9, time.Now().Unix() - 14

// 翠玉机甲是否活着
var jadeX, jadeY, jadeColor = 475, 875, []string{"9f", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9", "a0"}

var count int = 0

var resetTime = time.Now().Unix()

func main() {
	//robotgo.ActivePid(10764)

	//robotgo.ActiveName("Guild Wars 2")
	//robotgo.SetFocus(robotgo.GetHWND())

	gw2bot()

	//fmt.Println(localInit(605))

	//reset()

}

// func init() {
// 	fmt.Println("init")
// 	targetX = localInit(targetX)
// 	targetY = localInit(targetY)

// 	combatStatusX = localInit(combatStatusX)
// 	combatStatusY = localInit(combatStatusY)

// 	hpX = localInit(hpX)
// 	hpY = localInit(hpY)

// 	waterX = localInit(waterX)
// 	waterY = localInit(waterY)

// 	jadeX = localInit(jadeX)
// 	jadeY = localInit(jadeY)

// }

func localInit(local int) int {

	// flocal := float32(local)

	// newLocal := int(flocal / 1)

	return local
}

func SetEnvCount(count int) {
	//os.Setenv("gw2EMCount", strconv.Itoa(count))
	os.WriteFile("gw2Count.txt", []byte(strconv.Itoa(count)), 0666)
}

func getEnvCount() int {

	content, _ := os.ReadFile("gw2Count.txt")
	count, error := strconv.Atoi(string(content))

	if error != nil {
		count = 0
	}
	return count
}

func gw2bot() {

	//等待2秒 进入激活窗口
	//robotgo.Sleep(2)

	robotgo.Sleep(2)

	if count != 0 {
		SetEnvCount(count)
	}

	count = getEnvCount()
	fmt.Println("打怪数量统计：", count)

	// w   mkw   mjw

	for i := 0; ; i++ {

		fmt.Println("--------------------重置剩余时间:", 800-(time.Now().Unix()-resetTime))
		if time.Now().Unix()-resetTime > 800 {
			fmt.Println("---------------------位置重置---------------------")
			reset()
		}

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
	// if count%10 == 0 {
	// 	robotgo.Sleep(10)
	// }
	//处理在水中
	processInWater()

	fmt.Println("自动寻怪")
	robotgo.KeyPress("m")
	robotgo.MilliSleep(50)

	if isGetTarget() {
		fmt.Println("找到目标怪")
	} else {
		fmt.Println("没有找到目标怪 换方向 向前走")

		// robotgo.KeyDown("alt")
		// robotgo.KeyPress("w")
		// robotgo.KeyUp("alt")

		//换方向
		turn(600)

		robotgo.MilliSleep(50)

		//跳着向前走
		runWithJump(6)

	}

}

// 尝试攻击
func tryAttack() bool {

	if !isGetTarget() {
		return false
	}

	for i := 0; i < 3; i++ {

		processInWater()

		if !isGetTarget() {
			return false
		}

		robotgo.KeyPress("q")
		robotgo.MilliSleep(100)
		robotgo.KeyPress("6")
		robotgo.MilliSleep(100)

		fmt.Println("调正视角")
		robotgo.MilliSleep(100)
		robotgo.KeyDown("j")
		robotgo.MilliSleep(100)
		robotgo.KeyUp("j")

		robotgo.MilliSleep(100)
		robotgo.KeyDown("k")
		robotgo.MilliSleep(100)
		robotgo.KeyUp("k")
		robotgo.MilliSleep(100)

		robotgo.MilliSleep(1000)

		fmt.Println("检测是否在战斗状态")

		if !isCombat() && isGetTarget() {
			fmt.Println("尝试攻击没有达到怪，尝试")
			runWithJump(3)
			//tryAttack()
		} else {
			//如果打到怪了退出 进行下一步
			return true
		}

	}

	robotgo.MilliSleep(100)
	if isGetTarget() && !isCombat() {
		fmt.Println("多次尝试没达到怪 放弃 掉头")

		robotgo.MilliSleep(100)
		robotgo.Move(100, 100)
		robotgo.MilliSleep(100)
		robotgo.Click()

		robotgo.KeyDown("alt")
		robotgo.MilliSleep(100)
		robotgo.KeyPress("w")
		robotgo.KeyUp("alt")

		runWithJump(6)

		turn(700)

		runWithJump(6)

	}
	return false

}

// 进行攻击
func doAttack() {
	for i := 0; i < 50; i++ {

		processHP()

		processInWater()

		processJade()

		robotgo.MilliSleep(1000)

		if !isGetTarget() {
			fmt.Println("目标消失,停止打怪")
			count++
			SetEnvCount(count)

			fmt.Println("打怪数量:--------------------------------------------------", count)
			break
		}

		if !isCombat() {
			fmt.Println("脱战,停止打怪")
			count++
			SetEnvCount(count)
			fmt.Println("打怪数量:--------------------------------------------------", count)
			break
		}

		nowSec := time.Now().Unix()
		//fmt.Println("E", nowSec, indexE)
		if nowSec-indexE > 6 {
			fmt.Println("放E 技能")
			robotgo.KeyPress("e")
			robotgo.MilliSleep(500)

			indexE = time.Now().Unix()
		}

		if nowSec-indexR > 9 {
			fmt.Println("放R技能")
			robotgo.KeyPress("r")
			robotgo.MilliSleep(500)
			indexR = time.Now().Unix()

		}

		if nowSec-indexF > 14 {
			fmt.Println("放F 技能")
			robotgo.KeyPress("f")
			robotgo.MilliSleep(500)
			indexF = time.Now().Unix()

		}

	}

}

// 机甲处理
func processJade() {
	color := robotgo.GetPixelColor(jadeX, jadeY, 0)
	color = color[0:2]
	fmt.Println(color)
	if !colorMatch(color, jadeColor) {
		fmt.Println("机甲死了 重新召唤")

		robotgo.Sleep(1)
		robotgo.KeyPress("G")
		robotgo.MilliSleep(700)
		robotgo.KeyPress("space")
		robotgo.MilliSleep(200)
		robotgo.KeyPress("space")
		robotgo.MilliSleep(200)
		robotgo.KeyPress("space")

		robotgo.Sleep(1)
		robotgo.KeyPress("b")
		robotgo.Sleep(1)

	}
}

// 处理在水中
func processInWater() {
	color := robotgo.GetPixelColor(waterX, waterY, 0)
	color = color[0:2]
	fmt.Println("水中颜色:", color)
	if colorMatch(color, waterColor) {
		fmt.Println("走到水中去了，掉头")

		turn(1000)

		robotgo.KeyDown("w")
		robotgo.MilliSleep(200)
		robotgo.KeyDown("space")

		robotgo.Sleep(10)

		robotgo.KeyUp("space")
		robotgo.KeyUp("w")

		runWithJump(15)

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
		robotgo.MilliSleep(1000)

		//翻滚
		robotgo.KeyPress("v")
		robotgo.MilliSleep(1000)

		//加血

		robotgo.KeyPress("t")
		robotgo.MilliSleep(850)

		robotgo.KeyDown("alt")
		robotgo.KeyPress("w")
		robotgo.KeyUp("alt")

		turn(700)

		runWithJump(8)

		turn(700)

		runWithJump(8)

		robotgo.Move(100, 100)
		robotgo.Click()

		robotgo.MilliSleep(300)

		robotgo.KeyPress("m")
	}
}

func isCombat() bool {

	combatc := ""

	combatcolor := robotgo.GetPixelColor(combatStatusX, combatStatusY, 0)
	combatc = combatcolor[0:1]
	fmt.Println("战斗状态颜色", combatcolor, "-----", combatc)
	if colorMatch(combatc, combatColor) {
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

func isGetTarget1(targetX int, targetY int, targetColor []string) bool {
	color := robotgo.GetPixelColor(targetX, targetY, 0)
	fmt.Print("颜色:", color)
	color = color[0:1]
	fmt.Println("--------", color)
	//robotgo.Move(targetX, targetY)
	return colorMatch(color, targetColor)
}

func isGetTarget() bool {

	return isGetTarget1(target2X, target2Y, targetColor) || isGetTarget1(targetX, targetY, targetColor)
}

func runWithJump(sec int) {

	robotgo.MilliSleep(200)
	robotgo.KeyDown("w")
	robotgo.MilliSleep(200)
	for i := 0; i < sec; i++ {
		if i%2 == 0 {
			robotgo.KeyPress("space")

		}
		robotgo.Sleep(1)
	}
	robotgo.KeyUp("w")
}

func turn(msec int) {
	robotgo.MilliSleep(100)

	direction := "j"
	if time.Now().Second()%2 == 0 {
		direction = "k"
	}
	robotgo.KeyDown(direction)
	robotgo.MilliSleep(msec)
	robotgo.KeyUp(direction)
}

func reset() {
	robotgo.Sleep(6)
	robotgo.KeyPress("enter")
	robotgo.TypeStr("[&BE4DAAA=]")
	robotgo.KeyPress("enter")

	if isCombat() {
		fmt.Println("---------------在战斗中 取消重置---------------------")
		return
	}

	robotgo.MilliSleep(200)
	robotgo.Move(localInit(378), localInit(771))
	robotgo.MilliSleep(200)
	robotgo.Click()

	robotgo.MilliSleep(2000)

	robotgo.Move(localInit(804), localInit(515))
	robotgo.Sleep(1)
	robotgo.Click()
	robotgo.MilliSleep(100)
	robotgo.Click()

	robotgo.Sleep(3)

	c := robotgo.GetPixelColor(localInit(57), localInit(167))
	fmt.Println("大地图判定颜色", c)
	c = c[0:1]
	if colorMatch(c, []string{"c", "d"}) {
		robotgo.KeyPress("esc")
	}

	robotgo.Sleep(3)

	robotgo.KeyUp("w")
	robotgo.MilliSleep(100)
	robotgo.KeyDown("w")
	runWithJump(20)
	robotgo.KeyUp("w")

	resetTime = time.Now().Unix()

}
