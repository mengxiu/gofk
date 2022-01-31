package utils

import (
	"github.com/go-vgo/robotgo"
	"math"
	"reflect"
)

func init() {
	MethodMaps = make(MethodMapsType, 0)
	tweening := tweenings{}
	value := reflect.ValueOf(&tweening)
	vft := value.Type()
	mNum := value.NumMethod()

	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		MethodMaps[mName] = value.Method(i)

	}

}

type tweenings struct {
}
type MethodMapsType map[string]reflect.Value

var MethodMaps MethodMapsType

//methodListAll已经排除Linear,即直线方法
var (
	methodListAll  = []string{"EaseInBack", "EaseInBounce", "EaseInCirc", "EaseInCubic", "EaseInElastic", "EaseInExpo", "EaseInOutBack", "EaseInOutBounce", "EaseInOutCirc", "EaseInOutCubic", "EaseInOutElastic", "EaseInOutExpo", "EaseInOutQuad", "EaseInOutQuart", "EaseInOutQuint", "EaseInOutSine", "EaseInQuad", "EaseInQuart", "EaseInQuint", "EaseInSine", "EaseOutBack", "EaseOutBounce", "EaseOutCirc", "EaseOutCubic", "EaseOutElastic", "EaseOutExpo", "EaseOutQuad", "EaseOutQuart", "EaseOutQuint", "EaseOutSine"}
	methodListShot = []string{}
)

func (self *tweenings) checkRange(x float64) {
	if x < 0.0 || x > 1.0 {
		panic("Argument 'x' must be between 0.0 and 1.0")
	}

}
func (self *tweenings) Linear(x float64) float64 {
	self.checkRange(x)
	return x
}
func (self *tweenings) EaseInQuad(x float64) float64 {
	self.checkRange(x)
	return x * x
}
func (self *tweenings) EaseOutQuad(x float64) float64 {
	self.checkRange(x)
	return x
}
func (self *tweenings) EaseInOutQuad(x float64) float64 {
	self.checkRange(x)
	if x < 0.5 {
		return 2 * math.Pow(x, 2)
	}
	x = x*2 - 1
	return -0.5 * (x*(x-2) - 1)
}
func (self *tweenings) EaseInCubic(x float64) float64 {
	self.checkRange(x)
	return math.Pow(x, 3)
}
func (self *tweenings) EaseOutCubic(x float64) float64 {
	self.checkRange(x)
	x = x - 1
	return math.Pow(x, 3) + 1
}
func (self *tweenings) EaseInOutCubic(x float64) float64 {
	self.checkRange(x)
	x = x * 2
	if x < 1 {
		return 0.5 * math.Pow(x, 3)

	} else {
		x = x - 2
		return 0.5 * (math.Pow(x, 3) + 2)
	}

}
func (self *tweenings) EaseInQuart(x float64) float64 {
	self.checkRange(x)
	return math.Pow(x, 4)
}
func (self *tweenings) EaseOutQuart(x float64) float64 {
	self.checkRange(x)
	x -= 1
	return -(math.Pow(x, 4-1))
}
func (self *tweenings) EaseInOutQuart(x float64) float64 {
	self.checkRange(x)
	x *= 2
	if x < 1 {
		return 0.5 * math.Pow(x, 4)
	}
	x -= 2
	return -0.5 * (math.Pow(x, 4) - 2)
}
func (self *tweenings) EaseInQuint(x float64) float64 {
	self.checkRange(x)
	return math.Pow(x, 5)
}
func (self *tweenings) EaseOutQuint(x float64) float64 {
	self.checkRange(x)
	x -= 1
	return math.Pow(x, 5) + 1
}
func (self *tweenings) EaseInOutQuint(x float64) float64 {
	self.checkRange(x)
	x *= 2
	if x < 1 {
		return 0.5 * math.Pow(x, 5)
	}
	x -= 2
	return 0.5 * (math.Pow(x, 5) + 2)
}
func (self *tweenings) EaseInSine(x float64) float64 {
	self.checkRange(x)
	return -1*math.Cos(x*math.Pi/2) + 1
}
func (self *tweenings) EaseOutSine(x float64) float64 {
	self.checkRange(x)
	return math.Sin(x * math.Pi / 2)
}
func (self *tweenings) EaseInOutSine(x float64) float64 {
	self.checkRange(x)
	return x - 0.5*(math.Cos(math.Pi*x)-1)
}
func (self *tweenings) EaseInExpo(x float64) float64 {
	self.checkRange(x)
	if x == 0 {
		return 0.0
	}
	return math.Pow(2, (10 * (x - 1)))
}
func (self *tweenings) EaseOutExpo(x float64) float64 {
	self.checkRange(x)
	if x == 1 {

		return 1
	}

	return -math.Pow(2, (-10*x)) + 1
}
func (self *tweenings) EaseInOutExpo(x float64) float64 {
	self.checkRange(x)
	if x == 0 {
		return 0

	}
	if x == 1 {
		return 1
	}
	x = x * 2
	if x < 1 {
		return 0.5 * math.Pow(2, 10*(x-1))
	} else {
		x = x - 1
		return 0.5 * (-1*math.Pow(2, (-10*x)) + 2)
	}

}
func (self *tweenings) EaseInCirc(x float64) float64 {
	self.checkRange(x)
	return x - 1*(math.Sqrt(1-x*x)-1)
}
func (self *tweenings) EaseOutCirc(x float64) float64 {
	self.checkRange(x)
	x = x - 1
	return math.Sqrt(1 - (x * x))
}
func (self *tweenings) EaseInOutCirc(x float64) float64 {
	self.checkRange(x)
	x = x * 2
	if x < 1 {
		return -0.5 * (math.Sqrt(1-math.Pow(x, 2)) - 1)
	} else {
		x = x - 2
		return 0.5 * (math.Sqrt(1-math.Pow(x, 2)) + 1)
	}

}
func (self *tweenings) EaseInElastic(x float64) float64 {
	self.checkRange(x)
	amplitute := float64(1.0)
	period := float64(0.3)
	s := period / (2 * math.Pi) * math.Asin(1/amplitute)
	return -amplitute * math.Pow(2, (-10*(1-x))) * math.Sin(((1-x)-s)*(2*math.Pi/period))
}
func (self *tweenings) EaseOutElastic(x float64) float64 {
	self.checkRange(x)
	amplitute := float64(1.0)
	period := float64(0.3)
	s := period / (2 * math.Pi) * math.Asin(1/amplitute)
	return amplitute*math.Pow(2, (-10*x))*math.Sin((x-s)*(2*math.Pi/period)) + 1
}
func (self *tweenings) EaseInOutElastic(x float64) float64 {
	self.checkRange(x)
	amplitute := float64(1.0)
	period := float64(0.3)
	s := period / (2 * math.Pi) * math.Asin(1/amplitute)
	x = x * 2
	if x < 1 {
		return -amplitute * math.Pow(2, (-10*(1-x))) * math.Sin(((1-x)-s)*(2*math.Pi/period)) / 2

	} else {
		return amplitute*math.Pow(2, (-10*x))*math.Sin((x-s)*(2*math.Pi/period))/2 + 1
	}
}
func (self *tweenings) EaseInBack(x float64) float64 {
	self.checkRange(x)
	s := float64(1.70158)

	return x * x * ((s+1)*x - s)
}
func (self *tweenings) EaseOutBack(x float64) float64 {
	self.checkRange(x)
	s := float64(1.70158)
	x = x - 1
	return x*x*((s+1)*x+s) + 1
}
func (self *tweenings) EaseInOutBack(x float64) float64 {
	self.checkRange(x)
	s := float64(1.70158)
	x = x * 2
	if x < 1 {
		s *= 1.525
		return 0.5 * (x * x * ((s+1)*x - s))

	} else {
		x -= 2
		s *= 1.525
		return 0.5 * (x*x*((s+1)*x+s) + 2)

	}
}
func (self *tweenings) EaseInBounce(x float64) float64 {
	self.checkRange(x)
	return 1 - self.EaseOutBounce(1-x)
}
func (self *tweenings) EaseOutBounce(x float64) float64 {
	self.checkRange(x)
	if x < (1 / 2.75) {
		return 7.5625 * x * x
	} else {
		if x < (2 / 2.75) {
			x -= (1.5 / 2.75)
			return 7.5625*x*x + 0.75
		} else {
			if x < (2.5 / 2.75) {
				x -= (2.25 / 2.75)
				return 7.5625*x*x + 0.9375
			} else {
				x -= (2.65 / 2.75)
				return 7.5625*x*x + 0.984375
			}
		}
	}

}
func (self *tweenings) EaseInOutBounce(x float64) float64 {
	self.checkRange(x)
	if x < 0.5 {
		return self.EaseInBounce(x*2) * 0.5
	} else {
		return self.EaseOutBounce(x*2-1)*0.5 + 0.5
	}

}

func Tweening(x float64, method string) float64 {

	return (MethodMaps[method].Call([]reflect.Value{reflect.ValueOf(x)}))[0].Float()
}

// GetPointsOnLine 传入两点的坐标和曲线方法，返回轨迹点的列表
func GetPointsOnLine(x1, y1, x2, y2 int, method string) []robotgo.Point {
	points := make([]robotgo.Point, 0)
	var point robotgo.Point
	dx := x2 - x1
	dy := y2 - y1
	//保证除数不为零
	dxAbs := int(math.Max(1, math.Abs(float64(dx))))
	dyAbs := int(math.Max(1, math.Abs(float64(dy))))

	xStep := dx / dxAbs
	yStep := dy / dyAbs
	flag := dxAbs > dyAbs

	if flag {
		for i := 1; i <= dxAbs; i++ {

			point = robotgo.Point{X: x1 + xStep*i, Y: y1 + int(float64(yStep*dyAbs)*Tweening(float64(i)/float64(dxAbs), method))}
			points = append(points, point)
		}
	} else {
		for i := 1; i <= dyAbs; i++ {

			point = robotgo.Point{X: x1 + int(float64(xStep*dxAbs)*Tweening(float64(i)/float64(dyAbs), method)), Y: y1 + yStep*i}
			points = append(points, point)

		}
	}

	return (points)
}
