package cfg

import (
	"fmt"
	"image"
	_ "image/png" // 注册 PNG 解码器
	"os"
	"path/filepath"
	"saClient/src/common"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// UIAnimation 动画配置
type UIAnimation struct {
	Speed        float32 `yaml:"speed"`        // 动画速度（像素/帧）
	StaggerDelay int     `yaml:"staggerDelay"` // 依次弹出延迟（帧）
	Easing       string  `yaml:"easing"`       // 缓动函数
}

// UIButtonStyle 按钮样式配置
type UIButtonStyle struct {
	Width   int `yaml:"width"`   // 按钮宽度
	Height  int `yaml:"height"`  // 按钮高度
	Spacing int `yaml:"spacing"` // 按钮间距
	Margin  int `yaml:"margin"`  // 边缘距离
}

// UIGlobal 全局UI配置
type UIGlobal struct {
	Animation   UIAnimation   `yaml:"animation"`
	ButtonStyle UIButtonStyle `yaml:"buttonStyle"`
}

// UIMenuAnimation 菜单动画配置
type UIMenuAnimation struct {
	OpenType     string  `yaml:"openType"`     // 展开类型: stagger, simultaneous
	CloseType    string  `yaml:"closeType"`    // 收起类型: stagger, simultaneous
	Speed        float32 `yaml:"speed"`        // 动画速度
	StaggerDelay int     `yaml:"staggerDelay"` // 延迟帧数
}

// UIToggleButton 切换按钮配置
type UIToggleButton struct {
	ExpandIcon   string `yaml:"expandIcon"`   // 展开图标路径
	CollapseIcon string `yaml:"collapseIcon"` // 收起图标路径
	Width        int    `yaml:"width"`        // 宽度
	Height       int    `yaml:"height"`       // 高度

	// 运行时加载的资源
	ExpandImg   *ebitenv2.Image `yaml:"-"`
	CollapseImg *ebitenv2.Image `yaml:"-"`
}

// UIMenuItemRef 菜单项引用
type UIMenuItemRef struct {
	ID uint32 `yaml:"id"`

	// 运行时关联的元素
	Element *UIElement `yaml:"-"`
}

// UIPopupMenu 弹出菜单配置
type UIPopupMenu struct {
	ID              uint32          `yaml:"id"`
	Name            string          `yaml:"name"`
	Describe        string          `yaml:"describe"`
	Anchor          string          `yaml:"anchor"`          // 锚点: topLeft, topRight, bottomLeft, bottomRight, center
	OffsetX         int             `yaml:"offsetX"`         // X偏移
	OffsetY         int             `yaml:"offsetY"`         // Y偏移
	ExpandDirection string          `yaml:"expandDirection"` // 展开方向: left, right, up, down
	Animation       UIMenuAnimation `yaml:"animation"`
	ToggleButton    UIToggleButton  `yaml:"toggleButton"`
	Items           []UIMenuItemRef `yaml:"items"`
}

// UIElement UI元素定义
type UIElement struct {
	ID       uint32 `yaml:"id"`
	Type     string `yaml:"type"`     // icon, menuItem, panel
	Name     string `yaml:"name"`     // 名称
	Describe string `yaml:"describe"` // 描述
	IconPath string `yaml:"iconPath"` // 图标路径
	Action   string `yaml:"action"`   // 点击动作标识
	Hotkey   string `yaml:"hotkey"`   // 快捷键

	// 运行时加载的资源
	Icon *ebitenv2.Image `yaml:"-"`
}

var GUIMgr = newUIMgr()

// UIMgr UI配置管理器
type UIMgr struct {
	Global     UIGlobal
	PopupMenus *xmap.MapMgr[uint32, *UIPopupMenu]
	Elements   *xmap.MapMgr[uint32, *UIElement]
}

func newUIMgr() *UIMgr {
	return &UIMgr{
		PopupMenus: xmap.NewMapMgr[uint32, *UIPopupMenu](),
		Elements:   xmap.NewMapMgr[uint32, *UIElement](),
	}
}

// Load 加载配置文件
func (p *UIMgr) Load() error {
	cfgPath := filepath.Join(common.AppCfgDir, "ui.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "读取UI配置文件失败: %v %v", cfgPath, xruntime.Location())
	}

	var config struct {
		Global     UIGlobal       `yaml:"global"`
		PopupMenus []*UIPopupMenu `yaml:"popupMenus"`
		Elements   []*UIElement   `yaml:"elements"`
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return errors.WithMessagef(err, "解析UI配置文件失败: %v %v", cfgPath, xruntime.Location())
	}

	// 加载全局配置
	p.Global = config.Global

	// 设置全局默认值
	if p.Global.Animation.Speed == 0 {
		p.Global.Animation.Speed = 12.0
	}
	if p.Global.Animation.StaggerDelay == 0 {
		p.Global.Animation.StaggerDelay = 5
	}
	if p.Global.Animation.Easing == "" {
		p.Global.Animation.Easing = "linear"
	}
	if p.Global.ButtonStyle.Width == 0 {
		p.Global.ButtonStyle.Width = 60
	}
	if p.Global.ButtonStyle.Height == 0 {
		p.Global.ButtonStyle.Height = 60
	}
	if p.Global.ButtonStyle.Spacing == 0 {
		p.Global.ButtonStyle.Spacing = 10
	}
	if p.Global.ButtonStyle.Margin == 0 {
		p.Global.ButtonStyle.Margin = 20
	}

	// 加载UI元素
	for _, elem := range config.Elements {
		if elem.ID == 0 {
			return fmt.Errorf("UI元素ID不能为0 %v", xruntime.Location())
		}
		ok := p.Elements.AddIfNotExist(elem.ID, elem)
		if !ok {
			return fmt.Errorf("UI元素ID重复: %d %v", elem.ID, xruntime.Location())
		}
	}

	// 加载弹出菜单
	for _, menu := range config.PopupMenus {
		if menu.ID == 0 {
			return fmt.Errorf("弹出菜单ID不能为0 %v", xruntime.Location())
		}
		// 验证锚点
		if !isValidAnchor(menu.Anchor) {
			return fmt.Errorf("弹出菜单ID %d 锚点无效: %s %v", menu.ID, menu.Anchor, xruntime.Location())
		}
		// 验证展开方向
		if !isValidDirection(menu.ExpandDirection) {
			return fmt.Errorf("弹出菜单ID %d 展开方向无效: %s %v", menu.ID, menu.ExpandDirection, xruntime.Location())
		}
		// 设置动画默认值
		if menu.Animation.Speed == 0 {
			menu.Animation.Speed = p.Global.Animation.Speed
		}
		if menu.Animation.StaggerDelay == 0 {
			menu.Animation.StaggerDelay = p.Global.Animation.StaggerDelay
		}
		if menu.Animation.OpenType == "" {
			menu.Animation.OpenType = "stagger"
		}
		if menu.Animation.CloseType == "" {
			menu.Animation.CloseType = "simultaneous"
		}
		// 设置切换按钮默认值
		if menu.ToggleButton.Width == 0 {
			menu.ToggleButton.Width = p.Global.ButtonStyle.Width
		}
		if menu.ToggleButton.Height == 0 {
			menu.ToggleButton.Height = p.Global.ButtonStyle.Height
		}

		ok := p.PopupMenus.AddIfNotExist(menu.ID, menu)
		if !ok {
			return fmt.Errorf("弹出菜单ID重复: %d %v", menu.ID, xruntime.Location())
		}
	}

	return nil
}

// Check 检查配置
func (p *UIMgr) Check() error {
	var err error

	// 检查弹出菜单引用的元素是否存在
	p.PopupMenus.Foreach(func(id uint32, menu *UIPopupMenu) bool {
		for i := range menu.Items {
			itemRef := &menu.Items[i]
			elem := p.Elements.Get(itemRef.ID)
			if elem == nil {
				err = fmt.Errorf("弹出菜单ID %d 引用的元素ID %d 不存在 %v", menu.ID, itemRef.ID, xruntime.Location())
				return false
			}
			itemRef.Element = elem
		}
		return true
	})

	return err
}

// Assemble 组装配置（加载图片资源）
func (p *UIMgr) Assemble() error {
	var err error

	// 加载UI元素图标
	p.Elements.Foreach(func(id uint32, elem *UIElement) bool {
		if elem.IconPath != "" {
			img, loadErr := loadImage(elem.IconPath)
			if loadErr != nil {
				err = fmt.Errorf("UI元素ID %d 图标加载失败: %s %v", elem.ID, elem.IconPath, loadErr)
				return false
			}
			elem.Icon = img
		}
		return true
	})
	if err != nil {
		return err
	}

	// 加载弹出菜单切换按钮图标
	p.PopupMenus.Foreach(func(id uint32, menu *UIPopupMenu) bool {
		if menu.ToggleButton.ExpandIcon != "" {
			img, loadErr := loadImage(menu.ToggleButton.ExpandIcon)
			if loadErr != nil {
				err = fmt.Errorf("弹出菜单ID %d 展开图标加载失败: %s %v", menu.ID, menu.ToggleButton.ExpandIcon, loadErr)
				return false
			}
			menu.ToggleButton.ExpandImg = img
		}
		if menu.ToggleButton.CollapseIcon != "" {
			img, loadErr := loadImage(menu.ToggleButton.CollapseIcon)
			if loadErr != nil {
				err = fmt.Errorf("弹出菜单ID %d 收起图标加载失败: %s %v", menu.ID, menu.ToggleButton.CollapseIcon, loadErr)
				return false
			}
			menu.ToggleButton.CollapseImg = img
		}
		return true
	})

	return err
}

// CalcAnchorPosition 计算锚点位置
func (p *UIMgr) CalcAnchorPosition(anchor string, offsetX, offsetY int) (x, y int) {
	screenW := GCommon.ScreenMaxWidth
	screenH := GCommon.ScreenMaxHeight

	switch anchor {
	case "topLeft":
		return offsetX, offsetY
	case "topRight":
		return screenW + offsetX, offsetY
	case "bottomLeft":
		return offsetX, screenH + offsetY
	case "bottomRight":
		return screenW + offsetX, screenH + offsetY
	case "center":
		return screenW/2 + offsetX, screenH/2 + offsetY
	default:
		return offsetX, offsetY
	}
}

// isValidAnchor 验证锚点是否有效
func isValidAnchor(anchor string) bool {
	switch anchor {
	case "topLeft", "topRight", "bottomLeft", "bottomRight", "center":
		return true
	default:
		return false
	}
}

// isValidDirection 验证方向是否有效
func isValidDirection(dir string) bool {
	switch dir {
	case "left", "right", "up", "down":
		return true
	default:
		return false
	}
}

// loadImage 加载图片
func loadImage(path string) (*ebitenv2.Image, error) {
	// path 格式如 "res/ui/icon.menu.expand.png"，相对于 client 目录
	// AppCfgDir = client/cfg，所以 .. 回到 client 目录
	fullPath := filepath.Join(common.AppCfgDir, "..", path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "打开图片失败: %s", fullPath)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.WithMessagef(err, "解码图片失败: %s", fullPath)
	}

	return ebitenv2.NewImageFromImage(img), nil
}
