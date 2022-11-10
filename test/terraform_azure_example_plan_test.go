package test

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/terraform"

	//test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	//"path/filepath"
	"testing"
	//"regexp"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	//"github.com/gruntwork-io/terratest/modules/k8s"
	//"github.com/hashicorp/go-version"
	//"io/ioutil"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
)

var planstruct terraform.PlanStruct
var jsonPlan string

// An example of how to test the Terraform module in examples/terraform-aws-example using Terratest.
func TestTerraformAzureExamplePlan(t *testing.T) {
	t.Parallel()

	const (
		substr = "runs/run-"
	)
	// Make a copy of the terraform module to a temporary directory. This allows running multiple tests in parallel
	// against the same terraform module.
	//exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-azure-example")

	// Give this EC2 Instance a unique ID for a name tag so we can distinguish it from any other EC2 Instance running
	// in your AWS account
	expectedSAName := terraform.GetVariableAsStringFromVarFile(t, "../examples/terraform-azure-example/terraform.tfvars", "storageaccountname")
	expectedSAHttpSettings, _ := strconv.ParseBool(terraform.GetVariableAsStringFromVarFile(t, "../examples/terraform-azure-example/terraform.tfvars", "enable_https_traffic_only"))

	//expectedName1 := "ajustorageAccount1234123"
	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	//awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	//instanceType := aws.GetRecommendedInstanceType(t, awsRegion, []string{"t2.micro", "t3.micro"})

	// website::tag::1::Configure Terraform setting path to Terraform code, EC2 instance name, and AWS Region. We also
	// configure the options with default retryable errors to handle the most common retryable errors encountered in
	// terraform testing.
	//planFilePath := filepath.Join(exampleFolder, "plan.out")
	//planFilePath := "plan.out"
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-azure-example",

		// Variables to pass to our Terraform code using -var options
		/*Vars: map[string]interface{}{
			"storageaccountname": expectedName,
			"enable_https_traffic_only":false,
			//"instance_type": instanceType,
		},*/

		// Environment variables to set when running Terraform
		/*EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},*/

		// Configure a plan file path so we can introspect the plan and make assertions about it.
		//PlanFilePath: planFilePath,
	})

	// website::tag::2::Run `terraform init`, `terraform plan`, and `terraform show` and fail the test if there are any errors
	//plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)
	planTfc := terraform.InitAndPlan(t, terraformOptions)
	fmt.Println("------------------------------------Plan follows---------------------------------------")
	fmt.Print(planTfc)
	fmt.Println("------------------------------------Plan Ends------------------------------------------")
	i := strings.Index(planTfc, substr)
	fmt.Println(i)
	fmt.Println(planTfc[i+5 : i+5+20])
	runId := planTfc[i+5 : i+5+20]
	url := "https://app.terraform.io/api/v2/runs/#runid#/plan/json-output"
	url = strings.Replace(url, "#runid#", runId, -1)
	fmt.Println(url)
	token := "O0kIzeAfiRjazA.atlasv1.crll4YfhBTe7VtUrFyLvYbgF9YtxNzzqE8RwdBQpw05Ut059o3gw5r7AHEdQhndqfj4"

	bearer := "Bearer " + token

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERRO] -", err)
	} else {
		defer resp.Body.Close()
		jsonPlan, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("==================================Response Body====================")
		fmt.Println(jsonPlan)
		fmt.Println("==================================Response Body====================")
		fmt.Println("==================================Plan From API====================")
		fmt.Println(string(jsonPlan))
		fmt.Println("==================================Plan From API Ends================")
		//json.Unmarshal([]byte(string(data)), &planstruct)
		err := os.WriteFile("cloudjson.json", []byte(string(jsonPlan)), 0666)
		if err != nil {

		}
	}

	// planStruct := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t,terraformOptions)

	//assert.Equal(t, 1, len(planStruct.ResourceChangesMap))
	/*fmt.Println("Point 1")
	plan2, err :=parsePlanJson(string(jsonPlan))
	fmt.Println(plan2)*/
	/*plan2 := &terraform.PlanStruct{}

	json.Unmarshal([]byte(string(jsonPlan)), &plan2.RawPlan)
	fmt.Println(plan2.RawPlan)*/
	/*plannedValues := plan2.RawPlan.PlannedValues
	fmt.Println(plannedValues)*/

	// website::tag::3::Use the go struct to introspect the plan values.

	//azuretags := azureResource.AttributeValues["tags"].(map[string]interface{})
	//assert.Equal(t, map[string]interface{}{"Name": expectedName}, azuretags)

	// website::tag::4::Alternatively, you can get the direct JSON output and use jsonpath to extract the data.
	// jsonpath only returns lists.
	//var jsonEC2Tags []map[string]interface{}
	//jsonOut := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := ioutil.ReadFile("cloudjson.json")
	//terraform.parsePlanJson(content)
	plan2 := &planstruct
	json.Unmarshal([]byte(content), &plan2.RawPlan)
	//planstruct.ResourcePlannedValuesMap = terraform.parsePlannedValues(planstruct)
	//planstruct.ResourceChangesMap = terraform.parseResourceChanges(planstruct)
	fmt.Println("+++++++++++++++++++++Raw Plan+++++++++++++++++++++++++++")
	fmt.Println(plan2.RawPlan)
	fmt.Println("+++++++++++++++++++++Raw Plan End+++++++++++++++++++++++")
	plan2.ResourcePlannedValuesMap = parsePlannedValues(plan2)
	plan2.ResourceChangesMap = parseResourceChanges(plan2)
	fmt.Println("+++++++++++++++++++++Resourse Planned++++++++++++++++++++++++++++++++")
	fmt.Println(plan2.ResourcePlannedValuesMap)
	fmt.Println("+++++++++++++++++++++Resourse Planned Ends+++++++++++++++++++++++++++")
	fmt.Println("+++++++++++++++++++++Resourse Changes++++++++++++++++++++++++++++++++")
	fmt.Println(plan2.ResourceChangesMap)
	fmt.Println("+++++++++++++++++++++Resourse Changes Ends+++++++++++++++++++++++++++")

	terraform.RequirePlannedValuesMapKeyExists(t, plan2, "azurerm_storage_account.aju-storageaccount")
	azureResource := planstruct.ResourcePlannedValuesMap["azurerm_storage_account.aju-storageaccount"]
	azurestoreagename := azureResource.AttributeValues["name"]
	enable_https_traffic_only := azureResource.AttributeValues["enable_https_traffic_only"]
	assert.Equal(t, expectedSAName, azurestoreagename)
	assert.Equal(t, expectedSAHttpSettings, enable_https_traffic_only)
	fmt.Println("Expected Name: " + expectedSAName)
	fmt.Println("Storage Account Name: ")
	fmt.Print(azurestoreagename)

	/*k8s.UnmarshalJSONPath(
		t,
		[]byte(content),
		"{ .planned_values.root_module.resources[0].values.tags }",
		&jsonEC2Tags,
	)
	assert.Equal(t, map[string]interface{}{"Name": expectedName}, jsonEC2Tags[0])*/

}
func parsePlannedValues(plan *terraform.PlanStruct) map[string]*tfjson.StateResource {
	plannedValues := plan.RawPlan.PlannedValues
	if plannedValues == nil {
		// No planned values, so return empty map.
		return map[string]*tfjson.StateResource{}
	}

	rootModule := plannedValues.RootModule
	if rootModule == nil {
		// No module resources, so return empty map.
		return map[string]*tfjson.StateResource{}
	}
	return parseModulePlannedValues(rootModule)
}
func parseModulePlannedValues(module *tfjson.StateModule) map[string]*tfjson.StateResource {
	out := map[string]*tfjson.StateResource{}
	for _, resource := range module.Resources {
		// NOTE: the Address attribute of the module resource always returns the full address, even when the resource is
		// nested within sub modules.
		out[resource.Address] = resource
	}

	// NOTE: base case of recursion is when ChildModules is empty list.
	for _, child := range module.ChildModules {
		// Recurse in to the child module. We take a recursive approach here despite limitations of the recursion stack
		// in golang due to the fact that it is rare to have heavily deep module calls in Terraform. So we optimize for
		// code readability as opposed to performance.
		childMap := parseModulePlannedValues(child)
		for k, v := range childMap {
			out[k] = v
		}
	}
	return out
}
func parseResourceChanges(plan *terraform.PlanStruct) map[string]*tfjson.ResourceChange {
	out := map[string]*tfjson.ResourceChange{}
	for _, change := range plan.RawPlan.ResourceChanges {
		out[change.Address] = change
	}
	return out
}
