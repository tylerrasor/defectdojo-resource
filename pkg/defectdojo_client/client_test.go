package defectdojo_client_test

// this test doesn't matter anymore, but I'll need the pattern of how to mock out the api
// func TestNewClientFailsWithBunkCreds(t *testing.T) {
// 	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			w.WriteHeader(http.StatusBadRequest)
// 		}
// 		if r.Header.Get("Authorization") == "Token bad_key" {
// 			w.WriteHeader(http.StatusUnauthorized)
// 		}
// 	}))
// 	defer mock_server.Close()

// 	client, err := defectdojo_client.NewDefectdojoClient(mock_server.URL, "username", "", "bad_key")

// 	assert.NotNil(t, err)
// 	assert.Errorf(t, err, "error in authenticating to defectdojo instance at `%s`: ", mock_server.URL)
// 	assert.Nil(t, client)
// }

// func TestTryToAuthSucceeds(t *testing.T) {
// 	source := models.Source{
// 		DefectDojoUrl: "http://something",
// 		Username:      "admin",
// 		ApiKey:        "good_key",
// 	}

// 	err := client.TryToAuth(&source)
// 	assert.Nil(t, err)
// }

// func TestTryToAuthFailsSetsCorrectErrorMessage(t *testing.T) {
// 	source := models.Source{
// 		DefectDojoUrl: "http://something",
// 		Username:      "admin",
// 		ApiKey:        "bad_key",
// 	}

// 	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "hello client")
// 	}))
// 	defer mock_server.Close()

// 	err := client.tryToAuth(&source)
// 	assert.NotNil(t, err)
// }
