package api

import (
	"testing"
)

func TestParseValidationRedirectUriError(t *testing.T) {
	inputJSON := `
	{
		"code": "BadRequest",
		"message": "Validation Error",
		"details": [
			{
				"location": "body",
				"msg": "redirectUris is not array of strings",
				"param": "redirectUris"
			}
		]
	}`

	apiErr, err := ParseError([]byte(inputJSON))
	if err != nil {
		t.Fatalf("Failed to parse error: %v", err)
	}

	expectedStatusCode := 0 // Since "statusCode" is not provided in the input JSON
	expectedMethod := ""
	expectedURL := ""
	expectedCode := "BadRequest"
	expectedMessage := "Validation Error"
	expectedDetails := []ErrorDetail{
		{
			Location: "body",
			Msg:      "redirectUris is not array of strings",
			Param:    "redirectUris",
		},
	}

	// Check each field of ApiError
	if apiErr.StatusCode != expectedStatusCode {
		t.Fatalf("Expected StatusCode: %d, but got: %d", expectedStatusCode, apiErr.StatusCode)
	}
	if apiErr.Method != expectedMethod {
		t.Fatalf("Expected Method: %s, but got: %s", expectedMethod, apiErr.Method)
	}
	if apiErr.Url != expectedURL {
		t.Fatalf("Expected URL: %s, but got: %s", expectedURL, apiErr.Url)
	}
	if apiErr.Code != expectedCode {
		t.Fatalf("Expected Code: %s, but got: %s", expectedCode, apiErr.Code)
	}
	if apiErr.Message != expectedMessage {
		t.Fatalf("Expected Message: %s, but got: %s", expectedMessage, apiErr.Message)
	}

	// Check the details field
	if len(apiErr.Details) != len(expectedDetails) {
		t.Fatalf("Expected %d error details, but got %d", len(expectedDetails), len(apiErr.Details))
	}

	for i, detail := range apiErr.Details {
		if detail.Location != expectedDetails[i].Location {
			t.Fatalf("Expected Detail Location: %s, but got: %s", expectedDetails[i].Location, detail.Location)
		}
		if detail.Msg != expectedDetails[i].Msg {
			t.Fatalf("Expected Detail Msg: %s, but got: %s", expectedDetails[i].Msg, detail.Msg)
		}
		if detail.Param != expectedDetails[i].Param {
			t.Fatalf("Expected Detail Param: %s, but got: %s", expectedDetails[i].Param, detail.Param)
		}
	}

	expectedError := `Got 'BadRequest' error: Validation Error
Error in body with 'redirectUris': redirectUris is not array of strings`

	if apiErr.Error() != expectedError {
		t.Fatalf("Unexpected error format. Expected: \n%s\n\nActual:\n%s", apiErr.Error(), expectedError)
	}
}

func TestParseValidationProofTypeError(t *testing.T) {
	proofTypeError := `{"code":"BadRequest","message":"Validation Error","details":[{"value":["BbsBlsSignature2020"],"msg":"must be a valid credential proof type","param":"proofType","location":"body"}]}`

	apiErr, err := ParseError([]byte(proofTypeError))
	if err != nil {
		t.Fatalf("Failed to parse error: %v", err)
	}

	expectedError := "Got 'BadRequest' error: Validation Error\nError in body with 'proofType': must be a valid credential proof type"

	if apiErr.Error() != expectedError {
		t.Fatalf("Unexpected error format. Expected: \n%s\nActual:\n%s", apiErr.Error(), expectedError)
	}
}
