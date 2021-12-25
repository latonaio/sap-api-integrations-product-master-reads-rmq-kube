package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-product-master-reads-rmq-kube/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library/logger"

	"golang.org/x/xerrors"
)

type RMQOutputter interface {
	Send(sendQueue string, payload map[string]interface{}) error
}

type SAPAPICaller struct {
	baseURL      string
	apiKey       string
	outputQueues []string
	outputter    RMQOutputter
	log          *logger.Logger
}

func NewSAPAPICaller(baseUrl string, outputQueueTo []string, outputter RMQOutputter, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL:      baseUrl,
		apiKey:       GetApiKey(),
		outputQueues: outputQueueTo,
		outputter:    outputter,
		log:          l,
	}
}

func (c *SAPAPICaller) AsyncGetProductMaster(product, plant, mrpArea, valuationArea, productSalesOrg, productDistributionChnl string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "General":
			func() {
				c.General(product)
				wg.Done()
			}()
		case "Plant":
			func() {
				c.Plant(product, plant)
				wg.Done()
			}()
		case "MRPArea":
			func() {
				c.MRPArea(product, plant, mrpArea)
				wg.Done()
			}()
		case "Procurement":
			func() {
				c.Procurement(product, plant)
				wg.Done()
			}()
		case "WorkScheduling":
			func() {
				c.WorkScheduling(product, plant)
				wg.Done()
			}()
		case "SalesPlant":
			func() {
				c.SalesPlant(product, plant)
				wg.Done()
			}()
		case "Accounting":
			func() {
				c.Accounting(product, valuationArea)
				wg.Done()
			}()
		case "SalesOrganization":
			func() {
				c.SalesOrganization(product, productSalesOrg, productDistributionChnl)
				wg.Done()
			}()
		case "ProductDesc":
			func() {
				c.ProductDesc(product)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) General(product string) {
	data, err := c.callProductSrvAPIRequirementGeneral("A_Product", product)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "General"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementGeneral(api, product string) ([]sap_api_output_formatter.General, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithGeneral(req, product)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToGeneral(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) Plant(product, plant string) {
	data, err := c.callProductSrvAPIRequirementPlant("A_ProductPlant", product, plant)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "Plant"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementPlant(api, product, plant string) ([]sap_api_output_formatter.Plant, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithPlant(req, product, plant)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToPlant(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) MRPArea(product, plant, mrpArea string) {
	data, err := c.callProductSrvAPIRequirementMRPArea("A_ProductPlantMRPArea", product, plant, mrpArea)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "MRPArea"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementMRPArea(api, product, plant, mrpArea string) ([]sap_api_output_formatter.MRPArea, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithMRPArea(req, product, plant, mrpArea)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToMRPArea(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) Procurement(product, plant string) {
	data, err := c.callProductSrvAPIRequirementProcurement("A_ProductPlantProcurement", product, plant)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "Procurement"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementProcurement(api, product, plant string) ([]sap_api_output_formatter.Procurement, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithProcurement(req, product, plant)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToProcurement(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) WorkScheduling(product, plant string) {
	data, err := c.callProductSrvAPIRequirementWorkScheduling("A_ProductWorkScheduling", product, plant)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "WorkScheduling"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementWorkScheduling(api, product, plant string) ([]sap_api_output_formatter.WorkScheduling, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithWorkScheduling(req, product, plant)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToWorkScheduling(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) SalesPlant(product, plant string) {
	data, err := c.callProductSrvAPIRequirementSalesPlant("A_ProductPlantSales", product, plant)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "SalesPlant"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementSalesPlant(api, product, plant string) ([]sap_api_output_formatter.SalesPlant, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithSalesPlant(req, product, plant)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToSalesPlant(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) Accounting(product, valuationArea string) {
	data, err := c.callProductSrvAPIRequirementAccounting("A_ProductValuationAccount", product, valuationArea)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "Accounting"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementAccounting(api, product, valuationArea string) ([]sap_api_output_formatter.Accounting, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithAccounting(req, product, valuationArea)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToAccounting(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) SalesOrganization(product, productSalesOrg, productDistributionChnl string) {
	data, err := c.callProductSrvAPIRequirementSalesOrganization("A_ProductSalesDelivery", product, productSalesOrg, productDistributionChnl)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "SalesOrganization"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementSalesOrganization(api, product, productSalesOrg, productDistributionChnl string) ([]sap_api_output_formatter.SalesOrganization, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithSalesOrganization(req, product, productSalesOrg, productDistributionChnl)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToSalesOrganization(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) ProductDesc(product string) {
	data, err := c.callProductSrvAPIRequirementProductDesc(fmt.Sprintf("A_ProductDescription(Product='%s',Language='JA')", product))
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "ProductDesc"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProductSrvAPIRequirementProductDesc(api string) ([]sap_api_output_formatter.ProductDesc, error) {
	url := strings.Join([]string{c.baseURL, "API_PRODUCT_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToProductDesc(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithGeneral(req *http.Request, product string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s'", product))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithPlant(req *http.Request, product, plant string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s' and Plant eq '%s'", product, plant))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithMRPArea(req *http.Request, product, plant, mrpArea string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s' and Plant eq '%s' and MRPArea eq '%s'", product, plant, mrpArea))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithProcurement(req *http.Request, product, plant string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s' and Plant eq '%s'", product, plant))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithWorkScheduling(req *http.Request, product, plant string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s' and Plant eq '%s'", product, plant))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithSalesPlant(req *http.Request, product, plant string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s' and Plant eq '%s'", product, plant))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithAccounting(req *http.Request, product, valuationArea string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s' and ValuationArea eq '%s'", product, valuationArea))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithSalesOrganization(req *http.Request, product, productSalesOrg, productDistributionChnl string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s' and ProductSalesOrg eq '%s' and ProductDistributionChnl eq '%s'", product, productSalesOrg, productDistributionChnl))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithProductDesc(req *http.Request, product string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Product eq '%s'", product))
	req.URL.RawQuery = params.Encode()
}