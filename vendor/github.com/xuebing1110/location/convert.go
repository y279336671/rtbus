package location

import (
	"math"
)

const (
	PI  = 3.14159265358979324
	XPI = 3.14159265358979324 * 3000.0 / 180.0
)

//WGS-84 to GCJ-02
func GCJEncrypt(wgsLat, wgsLon float64) (float64, float64) {
	if OutOfChina(wgsLat, wgsLon) {
		return wgsLat, wgsLon
	}

	lat, lon := delta(wgsLat, wgsLon)
	return wgsLat + lat, wgsLon + lon
}

//GCJ-02 to WGS-84
func GCJDencrypt(gcjLat, gcjLon float64) (float64, float64) {
	if OutOfChina(gcjLat, gcjLon) {
		return gcjLat, gcjLon
	}

	lat, lon := delta(gcjLat, gcjLon)
	return gcjLat - lat, gcjLon - lon
}

//BD-09 to GCJ-02
func BdDencrypt(bdLat, bdLon float64) (float64, float64) {
	x := bdLon - 0.0065
	y := bdLat - 0.006
	var z = math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*XPI)
	var theta = math.Atan2(y, x) - 0.000003*math.Cos(x*XPI)
	var gcjLon = z * math.Cos(theta)
	var gcjLat = z * math.Sin(theta)
	return gcjLat, gcjLon
}

//GCJ-02 to BD-09
func BdEncrypt(gcjLat, gcjLon float64) (float64, float64) {
	x := gcjLon
	y := gcjLat
	var z = math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*XPI)
	var theta = math.Atan2(y, x) + 0.000003*math.Cos(x*XPI)
	bdLon := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006
	return bdLat, bdLon
}

//WGS-84 to Web mercator
//mercatorLat -> y mercatorLon -> x
func MercatorEncrypt(wgsLat, wgsLon float64) (float64, float64) {
	var x = wgsLon * 20037508.34 / 180.
	var y = math.Log(math.Tan((90.+wgsLat)*PI/360.)) / (PI / 180.)
	y = y * 20037508.34 / 180.
	return y, x
}

// Web mercator to WGS-84
// mercatorLat -> y mercatorLon -> x
func MercatorDecrypt(mercatorLat, mercatorLon float64) (float64, float64) {
	var x = mercatorLon / 20037508.34 * 180.
	var y = mercatorLat / 20037508.34 * 180.
	y = 180 / PI * (2*math.Atan(math.Exp(y*PI/180.)) - PI/2)
	return y, x
}

// two point's distance
func Distance(latA, lonA, latB, lonB float64) float64 {
	var earthR = 6371000.
	var x = math.Cos(latA*PI/180.) * math.Cos(latB*PI/180.) * math.Cos((lonA-lonB)*PI/180)
	var y = math.Sin(latA*PI/180.) * math.Sin(latB*PI/180.)
	var s = x + y
	if s > 1 {
		s = 1
	} else if s < -1 {
		s = -1
	}
	var alpha = math.Acos(s)
	var distance = alpha * earthR
	return distance
}

func OutOfChina(lat, lon float64) bool {
	if lon < 72.004 || lon > 137.8347 {
		return true
	} else if lat < 0.8293 || lat > 55.8271 {
		return true
	} else {
		return false
	}
}

func delta(lat, lon float64) (float64, float64) {
	a := 6378245.0               //  a: 卫星椭球坐标投影到平面地图坐标系的投影因子。
	ee := 0.00669342162296594323 //  ee: 椭球的偏心率。
	dLat := transformLat(lon-105.0, lat-35.0)
	dLon := transformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * PI
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * PI)
	dLon = (dLon * 180.0) / (a / sqrtMagic * math.Cos(radLat) * PI)
	return dLat, dLon
}

func transformLon(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*PI) + 20.0*math.Sin(2.0*x*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*PI) + 40.0*math.Sin(x/3.0*PI)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*PI) + 300.0*math.Sin(x/30.0*PI)) * 2.0 / 3.0
	return ret
}

func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*PI) + 20.0*math.Sin(2.0*x*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*PI) + 40.0*math.Sin(y/3.0*PI)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*PI) + 320*math.Sin(y*PI/30.0)) * 2.0 / 3.0
	return ret
}
