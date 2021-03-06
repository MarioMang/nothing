package main

import (
	"github.com/sciter-sdk/go-sciter"
	"fmt"
	"com/ectongs/preplanui/utils"
	"com/ectongs/preplanui/mq"
	"regexp"
	"com/ectongs/preplanui/gui"
	"com/ectongs/preplanui/conf"
)

const (
	Width  = 600
	Height = 600

	CHECKED   = 0x40000060
	UNCHECKED = 0x40000000
	NOTHING   = 0x40000000
)

func OpenSettingDialog() {
	//创建window窗口
	//参数一表示创建窗口的样式
	//SW_TITLEBAR 顶层窗口，有标题栏
	//SW_RESIZEABLE 可调整大小
	//SW_CONTROLS 有最小/最大按钮
	//SW_MAIN 应用程序主窗口，关闭后其他所有窗口也会关闭
	//SW_ENABLE_DEBUG 可以调试
	//参数二表示创建窗口的矩形

	dialog := gui.NewLoginDialod(gui.EW_TITLEBAR|
		gui.EW_CONTROLS|
		gui.EW_MAIN|
		gui.EW_ENABLE_DEBUG, Width, Height, "设置")
	dialog.SetFile("views\\conf.html")

	//CallbackHandler是一个结构，里面定义了一些方法
	//你可以通过实现这些方法，自定义自已的回调
	cb := &sciter.CallbackHandler{
		//加载数据开始
		OnLoadData: func(p *sciter.ScnLoadData) int {
			//显示加载资源的uri
			fmt.Println("加载:", p.Uri());
			return sciter.LOAD_OK;
		},
		//加载数据过程中
		OnDataLoaded: func(p *sciter.ScnDataLoaded) int {
			fmt.Println("加载中:", p.Uri());
			return sciter.LOAD_OK;
		},
	};
	dialog.SetCallBackHandlers(cb)

	dialog.SetElementHandlers(elementHandlers)

	dialog.Show()
	dialog.Run()
}

func elementHandlers(root *sciter.Element) {

	saveBtn := root.MustSelectById("savebtn");
	saveBtn.OnClick(func() {
		phone := root.MustSelectById("phone")
		elements, err := phone.Select("option")
		if err != nil {
			fmt.Println("setElementHandlers.phone.Select:", err.Error())
		}

		var val *sciter.Value = nil
		for _, elem := range elements {
			state, err := elem.State()
			if err != nil {
				fmt.Println("setElementHandlers.saveBtn.State:", err.Error())
			}

			if state == CHECKED {
				val, err = elem.GetValue()
				if err != nil {
					fmt.Println("setElementHandlers.elem.GetValue:", err.Error())
				}
			}
		}

		var phoneType string
		var serverAddr string
		var inner string
		var outter string
		var pr string

		//设置电话类型
		if val == nil {
			phoneType = conf.PhoneType()
		} else {
			phoneType = val.String()
		}

		//设置服务器地址
		sval, err := root.MustSelectById("server").GetValue()
		if err != nil {
			fmt.Println("setElementHandlers.root.MustSelectById:", err.Error())
		}
		if ok, err := regexp.MatchString("[a-zA-z]+://[^\\s]*", sval.String()); err != nil || !ok {
			utils.MsgBoxWithWarning(nil, "请输入正确的服务器地址")
			serverAddr = conf.ServerAddress()
		} else {
			serverAddr = sval.String()
		}

		//设置市内电话区号
		inner_teleno, err := root.MustSelectById("inner_teleno").GetValue()
		if err != nil {
			fmt.Println("setElementHandlers.root.MustSelectById:", err.Error())
		}
		if str := inner_teleno.String(); str == "" {
			inner = conf.InnerTeleNo()
		} else {
			inner = str
		}

		//设置市外电话区号
		outter_teleno, err := root.MustSelectById("outter_teleno").GetValue()
		if err != nil {
			fmt.Println("setElementHandlers.root.MustSelectById:", err.Error())
		}
		if str := outter_teleno.String(); str == "" {
			outter = conf.OutterTeleNo()
		} else {
			outter = str
		}

		//设置当前城市
		province, err := root.MustSelectById("province").GetValue()
		if err != nil {
			fmt.Println("setElementHandlers.root.MustSelectById:", err.Error())
		}
		if str := province.String(); str == "" {
			pr = conf.Province()
		} else {
			pr = str
		}

		conf := fmt.Sprintf("%s(|=|)%s(|=|)%s(|=|)%s(|=|)%s", phoneType, serverAddr, inner, outter, pr)
		fmt.Println(conf)
		mq.NotifyConf(conf)
	})
}
