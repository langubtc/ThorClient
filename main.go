package main

import (
	"ThorGui/config"
	"ThorGui/thormonitor"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// 创建自定义窗体
type MyMainWindow struct {
	*walk.MainWindow
}

// 矿池设置结构体
type MinerPool struct {
	Worker  string
	PoolStr string
	Wallet  string
}

type Foo struct {
	Index       int
	Worker      string
	MinerIp     string
	Wallet      string
	MinerStatus string
	ServerIp    string
	ServerStr   string
	ServerPort  string
	ScanTime    string
}

type FooModel struct {
	walk.ReflectTableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.ResetRows()
	return m
}

func (m *FooModel) Items() interface{} {
	return m.items
}

func (m *FooModel) GetByIndex(index int64) *Foo {
	return m.items[index]
}

func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {

		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.MinerIp < b.MinerIp)
		case 2:
			return c(a.Wallet < b.Wallet)

		case 3:
			return c(a.Worker < b.Worker)
		case 4:
			return c(a.ServerStr < b.ServerStr)
		case 5:
			return c(a.MinerStatus < b.MinerStatus)
		case 6:
			return c(a.ServerIp < b.ServerIp)
		case 7:
			return c(a.ServerPort < b.ServerPort)
		case 8:
			return c(a.ScanTime < b.ScanTime)

		}
		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func jsonRead() []Foo {
	file, err := ioutil.ReadFile("./config/result.json")
	if err != nil {
		log.Fatal("err")
	}

	var foos []Foo
	err = json.Unmarshal(file, &foos)
	if err != nil {
		fmt.Println(err.Error())

	}

	return foos
}

func searchMiner(search string) []Foo {
	file, err := ioutil.ReadFile("./config/result.json")
	if err != nil {
		log.Fatal("err")
	}

	var foos []Foo
	var saerch []Foo

	err = json.Unmarshal(file, &foos)
	if err != nil {
		fmt.Println(err.Error())

	}

	//判断查询条件，筛选IP或者矿工号
	for i := 0; i < len(foos); i++ {
		if foos[i].Worker == search || foos[i].MinerIp == search {
			saerch = append(saerch, foos[i])
		} else {
			continue
		}
	}

	return saerch
}

func (m *FooModel) returnData() {

	m.items = make([]*Foo, 1000)

	miners := jsonRead()

	for i := 0; i < len(miners); i++ {
		var tminer Foo
		tminer = miners[i]
		m.items[i] = &tminer
	}

	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)

}

func (m *FooModel) ResetRows() {

	miners := jsonRead()
	m.items = make([]*Foo, len(miners))

	for i := 0; i < len(miners); i++ {
		var tminer Foo
		tminer = miners[i]
		m.items[i] = &tminer

	}

	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)

}

func (m *FooModel) filerRows(miner string) {

	//now := time.Now()

	miners := searchMiner(miner)
	m.items = make([]*Foo, len(miners))

	for i := 0; i < len(miners); i++ {
		var tminer Foo
		tminer = miners[i]
		m.items[i] = &tminer

	}

	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)

}

func MinerConfigJson() *MinerPool {

	file, err := ioutil.ReadFile("./config/minerConfig.json")
	if err != nil {
		log.Fatal("err")
	}

	var miner MinerPool
	err = json.Unmarshal(file, &miner)
	if err != nil {
		fmt.Println(err.Error())

	}

	return &miner

}

func MinerRunDetail(owner walk.Form, detail string) {
	var outTE *walk.TextEdit

	replacer := strings.NewReplacer("\033[0m", "",
		"\033[36m", "",
		"\033[91m", "",
		"\033[32m", "",
		"\033[97m", "",
		"\033[94m", "",
		"\033[33m", "",
		"\033[35m", "",
		"\033[34m", "",
		"\033[1;97m", "",
		"\033[1;36m", "",
	)
	newdetail := replacer.Replace(detail)
	fmt.Println("----------------------------------------------")
	fmt.Println(newdetail)

	var dialog = Dialog{}

	dialog.Title = "矿机详情"
	dialog.MinSize = Size{900, 800}
	dialog.Layout = VBox{}
	childrens := []Widget{
		HSplitter{
			Children: []Widget{
				TextEdit{
					Background: SolidColorBrush{Color: walk.RGB(0, 0, 0)},
					Font: Font{
						PointSize: 11,
					},
					TextColor: walk.RGB(255, 255, 255),
					AssignTo:  &outTE,
					ReadOnly:  true,
					Text:      newdetail,
					VScroll:   true,
				},
			},
		},
	}
	dialog.Children = childrens

	dialog.Run(owner)

}

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("revocer....")
		}
	}()

	//var inTE, outTE *walk.TextEdit // 声明两个文本编辑控件
	var tv *walk.TableView
	var usernameTE *walk.LineEdit
	var ScanNet *walk.LineEdit
	var poolIp *walk.LineEdit
	var Miner *walk.LineEdit
	var WorkerName *walk.LineEdit
	var result *walk.Label

	// 定义获取旷工配置
	conf_miner := MinerConfigJson()

	//全局配置
	sys_config := config.LoadConfig()

	font := Font{
		//Family: "Times New Roman",
		PointSize: 10,
		//Bold: true,
	}

	//new一个自定义窗体的指针
	mw := new(MyMainWindow)

	model := NewFooModel()

	titleStr := fmt.Sprintf("THOR矿机管理 %s     [当前工作账户%s]", sys_config.Version, sys_config.User)
	versionStr := fmt.Sprintf("当前版本：%s", sys_config.Version)

	//主窗口对象
	MainWindow{
		AssignTo: &mw.MainWindow, //赋值给mw
		Font:     font,
		Icon:     "icon/thor.ico", //设置Icon
		Title:    titleStr,        // 窗口标题设置
		Size:     Size{1200, 800}, //窗体的大小
		MinSize:  Size{600, 600},
		Layout:   VBox{}, // 窗体的布局形式

		// 定义菜单
		MenuItems: []MenuItem{
			Menu{
				Text: "&系统",
				Items: []MenuItem{
					Action{
						Text:        "退出",
						OnTriggered: func() { mw.Close() },
					},
				},
			},

			Menu{
				Text: "&关于",
				Items: []MenuItem{
					Action{
						Text: "版本",
						OnTriggered: func() {
							walk.MsgBox(mw, "版本", versionStr, walk.MsgBoxIconInformation)
						},
					},
					Action{
						Text: "联系我们",
						OnTriggered: func() {
							walk.MsgBox(mw, "联系我们", "开发者:蓝谷骑兵\nQQ:923401910", walk.MsgBoxIconInformation)
						},
					},
				},
			},
		},

		//定义vbox的所有控件
		Children: []Widget{ //定义控件
			VSplitter{
				Children: []Widget{
					GroupBox{
						Layout: Grid{Columns: 6, Spacing: 20},

						Children: []Widget{
							Label{
								Text: "扫描网段:",
							},
							LineEdit{
								MaxSize:     Size{200, 30},
								AssignTo:    &ScanNet,
								Text:        sys_config.IpRange,
								ToolTipText: "输入网段样例:192.168.1.0/24",
							},

							PushButton{ //按钮控件
								Text:    "开始扫描",
								MaxSize: Size{200, 80},
								Image:   "./icon/scan.png",
								OnClicked: func() {
									//定义个空数组用户传参
									//var iplist [] string

									if ScanNet.Text() != "" {

										iplist, err := thormonitor.MinerIPFunc(ScanNet.Text())
										if err != nil {
											walk.MsgBox(mw, "提示", "IP范围错误!", walk.MsgBoxIconError)
										}
										total, active, inactive, _ := thormonitor.RunMonitor("scan", iplist, "", "")

										resultStr := fmt.Sprintf("%d/%d/%d", total, active, inactive)

										model.ResetRows()

										fmt.Println(resultStr)

										if err != nil {
											walk.MsgBox(mw, "提示", "系统错误!", walk.MsgBoxIconError)
										} else {
											walk.MsgBox(mw, "提示", "扫描完成", walk.MsgBoxIconInformation)
										}
										err = result.SetText(resultStr)
									} else {
										walk.MsgBox(mw, "提示", "网段错误!", walk.MsgBoxIconError)
									}

								},
							},

							Label{
								ToolTipText: "扫描结果",
								Text:        "扫描总数/成功/失败:",
							},

							Label{
								Text:     "0/0/0",
								Name:     "result",
								MinSize:  Size{50, 30},
								AssignTo: &result,
							},
						},
					},
					GroupBox{
						Layout: Grid{Columns: 4, Spacing: 10},
						Children: []Widget{
							Label{
								Text: "矿池接入IP:",
							},
							LineEdit{
								MinSize:  Size{160, 0},
								MaxSize:  Size{200, 30},
								AssignTo: &poolIp,
								Text:     conf_miner.PoolStr,
							},
							Label{
								Text: "矿工钱包:",
							},
							LineEdit{
								MinSize:  Size{160, 0},
								MaxSize:  Size{400, 30},
								AssignTo: &Miner,
								Text:     conf_miner.Wallet,
							},
							Label{
								Text: "矿工名:",
							},
							LineEdit{
								MinSize:  Size{160, 0},
								MaxSize:  Size{200, 30},
								AssignTo: &WorkerName,
								Text:     conf_miner.Worker,
							},
							PushButton{ //按钮控件
								Text:    "保存",
								MaxSize: Size{50, 30},
								OnClicked: func() {

									//保存设置到JSON文件
									data := MinerPool{
										Worker:  WorkerName.Text(),
										Wallet:  Miner.Text(),
										PoolStr: poolIp.Text(),
									}

									file, _ := json.MarshalIndent(data, "", "")
									_ = ioutil.WriteFile("./config/minerConfig.json", file, 0644)
									walk.MsgBox(mw, "提示", "设置成功", walk.MsgBoxIconInformation)

								},
							},
						},
					},

					GroupBox{
						Layout: Grid{Columns: 8, Spacing: 0},
						Children: []Widget{
							LineEdit{
								MinSize:     Size{160, 0},
								MaxSize:     Size{200, 30},
								AssignTo:    &usernameTE,
								Text:        "",
								ToolTipText: "输入IP或者矿工号",
							},
							PushButton{ //按钮控件
								Text:    "查询",
								Image:   "./icon/search.png",
								MaxSize: Size{50, 30},
								OnClicked: func() {
									//如果为空就获取所有JSON数据
									if usernameTE.Text() == "" {
										model.ResetRows()
									} else {
										model.filerRows(usernameTE.Text())
									}
								},
							},

							PushButton{
								Text:  "修改",
								Image: "./icon/edit.png",
								OnClicked: func() {
									var iplist []string
									indexs := tv.SelectedIndexes()
									if len(indexs) == 0 {
										walk.MsgBox(mw, "提示", "请选择矿机!", walk.MsgBoxIconError)
										return
									}

									// 遍历选择的矿机进行重启
									for i := 0; i < len(indexs); i++ {

										var minerIpList []string
										itemValue := model.GetByIndex(int64(indexs[i]))
										iplist = append(iplist, itemValue.MinerIp)

										//分台设置配置
										minerIpList = append(minerIpList, itemValue.MinerIp)

										fmt.Println(itemValue.ServerStr, poolIp.Text())
										if poolIp.Text() != "" {
											_, _, _, _ = thormonitor.RunMonitor("update", iplist, itemValue.ServerStr, poolIp.Text())
										}
										if Miner.Text() != "" {
											_, _, _, _ = thormonitor.RunMonitor("update", iplist, itemValue.Wallet, Miner.Text())
										}
										if WorkerName.Text() != "" {
											_, _, _, _ = thormonitor.RunMonitor("update", iplist, itemValue.Worker, WorkerName.Text())
										}

									}

									//批量重启
									total, active, inactive, _ := thormonitor.RunMonitor("reboot", iplist, "", "")
									rebootMiner := fmt.Sprintf("设置成功,选择%d台矿机,成功%d,失败%d", total, active, inactive)

									walk.MsgBox(mw, "提示", rebootMiner, walk.MsgBoxIconInformation)
								},
							},

							PushButton{
								Text:  "重启",
								Image: "./icon/restart.png",
								OnClicked: func() {

									var iplist []string
									indexs := tv.SelectedIndexes()
									if len(indexs) == 0 {
										walk.MsgBox(mw, "提示", "请选择矿机!", walk.MsgBoxIconError)
										return
									}

									// 遍历选择的矿机进行重启
									for i := 0; i < len(indexs); i++ {
										itemValue := model.GetByIndex(int64(indexs[i]))
										iplist = append(iplist, itemValue.MinerIp)
									}

									//重启结果回显
									total, active, inactive, _ := thormonitor.RunMonitor("reboot", iplist, "", "")

									rebootMiner := fmt.Sprintf("重启成功,选择%d台矿机,成功%d,失败%d", total, active, inactive)

									walk.MsgBox(mw, "提示", rebootMiner, walk.MsgBoxIconInformation)
								},
							},

							PushButton{
								Text:  "查看状态",
								Image: "./icon/refresh.png",
								OnClicked: func() {
									var iplist []string

									indexs := tv.SelectedIndexes()
									if len(indexs) == 0 {
										walk.MsgBox(mw, "提示", "请选择矿机!", walk.MsgBoxIconError)
										return
									}
									//判断选择是否是单选
									if len(indexs) > 1 {
										walk.MsgBox(mw, "提示", "只能单选一台矿机!", walk.MsgBoxIconError)
									} else {

										// 遍历选择的矿机进行重启
										for i := 0; i < len(indexs); i++ {
											itemValue := model.GetByIndex(int64(indexs[i]))
											iplist = append(iplist, itemValue.MinerIp)
										}

										_, _, _, statsResult := thormonitor.RunMonitor("stats", iplist, "", "")
										fmt.Println(statsResult)
										MinerRunDetail(mw, statsResult)
									}
								},
							},
						},
					},
					GroupBox{
						Layout: VBox{},
						Children: []Widget{

							TableView{
								AssignTo:         &tv,
								AlternatingRowBG: true,
								CheckBoxes:       false, //开启前面的勾选框
								CustomRowHeight:  25,    // 自定义表格行高
								ColumnsOrderable: true,
								MultiSelection:   true,
								Columns: []TableViewColumn{

									{Title: "#", Name: "Index", Width: 50},
									{Title: "矿机IP", Width: 100, Name: "MinerIp"},
									{Title: "矿工钱包", Width: 350, Name: "Wallet"},
									{Title: "矿工号", Name: "Worker", Width: 50},
									{Title: "矿池接入", Alignment: AlignFar, Width: 150, Name: "ServerStr"},
									{Title: "矿池IP", Alignment: AlignFar, Width: 100, Name: "ServerIp"},
									{Title: "状态", Alignment: AlignFar, Width: 100, Name: "MinerStatus"},
									{Title: "端口", Alignment: AlignFar, Name: "ServerPort", Width: 50},
									{Title: "扫描时间", Format: "2006-01-02 15:04:05", Width: 150, Name: "ScanTime"},
								},
								Model: model,

								OnSelectedIndexesChanged: func() {
									fmt.Printf("selected:%v\n", tv.SelectedIndexes())
								},
							},
							// 混合样式
							Composite{
								Layout: HBox{},
								Children: []Widget{
									Label{
										ToolTipText: "版本",
										Text:        "Build Date: 2020-12-01  By: 蓝谷骑兵",
									},
								},
							},
						},
					},
				},
			},
		},
	}.Create()

	mw.Run()

}
