package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"neoway-challenge/dao"
	"neoway-challenge/dao/models"

	"github.com/Nhanderu/brdoc"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
)

type ShoppingData = models.ShoppingData

func ConvertAndSaveData(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		errMsg := fmt.Sprintf("Error reading request file: %s", err)
		log.Fatal(errMsg)
	}

	bulkData, err := parseDataFileToStruct(file)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing data file: %s", err)
		ctx.AbortWithStatusJSON(400, gin.H{"message": errMsg})
		return
	}

	if len(bulkData) == 0 {
		errMsg := fmt.Sprintf("Something went wrong while parsing and validating values")
		ctx.AbortWithStatusJSON(400, gin.H{"message": errMsg})
		return
	}

	err = dao.Save(bulkData)
	if err != nil {
		errMsg := fmt.Sprintf("Error saving data to db: %s", err)
		ctx.AbortWithStatusJSON(500, gin.H{"message": errMsg})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Data saved succesfully into the database.",
	})
}

func parseDataFileToStruct(requestDataFile *multipart.FileHeader) ([]interface{}, error) {
	// read the file extension to later decide how to parse the file lines (CSV vs TXT)
	filename := filepath.Base(requestDataFile.Filename)
	fileExtension := filepath.Ext(filename)

	// converts the *multipart.FileHeader file from gin into a os.File that can be used by scanner
	file, _ := requestDataFile.Open()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// skip the first line
	scanner.Scan()

	// create a map that will carry the structs to be saved in the DB
	var bulkData []interface{}
	for scanner.Scan() {
		line := scanner.Text()

		// handle data splitting for both CSV and TXT files
		var dataMap []string
		if fileExtension == ".csv" {
			log.Println("Found csv file!")
			dataMap = strings.Split(line, ",")
		} else if fileExtension == ".txt" {
			dataMap = convertTxtToMap(line)
		} else {
			return nil, errors.New("This API only supports TXT and CSV files.")
		}

		// the line will only be added if the CPF and CNPJs contained in it are valid
		if brdoc.IsCPF(dataMap[0]) && brdoc.IsCNPJ(dataMap[6]) && brdoc.IsCNPJ(dataMap[7]) {
			sanitizedData, err := sanitizeData(dataMap)

			var data = ShoppingData{
				CPF:                sanitizedData.CPF,
				Private:            sanitizedData.Private,
				Incompleto:         sanitizedData.Incompleto,
				DataUltimaCompra:   sanitizedData.DataUltimaCompra,
				TicketMedio:        sanitizedData.TicketMedio,
				TicketUltimaCompra: sanitizedData.TicketUltimaCompra,
				LojaMaisFrequente:  sanitizedData.LojaMaisFrequente,
				LojaUltimaCompra:   sanitizedData.LojaUltimaCompra,
			}

			// make a per field validation on the final struct before appending to the db payload
			if err = validator.Validate(data); err != nil {
				errMsg := fmt.Sprintf("Error validating data: %s", err)
				log.Println(errMsg)
			}

			bulkData = append(bulkData, data)
		}
	}

	return bulkData, nil

}

func convertTxtToMap(rawLine string) []string {
	return strings.Fields(rawLine)
}
func convertCsvToMap(rawLine string) []string {
	dataMap := strings.Fields(rawLine)
	return dataMap
}

func sanitizeData(data []string) (*ShoppingData, error) {
	var sanitizedData ShoppingData

	// sanitize the CPF by removing the separators
	cpf := data[0]
	if brdoc.IsCPF(cpf) {
		cpf = formatterOnlyNumbers(cpf)
	} else {
		return nil, errors.New("The document provided isnt a valid CPF")
	}
	sanitizedData.CPF = cpf

	var err error
	sanitizedData.Private, err = strconv.Atoi(data[1])
	if err != nil {
		return nil, errors.New("Coundt parse Private field for this line")
	}
	sanitizedData.Incompleto, err = strconv.Atoi(data[2])
	if err != nil {
		return nil, errors.New("Coundt parse Incompleto field for this line")
	}

	// converts null dates int NULL type the date into a standard go time.Time format
	stringDate := data[3]
	var date *time.Time
	if stringDate != "NULL" {
		format := "2006-01-02"
		parsedDate, _ := time.Parse(format, stringDate)
		date = &parsedDate
	} else {
		date = nil

	}
	sanitizedData.DataUltimaCompra = date

	// parse currency values to INT type, which on read is parsed considering the cents, the last 2 decimal points
	ticketMedio := formatterOnlyNumbers(data[4])
	sanitizedData.TicketMedio, err = strconv.Atoi(ticketMedio)
	if err != nil {
		return nil, errors.New("The document provided isnt a valid CNPJ")
	}

	ticketUltimaCompra := formatterOnlyNumbers(data[5])
	sanitizedData.TicketUltimaCompra, err = strconv.Atoi(ticketUltimaCompra)
	if err != nil {
		return nil, errors.New("The document provided isnt a valid CNPJ")
	}

	// sanitize CNPJs by removing separators
	var cnpj *string
	lojaMaisFrequente := data[6]
	if lojaMaisFrequente != "NULL" {
		lojaMaisFrequente = formatterOnlyNumbers(lojaMaisFrequente)
		cnpj = &lojaMaisFrequente

	} else {
		return nil, errors.New("The document provided isnt a valid CNPJ")
	}
	sanitizedData.LojaMaisFrequente = cnpj

	lojaUltimaCompra := data[7]
	if lojaUltimaCompra != "NULL" {
		lojaUltimaCompra = formatterOnlyNumbers(lojaMaisFrequente)
		cnpj = &lojaMaisFrequente
	} else {
		return nil, errors.New("The document provided isnt a valid CNPJ")
	}
	sanitizedData.LojaUltimaCompra = cnpj

	return &sanitizedData, nil

}

func formatterOnlyNumbers(str string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(str, "")
}
