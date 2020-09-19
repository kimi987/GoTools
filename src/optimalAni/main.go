package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	// "time"
)

type AnimationClip1 struct {
	M_objectHideFlags int    `yaml:"m_ObjectHideFlags"`
	M_name            string `yaml:"m_Name"`
}
type AnimationClip struct {
	ObjectHideFlags          int                 `yaml:"m_ObjectHideFlags"`
	PrefabParentObject       map[string]int      `yaml:"m_PrefabParentObject,flow"`
	PrefabInternal           map[string]int      `yaml:"m_PrefabInternal,flow"`
	Name                     string              `yaml:"m_Name"`
	SerializedVersion        int                 `yaml:"serializedVersion"`
	Legacy                   int                 `yaml:"m_Legacy"`
	Compressed               int                 `yaml:"m_Compressed"`
	UseHighQualityCurve      int                 `yaml:"m_UseHighQualityCurve"`
	RotationCurves           []*RotationCurves   `yaml:"m_RotationCurves"`
	CompressedRotationCurves []*RotationCurves   `yaml:"m_CompressedRotationCurves"`
	EulerCurves              []*RotationCurves   `yaml:"m_EulerCurves"`
	PositionCurves           []*RotationCurves   `yaml:"m_PositionCurves"`
	ScaleCurves              []*RotationCurves   `yaml:"m_ScaleCurves"`
	FloatCurves              []*FloatCurve       `yaml:"m_FloatCurves"`
	PPtrCurves               []*RotationCurves   `yaml:"m_PPtrCurves"`
	SampleRate               int                 `yaml:"m_SampleRate"`
	WrapMode                 int                 `yaml:"m_WrapMode"`
	Bounds                   Bounds              `yaml:"m_Bounds"`
	ClipBindingConstant      ClipBindingConstant `yaml:"m_ClipBindingConstant"`

	AnimationClipSettings AnimationClipSettings `yaml:"m_AnimationClipSettings"`

	EditorCurves      []*FloatCurve `yaml:"m_EditorCurves"`
	EulerEditorCurves []*FloatCurve `yaml:"m_EulerEditorCurves"`

	HasGenericRootTransform float32 `yaml:"m_HasGenericRootTransform"`
	HasMotionFloatCurves    float32 `yaml:"m_HasMotionFloatCurves"`
	GenerateMotionCurves    float32 `yaml:"m_GenerateMotionCurves"`
	Events                  []Event `yaml:"m_Events"`
}
type Event struct {
	Time                     float32        `yaml:"time"`
	FunctionName             string         `yaml:"functionName"`
	Data                     string         `yaml:"data,flow"`
	ObjectReferenceParameter map[string]int `yaml:"objectReferenceParameter,flow"`
	FloatParameter           float32        `yaml:"floatParameter"`
	IntParameter             int            `yaml:"intParameter"`
	MessageOptions           float32        `yaml:"messageOptions"`
}

type AnimationClipSettings struct {
	SerializedVersion         int            `yaml:"serializedVersion"`
	AdditiveReferencePoseClip map[string]int `yaml:"m_AdditiveReferencePoseClip,flow"`
	AdditiveReferencePoseTime int            `yaml:"m_AdditiveReferencePoseTime"`
	StartTime                 float32        `yaml:"m_StartTime"`
	StopTime                  float32        `yaml:"m_StopTime"`
	OrientationOffsetY        float32        `yaml:"m_OrientationOffsetY"`
	Level                     int            `yaml:"m_Level"`
	CycleOffset               float32        `yaml:"m_CycleOffset"`
	HasAdditiveReferencePose  float32        `yaml:"m_HasAdditiveReferencePose"`
	LoopTime                  float32        `yaml:"m_LoopTime"`
	LoopBlend                 float32        `yaml:"m_LoopBlend"`

	LoopBlendOrientation    float32 `yaml:"m_LoopBlendOrientation"`
	LoopBlendPositionY      float32 `yaml:"m_LoopBlendPositionY"`
	LoopBlendPositionXZ     float32 `yaml:"m_LoopBlendPositionXZ"`
	KeepOriginalOrientation float32 `yaml:"m_KeepOriginalOrientation"`

	KeepOriginalPositionY  float32 `yaml:"m_KeepOriginalPositionY"`
	KeepOriginalPositionXZ float32 `yaml:"m_KeepOriginalPositionXZ"`
	HeightFromFeet         float32 `yaml:"m_HeightFromFeet"`
	Mirror                 float32 `yaml:"m_Mirror"`
}

type ClipBindingConstant struct {
	GenericBindings  []MGenericBindings `yaml:"genericBindings"`
	PptrCurveMapping []int              `yaml:"pptrCurveMapping"`
}

type MGenericBindings struct {
	SerializedVersion int            `yaml:"serializedVersion"`
	Path              int            `yaml:"path"`
	Attribute         int            `yaml:"attributes"`
	Script            map[string]int `yaml:"script"`
	CustomType        int            `yaml:"customType"`
	IsPPtrCurve       int            `yaml:"isPPtrCurve"`
}

type Bounds struct {
	Center map[string]int `yaml:"m_Center,flow"`
	Extent map[string]int `yaml:"m_Extent,flow"`
}

type RotationCurves struct {
	Curve Curve  `yaml:"curve"`
	Path  string `yaml:"path"`
}

type FloatCurve struct {
	Curve     FCurve         `yaml:"curve"`
	attribute string         `yaml:"attribute"`
	Path      string         `yaml:"path"`
	ClassID   int            `yaml:"classID"`
	Script    map[string]int `yaml:"script,flow"`
}

type FCurve struct {
	SerializedVersion int        `yaml:"serializedVersion"`
	MCurve            []*FMCurve `yaml:"m_Curve"`
	PreInfinity       int        `yaml:"m_PreInfinity"`
	PostInfinity      int        `yaml:"m_PostInfinity"`
	RotationOrder     int        `yaml:"m_RotationOrder"`
}

type FMCurve struct {
	SerializedVersion int     `yaml:"serializedVersion"`
	Time              float32 `yaml:"time"`
	Value             float32 `yaml:"value"`
	InSlope           float32 `yaml:"inSlope"`
	OutSlope          float32 `yaml:"outSlope"`
	TangentMode       int     `yaml:"tangentMode"`

	LastMCurve *MCurve `yaml:"-"`
}

type Curve struct {
	SerializedVersion int       `yaml:"serializedVersion"`
	MCurve            []*MCurve `yaml:"m_Curve"`
	PreInfinity       int       `yaml:"m_PreInfinity"`
	PostInfinity      int       `yaml:"m_PostInfinity"`
	RotationOrder     int       `yaml:"m_RotationOrder"`
}

type MCurve struct {
	SerializedVersion int                `yaml:"serializedVersion"`
	Time              float64            `yaml:"time"`
	Value             map[string]float32 `yaml:"value,flow"`
	InSlope           map[string]float32 `yaml:"inSlope,flow"`
	OutSlope          map[string]float32 `yaml:"outSlope,flow"`
	TangentMode       int                `yaml:"tangentMode"`

	LastMCurve *MCurve `yaml:"-"`
}

type NewMCurve struct {
	SerializedVersion int     `yaml:"serializedVersion"`
	Time              float64 `yaml:"time"`
	Value             string  `yaml:"value"`
	InSlope           string  `yaml:"inSlope"`
	OutSlope          string  `yaml:"outSlope"`
	TangentMode       int     `yaml:"tangentMode"`
}

type conf struct {
	AnimationClip AnimationClip `yaml:"AnimationClip"`
}

const Head = `%YAML 1.1
%TAG !u! tag:unity3d.com,2011:
--- !u!74 &7400000`

type TimeRankMaps struct {
	TimeRankMap map[string][]*MCurve
	Data        map[string]*MCurve
}

var TimeRankMap *TimeRankMaps

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])

		if err != nil {
			fmt.Println("[Error] err = ", err)
			return
		}
		a := conf{}
		fmt.Println("start optimal path = ", os.Args[1])
		err = yaml.Unmarshal(data, &a)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		CompressData(&a.AnimationClip.RotationCurves, "RotationCurves")
		CompressData(&a.AnimationClip.CompressedRotationCurves, "CompressedRotationCurves")
		CompressData(&a.AnimationClip.EulerCurves, "EulerCurves")
		CompressData(&a.AnimationClip.PositionCurves, "PositionCurves")
		CompressData(&a.AnimationClip.ScaleCurves, "ScaleCurves")
		CompressFloatData(&a.AnimationClip.FloatCurves, "FloatCurves")
		CompressData(&a.AnimationClip.PPtrCurves, "PPtrCurves")

		data, err = yaml.Marshal(a)

		if err != nil {
			fmt.Println("err = ", err)
			return
		}

		dataString := string(data)

		dataString = strings.ReplaceAll(dataString, "\"y\"", "y")
		dataString = strings.ReplaceAll(dataString, "\"\"", "")
		dataString = fmt.Sprintf("%s\n%s", Head, dataString)
		ioutil.WriteFile(os.Args[1], []byte(dataString), os.ModePerm)

		return
	}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			if strings.Contains(info.Name(), ".anim") && !strings.Contains(info.Name(), ".meta") {
				//read anim files
				// if time.Now().Unix() - info.ModTime().Unix() > 18000 {
				// 	return nil
				// }

				data, _ := ioutil.ReadFile(path)
				a := conf{}
				fmt.Println("start optimal path = ", path)
				err := yaml.Unmarshal(data, &a)
				if err != nil {
					fmt.Println("err = ", err)
					return nil
				}
				// TimeRankMap =  &TimeRankMaps{
				// 	TimeRankMap:make(map[string][]*MCurve),
				// 	Data: make(map[string]*MCurve),
				// }

				// AddRankTime(&a.AnimationClip.RotationCurves)
				// AddRankTime(&a.AnimationClip.CompressedRotationCurves)
				// AddRankTime(&a.AnimationClip.EulerCurves)
				// AddRankTime(&a.AnimationClip.PositionCurves)
				// AddRankTime(&a.AnimationClip.ScaleCurves)
				// AddRankTime(&a.AnimationClip.FloatCurves)
				// AddRankTime(&a.AnimationClip.PPtrCurves)

				CompressData(&a.AnimationClip.RotationCurves, "RotationCurves")
				CompressData(&a.AnimationClip.CompressedRotationCurves, "CompressedRotationCurves")
				CompressData(&a.AnimationClip.EulerCurves, "EulerCurves")
				CompressData(&a.AnimationClip.PositionCurves, "PositionCurves")
				CompressData(&a.AnimationClip.ScaleCurves, "ScaleCurves")
				CompressFloatData(&a.AnimationClip.FloatCurves, "FloatCurves")
				CompressData(&a.AnimationClip.PPtrCurves, "PPtrCurves")

				data, err = yaml.Marshal(a)

				if err != nil {
					fmt.Println("err = ", err)
					return nil
				}

				dataString := string(data)

				dataString = strings.ReplaceAll(dataString, "\"y\"", "y")
				dataString = strings.ReplaceAll(dataString, "\"\"", "")
				dataString = fmt.Sprintf("%s\n%s", Head, dataString)
				ioutil.WriteFile(path, []byte(dataString), os.ModePerm)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("file walk err = ", err)
	}
}

func AddRankTime(data *[]*RotationCurves) {
	for _, v := range *data {
		ranks := TimeRankMap.TimeRankMap[v.Path]
		fmt.Println("len(ranks) 1 = ", len(ranks))
		for _, v1 := range v.Curve.MCurve {
			ranks = AddCurve(&ranks, v1)
			TimeRankMap.Data[fmt.Sprintf("%s_%f", v.Path, v1.Time)] = v1
		}
		fmt.Println("len(ranks) 2= ", len(ranks))
		TimeRankMap.TimeRankMap[v.Path] = ranks
	}
}

func AddCurve(ranks *[]*MCurve, curve *MCurve) []*MCurve {
	if len(*ranks) == 0 {
		*ranks = append(*ranks, curve)
		return *ranks
	}
	isAdd := false
	for k, v := range *ranks {
		if v.Time > curve.Time {
			if k == 0 {
				*ranks = append([]*MCurve{curve}, (*ranks)[k:]...)
			} else {
				*ranks = append((*ranks)[:k], append([]*MCurve{curve}, (*ranks)[k:]...)...)
			}

			curve.LastMCurve = v.LastMCurve
			v.LastMCurve = curve
			isAdd = true
			break
		}
	}

	if !isAdd {
		curve.LastMCurve = (*ranks)[len(*ranks)-1]
		*ranks = append(*ranks, curve)
	}

	return *ranks
}

func CompressData(data *[]*RotationCurves, name string) {

	for k1, v := range *data {
		var newCures []*MCurve
		// var lastValue map[string]float32
		// var lastInSlope map[string]float32
		// var lastOutSlope map[string]float32
		for k, v1 := range v.Curve.MCurve {
			if k != 0 && k != len(v.Curve.MCurve)-1 {
				if CheckMapValueIsSame(v.Curve.MCurve[k-1].Value, v1.Value) && CheckMapValueIsSame(v.Curve.MCurve[k-1].InSlope, v1.InSlope) &&
					CheckMapValueIsSame(v.Curve.MCurve[k-1].OutSlope, v1.OutSlope) {
					if CheckMapValueIsSame(v.Curve.MCurve[k+1].Value, v1.Value) && CheckMapValueIsSame(v.Curve.MCurve[k+1].InSlope, v1.InSlope) &&
						CheckMapValueIsSame(v.Curve.MCurve[k+1].OutSlope, v1.OutSlope) {
						continue
					}
				}
			}
			// if k != len(v.Curve.MCurve)-1 {
			// 	if CheckMapValueIsSame(lastValue, v1.Value) && CheckMapValueIsSame(lastInSlope, v1.InSlope) &&
			// 		CheckMapValueIsSame(lastOutSlope, v1.OutSlope)  {
			// 		// v_check := TimeRankMap.Data[fmt.Sprintf("%s_%f", v.Path, v1.Time)]
			// 		// if v_check != nil {
			// 		// 	fmt.Println("v_check != nil")
			// 		// 	if v_check.LastMCurve != nil {
			// 		// 		fmt.Println("lastMCurve != nil")
			// 		// 	}
			// 		// }
			// 		if v1 != nil && v1.LastMCurve != nil {
			// 			fmt.Println("lastMCurve != nil")
			// 			if CheckMapValueIsSame(v1.LastMCurve.Value, v1.Value)  &&
			// 				CheckMapValueIsSame(v1.LastMCurve.InSlope, v1.InSlope) && CheckMapValueIsSame(v1.LastMCurve.OutSlope, v1.OutSlope) {
			// 					continue;
			// 				}
			// 		}
			// 		// continue
			// 	}
			// }

			// lastValue = v1.Value
			// lastInSlope = v1.InSlope
			// lastOutSlope = v1.OutSlope
			newCures = append(newCures, v1)
		}
		(*data)[k1].Curve.MCurve = newCures
	}
}

func CompressFloatData(data *[]*FloatCurve, name string) {

	for k1, v := range *data {
		var newCures []*FMCurve
		// var lastValue map[string]float32
		// var lastInSlope map[string]float32
		// var lastOutSlope map[string]float32
		for k, v1 := range v.Curve.MCurve {
			if k != 0 && k != len(v.Curve.MCurve)-1 {
				if v.Curve.MCurve[k-1].Value == v1.Value && v.Curve.MCurve[k-1].InSlope == v1.InSlope &&
					v.Curve.MCurve[k-1].OutSlope == v1.OutSlope {
					if v.Curve.MCurve[k+1].Value == v1.Value && v.Curve.MCurve[k+1].InSlope == v1.InSlope &&
						v.Curve.MCurve[k+1].OutSlope == v1.OutSlope {
						continue
					}
				}
			}

			newCures = append(newCures, v1)
		}
		(*data)[k1].Curve.MCurve = newCures
	}
}

func CheckMapValueIsSame(val1, val2 map[string]float32) bool {
	if val1 == nil {
		return false
	}
	for k := range val1 {
		if val1[k] != val2[k] {
			return false
		}
	}
	return true
}
