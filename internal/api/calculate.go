package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"math"
	"net/http"
	"time"
)

var token = "Hg45ArLPEqiQ7-weSrTGo"
var threshold = 0.25
const float64EqualityThreshold = 1e-9

type calculationParam struct {
	CalculationID   int `json:"calculation_id"`
}

type outputParam struct {
	CalculationID    int `json:"calculation_id"`
	OutputParam float64 `json:"output_param"`
	OutputErrorParam string `json:"output_error_param"`
}

type calcReq struct {
	ID                   int          `json:"id"`
	InputFirstParam        float64          `json:"input_first_param"`
	InputSecondParam float64          `json:"input_second_param"`
	Params               []calculationParam `json:"calculations"`
}

func Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("as_here0")

	req := &calcReq{}
	fmt.Println("as_here1")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	go worker(req.ID, req.InputFirstParam, req.InputSecondParam, req.Params)

}

func worker(id int, input_first_param float64, input_second_param float64, params []calculationParam) {
	time.Sleep(6 * time.Second)

	var results []outputParam

	for _, param := range params {
		if param.CalculationID != 2 && param.CalculationID != 3 && param.CalculationID != 4 {
			result_type_identifier := calculated_operation_probapility()
			if almostEqual(result_type_identifier, -1.0) {
				resultCalculating := "error"
				result := outputParam{
					CalculationID:    param.CalculationID,
					OutputParam:      result_type_identifier,
					OutputErrorParam: resultCalculating,
				}
				results = append(results, result)
				continue
			}
		}

		var resultCalculating float64

		if param.CalculationID == 2 {
			resultCalculating = float64(gcd(int(input_first_param), int(input_second_param)))
		} else if param.CalculationID == 3 {
			resultCalculating = float64(lcm(int(input_first_param), int(input_second_param)))
		} else if param.CalculationID == 4 {
			resultCalculating = math.Pow(input_first_param, 1/input_second_param)
		} else {
			result_type_identifier := calculated_operation_probapility()
			resultCalculating = result_type_identifier * input_first_param * input_second_param
		}

		// fmt.Println(param.CalculationID)
		// fmt.Println(resultCalculating)

		// Создаем новый элемент outputParam и добавляем в слайс results
		result := outputParam{
			CalculationID:    param.CalculationID,
			OutputParam:      resultCalculating,
			OutputErrorParam: "",
		}
		results = append(results, result)
	}
	// fmt.Printf("Simulated Metro Load for %s: %d people\n", param.ModelID, resultLoading)

	// Create the JSON payload
	requestData := map[string]interface{}{
		"application_id": id,
		"results":        results,
		"token":          token,
	}

	// Convert data to JSON
	postBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
		return
	}

	postURL := "http://localhost:8000/api/applications_calculations/write_result_calculating/" 


	// Make the HTTP POST request
	resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println("Ошибка при создании POST запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка: Неожиданный код статуса ответа %d\n", resp.StatusCode)
		return
	}
}

func calculated_operation_probapility() float64 {
	propability := rand.Float64()
	if propability < threshold{
		return -1.0
	}
	return propability
}

func almostEqual(a float64, b float64) bool {
    return math.Abs(a - b) <= float64EqualityThreshold
}

// Функция для нахождения наибольшего общего делителя (НОД)
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Функция для нахождения наименьшего общего кратного (НОК)
func lcm(a, b int) int {
	return a * b / gcd(a, b)
}
