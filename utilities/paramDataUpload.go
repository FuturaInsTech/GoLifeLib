package utilities

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/types"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func UploadParamDataItems(file multipart.File) error {

	f, err := excelize.OpenReader(file)

	if err != nil {
		return err
	}

	firstSheet := f.WorkBook.Sheets.Sheet[0].Name
	//fmt.Printf("'%s' is first sheet of %d sheets.\n", firstSheet, f.SheetCount)

	headerFieldMap := make(map[string]string)
	for i := 0; i < 5; i++ {
		k := i / 3
		j := i % 3

		cellName1, _ := excelize.CoordinatesToCellName((j*2)+1, k+3)
		cellName2, _ := excelize.CoordinatesToCellName((j*2)+2, k+3)
		cellVal1, _ := f.GetCellValue(firstSheet, cellName1)
		cellVal2, _ := f.GetCellValue(firstSheet, cellName2)
		headerFieldMap[cellVal1] = cellVal2

	}

	var headerParam models.Param

	result := initializers.DB.Model(&models.Param{}).Find(&headerParam, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =? and is_valid =?", headerFieldMap["Company"], headerFieldMap["Param Name"], "HE", "", 0, true)

	if result.Error != nil {

		return errors.New("unable to read param header:" + result.Error.Error())

	}

	if result.RowsAffected == 0 {

		return errors.New("param header not found")

	}

	paramType := "0"
	extraDataExist, _ := headerParam.Data["extraDataExist"].(bool)

	if headerParam.Data["paramType"] == "dated" {
		paramType = "D"
	} else {

		if extraDataExist {
			paramType = "1"
		}
	}

	if paramType == "0" {
		return errors.New("invalid param type - does not have extradata")
	}
	totalHeaderElem := 5
	seqNo := 0
	if paramType == "D" {
		totalHeaderElem = 8
		for i := 5; i < 8; i++ {
			k := i / 3
			j := i % 3

			cellName1, _ := excelize.CoordinatesToCellName((j*2)+1, k+3)
			cellName2, _ := excelize.CoordinatesToCellName((j*2)+2, k+3)
			cellVal1, _ := f.GetCellValue(firstSheet, cellName1)
			cellVal2, _ := f.GetCellValue(firstSheet, cellName2)
			headerFieldMap[cellVal1] = cellVal2

		}
		seqNo, err = strconv.Atoi(headerFieldMap["Seq No"])
		if err != nil {
			return errors.New("sequence number not valid:" + err.Error() + " param:" + headerFieldMap["Param Name"])
		}

		if seqNo > 0 {
			var params_existing []models.Param
			result := initializers.DB.Order("seqno asc").Find(&params_existing, "company_id  = ? AND  name = ?  AND item = ? AND  rec_type = ?", headerFieldMap["Company"], headerFieldMap["Param Name"], headerFieldMap["Item Name"], "IT")

			if result.Error != nil {

				return errors.New("Failed to get dated items:" + result.Error.Error() + " param:" + headerFieldMap["Param Name"])

			}

			if result.RowsAffected != int64(seqNo) {

				return errors.New("sequence number not valid for param:" + headerFieldMap["Param Name"])

			}

		}
	}

	//check if param item already exists
	var param models.Param
	result = initializers.DB.First(&param, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", headerFieldMap["Company"], headerFieldMap["Param Name"], "IT", headerFieldMap["Item Name"], seqNo)
	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if !recordNotFound && result.Error != nil {
		return errors.New(result.Error.Error() + "param :" + headerFieldMap["Param Name"])
	}

	//extra data processing begins

	fields, subfields, err := GetParamExtraDataFields(headerFieldMap["Param Name"])

	if err != nil {
		return errors.New("error occuered while extracting fields of param:" + headerFieldMap["Param Name"] + err.Error())

	}

	var extradata map[string]interface{}
	//check if the field is array type
	if subfields != nil {

		k := totalHeaderElem / 3
		j := totalHeaderElem % 3

		if j > 0 {
			k++
		}

		keyArray := make([]string, 0)

		for i := 0; ; i++ {
			cellName, _ := excelize.CoordinatesToCellName(i+1, k+4)
			cellVal, _ := f.GetCellValue(firstSheet, cellName)

			if cellVal == "" {
				break
			}
			keyArray = append(keyArray, cellVal)

		}
		subFieldMap := make(map[string]reflect.StructField)
		for _, v := range subfields {

			subFieldMap[v.Name] = v
		}
		valArray := make([]interface{}, 0)
		for i := 0; ; i++ {

			breakloop := true
			valMap := make(map[string]interface{})
			for j := 0; j < len(keyArray); j++ {

				cellName, _ := excelize.CoordinatesToCellName(j+1, i+k+5)
				cellVal, _ := f.GetCellValue(firstSheet, cellName)

				if subFieldMap[keyArray[j]].Type.Kind() == reflect.String {
					valMap[keyArray[j]] = cellVal

				} else {

					value, err := GetFormattedField(cellVal, subFieldMap[keyArray[j]])
					if err != nil {
						return errors.New(err.Error() + " field:" + keyArray[j])
					} else {
						//make first character lower case in the key

						valMap[strings.ToLower(string(keyArray[j][0]))+keyArray[j][1:]] = value

					}

				}

				if cellVal != "" {
					breakloop = false
				}

			}

			if breakloop {
				break
			}

			valArray = append(valArray, valMap)

		}

		extradata = map[string]interface{}{
			//make first character lower case in the key

			strings.ToLower(string(fields[0].Name[0])) + fields[0].Name[1:]: valArray,
		}

	} else {

		extraDataFieldMap := make(map[string]string)
		for i := totalHeaderElem; ; i++ {
			k := i / 3
			j := i % 3

			cellName1, _ := excelize.CoordinatesToCellName((j*2)+1, k+3)
			cellName2, _ := excelize.CoordinatesToCellName((j*2)+2, k+3)
			cellVal1, _ := f.GetCellValue(firstSheet, cellName1)

			if cellVal1 == "" {
				break
			}

			cellVal2, _ := f.GetCellValue(firstSheet, cellName2)
			extraDataFieldMap[cellVal1] = cellVal2

		}

		data1, err := GetFormattedExtraData(fields, extraDataFieldMap)

		if err != nil {
			return errors.New("error occured while formatting data param:" + headerFieldMap["Param Name"] + err.Error())
		}

		extradata = data1

	}

	if recordNotFound {
		param.Data = extradata
		param.CreatedAt = time.Now()
		val, _ := strconv.Atoi(headerFieldMap["Company"])

		param.CompanyId = uint16(val)
		param.Name = headerFieldMap["Param Name"]
		param.Item = headerFieldMap["Item Name"]
		if paramType != "D" {
			param.EndDate = "0"
			param.StartDate = "0"
		} else {

			date, err := time.Parse("02/01/2006", headerFieldMap["Start Date"])
			if err == nil {
				param.StartDate = date.Format("20060102")
			} else {
				return errors.New("incorrect start date :" + headerFieldMap["Start Date"] + " " + result.Error.Error())
			}

			date, err = time.Parse("02/01/2006", headerFieldMap["End Date"])
			if err == nil {
				param.EndDate = date.Format("20060102")
			} else {
				return errors.New("incorrect end date :" + headerFieldMap["End Date"] + " " + result.Error.Error())
			}

		}
		param.Is_valid = true
		param.RecType = "IT"
		param.LastModUser = 1
		param.Seqno = uint16(seqNo)

		result := initializers.DB.Create(&param)

		if result.Error != nil {

			return errors.New("error while creating  param item :" + headerFieldMap["Param Name"] + "-" + headerFieldMap["Item Name"] + "error:" + result.Error.Error())

		}

	} else {

		param.Data = extradata
		param.UpdatedAt = time.Now()
		param.LastModUser = 1
		if paramType == "D" {

			date, err := time.Parse("02/01/2006", headerFieldMap["Start Date"])
			if err == nil {
				param.StartDate = date.Format("20060102")
			} else {
				return errors.New("incorrect start date :" + headerFieldMap["Start Date"] + " " + err.Error())
			}

			date, err = time.Parse("02/01/2006", headerFieldMap["End Date"])
			if err == nil {
				param.EndDate = date.Format("20060102")
			} else {
				return errors.New("incorrect end date :" + headerFieldMap["End Date"] + " " + err.Error())
			}

		}

		result = initializers.DB.Model(&param).Where("company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and seqno =?", headerFieldMap["Company"], headerFieldMap["Param Name"], "IT", headerFieldMap["Item Name"], seqNo).Updates(param)

		if result.Error != nil {
			return errors.New("failed to update param item :" + headerFieldMap["Param Name"] + "-" + headerFieldMap["Item Name"] + "error:" + result.Error.Error())

		}

	}

	if seqNo == 0 {

		if recordNotFound {

			var paramdesc models.ParamDesc
			paramdesc.CreatedAt = time.Now()
			val, _ := strconv.Atoi(headerFieldMap["Company"])
			paramdesc.CompanyId = uint16(val)
			paramdesc.Name = headerFieldMap["Param Name"]
			paramdesc.Item = headerFieldMap["Item Name"]
			paramdesc.LanguageId = 1
			paramdesc.RecType = "IT"
			paramdesc.Longdesc = headerFieldMap["Long Description"]
			paramdesc.Shortdesc = headerFieldMap["Short Description"]
			paramdesc.LastModUser = 1

			result = initializers.DB.Create(&paramdesc)

			if result.Error != nil {

				return errors.New("Failed to insert param desc:" + headerFieldMap["Param Name"] + "-" + headerFieldMap["Item Name"] + "error:" + result.Error.Error())

			}

		} else {

			var paramdesc models.ParamDesc

			result = initializers.DB.First(&paramdesc, "company_id  = ? AND  name = ?  AND  rec_type = ?  AND item = ?  and language_id =?", headerFieldMap["Company"], headerFieldMap["Param Name"], "IT", headerFieldMap["Item Name"], 1)

			if result.Error != nil {

				return errors.New("failed to get param desc :" + headerFieldMap["Param Name"] + "-" + headerFieldMap["Item Name"] + "error:" + result.Error.Error())

			}

			paramdesc.UpdatedAt = time.Now()
			paramdesc.LastModUser = 1
			paramdesc.Longdesc = headerFieldMap["Long Description"]
			paramdesc.Shortdesc = headerFieldMap["Short Description"]

			result = initializers.DB.Model(&paramdesc).Updates(paramdesc)

			if result.Error != nil {

				return errors.New("Failed to update param desc :" + headerFieldMap["Param Name"] + "-" + headerFieldMap["Item Name"] + "error:" + result.Error.Error())

			}

		}

	}

	return nil

}

func GetParamExtraDataFields(paramName string) ([]reflect.StructField, []reflect.StructField, error) {

	var data any
	var subfield any

	switch paramName {

	case "Q0011":
		data = types.Q0011Data{}
	case "Q0017":
		data = types.Q0017Data{}
	case "Q0016":
		data = types.Q0016Data{}
	case "Q0015":
		data = types.Q0015Data{}

	case "Q0014":
		data = types.Q0014Data{}

	case "Q0005":
		data = types.Q0005Data{}

	case "Q0006":
		data = types.Q0006Data{}

	case "P0044":
		data = types.P0044Data{}

	case "Q0010":
		data = types.Q0010Data{}
		subfield = types.Q0010{}

	default:
		return nil, nil, fmt.Errorf("unable to find extradata struct for the param:" + paramName)

	}

	typeOfT := reflect.TypeOf(data)
	if typeOfT.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("can't reflect the fields of non-struct type %s", typeOfT.Elem().Name())
	}

	fields := reflect.VisibleFields(reflect.TypeOf(data))

	var subfields []reflect.StructField = nil
	//If array type then get the subfield list
	if len(fields) == 1 && (fields[0].Type.Kind() == reflect.Slice || fields[0].Type.Kind() == reflect.Array) {

		subfields = reflect.VisibleFields(reflect.TypeOf(subfield))

	}

	return fields, subfields, nil

}

func GetFormattedExtraData(fields []reflect.StructField, dataMap map[string]string) (map[string]interface{}, error) {

	dataMap1 := make(map[string]interface{})

	for _, f := range fields {

		if f.Type.Kind() == reflect.String {
			dataMap1[strings.ToLower(string(f.Name[0]))+f.Name[1:]] = dataMap[f.Name]

		} else {
			value, err := GetFormattedField(dataMap[f.Name], f)

			if err != nil {
				return nil, errors.New(err.Error() + " field:" + f.Name)
			}
			//make first character lower case in the key
			dataMap1[strings.ToLower(string(f.Name[0]))+f.Name[1:]] = value
		}

	}

	return dataMap1, nil

}

func GetFormattedField(value string, field reflect.StructField) (interface{}, error) {

	switch field.Type.Kind() {

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return boolValue, nil

	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		intVar := 0
		if value != "" {
			val, err := strconv.Atoi(value)

			if err != nil {
				return nil, err
			}

			intVar = val
		}
		return intVar, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		var intVar uint64 = 0
		if value != "" {
			val, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, err
			}
			intVar = val
		}
		return intVar, nil
	case reflect.Float32, reflect.Float64:
		var floatVal float64 = 0
		if value != "" {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			floatVal = val
		}
		return floatVal, nil

	case reflect.Slice, reflect.Array:
		s := strings.Split(value, ",")
		switch field.Type.Elem().Kind() {

		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:

			intArry := make([]int, 0)
			for _, v := range s {

				val, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				intArry = append(intArry, val)
			}
			return intArry, nil
		case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:

			intArry := make([]uint64, 0)
			for _, v := range s {

				val, err := strconv.ParseUint(v, 10, 32)
				if err != nil {
					return nil, err
				}
				intArry = append(intArry, val)
			}
			return intArry, nil

		case reflect.Float32, reflect.Float64:
			floatArry := make([]float64, 0)
			for _, v := range s {

				val, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, err
				}
				floatArry = append(floatArry, val)
			}
			return floatArry, nil

		case reflect.String:
			return s, nil

		default:
			fmt.Println("unhandled Type array field:" + field.Name + " value:" + value)
			return s, nil
		}
	default:
		fmt.Println("unhandled Type field:" + field.Name + " value:" + value)
		return value, nil

	}
}
