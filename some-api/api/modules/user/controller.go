package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pep-api/api/tools/router"
	"pep-api/api/tools/validate"
	db "pep-api/db/schemes"
	"strconv"
)

const DominioApiMemed = ""
const TokenApiMemed = "$2y$10$qvIv80NNea9NQLbEp1b81.drn5sOw01NpkzCVqEQvWpRvKgkgWBzu"
const SecretApiMemed = "$2y$10$Tk8ujHtCWy3Zvl2sEvQ5O.1p30ci9kI9bYEQjpnqkyh1Wptok/NNa"
const ResToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwibG9naW4iOiJkb2UiLCJyb2xlcyI6WyJBRE1JTiJdfQ.B83AXv3yQ80Rca0_cZlQML63yEqo-afNzwUHgIDhucI"
const ResApiURL = "http://respshml.qualirede.com.br/v1"
const Bearer = "Bearer " + ResToken

// getProfissionalByMpi(mpi: string): Observable<Profissional> {
//     const url = `${this.baseURLRES}/profissionais/mpi=${mpi}`
//     return this.http.get<any>(url, this.optionsRES).pipe(map(r => (r ? r[0] : undefined)))
//   }

// getLocalById(id: string): Observable<LocalAtendimento> {
//     const url = `${this.baseURLRES}/locaisAtendimento/${id}`
//     return this.http.get<LocalAtendimento>(url, this.optionsRES).pipe(
//       map(r => (r ? { ...r, _id: id } : undefined)),
//       switchMap(r => this.getMunicipio(r))
//     )
//   }

//   getLocalByCNES(cnes: string): Observable<LocalAtendimento> {
//     const url = `${this.baseURLRES}/locaisAtendimento/cnes?=${cnes}`
//     return this.http.get<LocalAtendimento>(url, this.optionsRES).pipe(
//       map(r => (r ? r[0] : undefined)),
//       switchMap(r => this.getMunicipio(r))
//     )

// get(mpi: string): Observable<Models.DadosPessoais> {
//     const url = `${this.base}/pessoas/${mpi}`
//     return this.http.get<any>(url, this.options)
//   }
type (
	user struct {
		ID               int     `json:"id" validate:"required"`
		SsoID            int     `json:"id_sso" validate:"required"`
		Conselho         string  `json:"conselho,omitempty"`
		ConselhoCode     string  `json:"conselho_code,omitempty"`
		UF               string  `json:"uf,omitempty"`
		TokenMemed       string  `json:"token_memed,omitempty"`
		ProMPI           string  `json:"res_pro_mpi" validate:"required"`
		ProfissionalName string  `json:"res_pro_name" validate:"required"`
		CreatedAt        int     `json:"created_at" validate:"required"`
		LastUpdatedAt    int     `json:"updated_at" validate:"required"`
		CreatedBy        int     `json:"created_by" validate:"required"`
		LastUpdatedBy    int     `json:"updated_by" validate:"required"`
		IsActive         bool    `json:"active" validate:"required"`
		Deleted          bool    `json:"deleted" validate:"required"`
		Locals           []local `json:"locals" valitade:"requried"`
	}

	local struct {
		ID      int    `json:"id,omitempty"`
		ResID   string `json:"res_id" validate:"required"`
		ResCNES string `json:"res_cnes" validate:"required"`
		ResName string `json:"res_name" validate:"required"`
		Address string `json:"address" validate:"required"`
	}

	post struct {
		SsoID            int      `json:"id_sso" validate:"required"`
		ProMPI           string   `json:"res_pro_mpi" validate:"required"`
		Conselho         string   `json:"conselho,omitempty"`
		ConselhoCode     string   `json:"conselho_code,omitempty"`
		UF               string   `json:"uf,omitempty"`
		ProfissionalName string   `json:"res_pro_name" validate:"required"`
		Locals           []local  `json:"locals" validate:"required"`
		LocalModified    []string `json:"local_id,omitempty"`
	}

	put struct {
		ProMPI           string   `json:"res_pro_mpi,omitempty"`
		ProfissionalName string   `json:"res_pro_name,omitempty"`
		Conselho         string   `json:"conselho,omitempty"`
		ConselhoCode     string   `json:"conselho_code,omitempty"`
		UF               string   `json:"uf,omitempty"`
		Locals           []local  `json:"locals,omitempty"`
		LocalModified    []string `json:"local_id,omitempty"`
		IsActive         int      `json:"active"`
		Deleted          int      `json:"deleted"`
	}

	fromMed struct {
		ResProfissionaMPI  string `json:"res_pro_mpi" validate:"required"`
		ResBeneficiarioMPI string `json:"res_benf_mpi" validate:"required"`
		ResLocalCNES       string `json:"res_local_cnes" validate:"required"`
	}

	fromMedResponse struct {
		ResToken        string // TODO tirar dúvida com eduardo Biss
		ResProfissional interface{}
		ResLocal        interface{}
		ResBeneficiario interface{}
	}
)

var messages = map[string]string{
	"user-already-exists":     "Este email já está em uso por outro usuário",
	"unexpected-create":       "Não foi possível criar este usuário",
	"unexpected-create-local": "Não foi possível criar este local",
	"unexpected-update":       "Não foi possível atualizar este usuário",
	"unexpected-read-all":     "Não foi possível acessar este recurso",
	"unexpected-new-pass":     "Erro ao atribuir nova senha",
	"success-update":          "Usuário atualizado com sucesso",
	"success-create":          "Usuário criado com sucesso",
	"success-get":             "Usuários buscados.",
	"not-found":               "Usuário não existe",
	"local-not-found":         "Usuário não pertence a este local",
}

// Create -
func Create(ctx *router.Context) (int, *router.Response) {
	var user post
	b, _ := ioutil.ReadAll(ctx.Body)
	defer ctx.Body.Close()
	json.Unmarshal(b, &user)

	errMsg := validate.NewValidator(user)
	if errMsg != "" {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: errMsg, Error: "unexpected-create"}
	}

	ok, err := db.ExistsUser(map[string]interface{}{"id_sso": user.SsoID})

	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-create"], Error: "unexpected-create"}
	}

	if ok {
		return http.StatusConflict, &router.Response{Status: "error", Message: messages["user-already-exists"], Error: "user-already-exists"}
	}

	userLocals := put{
		LocalModified: user.LocalModified,
		Locals:        user.Locals,
	}
	err = verifyLocals(&userLocals)
	user.LocalModified = userLocals.LocalModified
	user.Locals = userLocals.Locals

	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-create-local"], Error: "unexpected-create-local"}
	}

	res, err := db.AddUser(user)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-create"], Error: "unexpected-create"}
	}

	return http.StatusOK, &router.Response{Status: "OK", Data: res, Message: messages["success-create"]} // TODO ajustar
}

func verifyLocals(user *put) error {

	for _, local := range user.Locals {
		existLocal, err := db.ExistsLocal(map[string]interface{}{"res_cnes": local.ResCNES})
		if err != nil {
			return err
		}

		if !existLocal {
			res, err := db.AddLocal(local)
			if err != nil {
				return err
			}
			user.LocalModified = append(user.LocalModified, res.(map[string]interface{})["id"].(string))
		}
		user.LocalModified = append(user.LocalModified, "(SELECT id FROM local WHERE res_cnes = '"+local.ResCNES+"')")
		user.Locals = nil
	}
	return nil
}

// ReadAll -
func ReadAll(ctx *router.Context) (int, *router.Response) {
	// TODO implement filter ssoID
	var paginator = make(map[string]interface{})
	var where = make(map[string]interface{})
	var users []interface{}
	var err error

	if ctx.Queries["limit"] != nil {
		paginator["limit"] = ctx.Queries["limit"]
	}
	if ctx.Queries["offset"] != nil {
		paginator["offset"] = ctx.Queries["offset"]
	}
	where["1"] = "1"
	if ctx.Queries["ssoid"] != nil {
		where["ssoid"] = ctx.Queries["ssoid"]
	}

	users, err = db.GetUsers(where, paginator)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}
	}
	return http.StatusOK, &router.Response{Status: "OK", Data: users, Message: messages["success-get"]}
}

// Read -
func Read(ctx *router.Context) (int, *router.Response) {

	// user, err := db.GetUserByEmailAndPass(ctx.Params["email"].(string), ctx.Params["password"].(string))
	intID, _ := strconv.Atoi(ctx.Params["id"].(string))
	usuario, err := db.GetUserByID(intID)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}
	}

	if usuario == nil || len(usuario.([]interface{})) <= 0 {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["not-found"], Error: "not-found"}
	}

	u, ok := usuario.(user)
	if !ok {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}
	}
	client := &http.Client{}
	var req *http.Request
	var response *http.Response
	var responseMemed interface{}

	urlRequest := `https://` + DominioApiMemed + `/v1/sinapse-prescricao/usuarios/` + u.ConselhoCode + u.UF + `?api-key=` + TokenApiMemed + `&secret-key=` + SecretApiMemed
	req, _ = http.NewRequest("GET", urlRequest, nil)
	req.Header.Add("Accept", "application/vnd.api+json")
	response, err = client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}
	}
	json.NewDecoder(response.Body).Decode(&responseMemed)
	u.TokenMemed = responseMemed.(map[string]interface{})["data"].(map[string]interface{})["attributes"].(map[string]interface{})["token"].(string)

	return http.StatusOK, &router.Response{Status: "OK", Data: u, Message: messages["success-get"]}
}

// ReadFromMed -
func ReadFromMed(ctx *router.Context) (int, *router.Response) {
	var u fromMed
	var fullResponse fromMedResponse
	b, _ := ioutil.ReadAll(ctx.Body)
	defer ctx.Body.Close()
	json.Unmarshal(b, &u)

	errMsg := validate.NewValidator(u)
	if errMsg != "" {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: errMsg, Error: "unexpected-create"}
	}

	// 1. Verificar usuário no banco: select users where proMPI = proMPI recebido
	users, err := db.GetUsers(map[string]interface{}{"res_pro_mpi": u.ResProfissionaMPI})
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}
	}

	// 1.1 se não tiver, erro: usuário não encontrado
	if users == nil || len(users) <= 0 {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["not-found"], Error: "not-found"}
	}

	client := &http.Client{}
	var req *http.Request
	var response *http.Response

	code, apiResponse, err := getProfissional(client, req, response, &fullResponse, u)
	if err != nil {
		return code, apiResponse
	}

	// 3. Verificar local
	// 3.1 Verifica se usuário tem permissão de local (select user_local where id usuario & local cnes )
	userHasLocal, err := db.UserHasLocal(map[string]interface{}{"id_usuario": users[0].(map[string]interface{})["id"].(string), "local_cnes": u.ResLocalCNES})
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}
	}

	// 3.2 Usuário não tem autorização para este local
	if !userHasLocal {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["local-not-found"], Error: "local-not-found"}
	}

	code, apiResponse, err = getLocal(client, req, response, &fullResponse, u)
	if err != nil {
		return code, apiResponse
	}

	code, apiResponse, err = getBeneficiario(client, req, response, &fullResponse, u)
	if err != nil {
		return code, apiResponse
	}
	// retorna tudo

	return http.StatusOK, &router.Response{Status: "OK", Data: fullResponse}
}

func getProfissional(client *http.Client, req *http.Request, response *http.Response, fullResponse *fromMedResponse, body fromMed) (int, *router.Response, error) {
	var err error
	// 2. Verifica se existe profissional no res
	// 2.1 request get res/pro...
	// const url = `${this.baseURLRES}/profissionais/mpi=${mpi}`
	req, _ = http.NewRequest("GET", ResApiURL+"/profissionais/mpi="+body.ResProfissionaMPI, nil)
	req.Header.Add("Authorization", Bearer)
	response, err = client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}, err
	}
	//TODO 2.2. Erro, profissional não encontrado,
	json.NewDecoder(response.Body).Decode(&fullResponse.ResProfissional)
	return 0, nil, nil
}

func getLocal(client *http.Client, req *http.Request, response *http.Response, fullResponse *fromMedResponse, body fromMed) (int, *router.Response, error) {
	var err error
	// 4. Pegar o local do res
	// 4.1 request get res/local...
	// const url = `${this.baseURLRES}/locaisAtendimento/cnes?=${cnes}`
	req, _ = http.NewRequest("GET", ResApiURL+"/locaisAtendimento/cnes?="+body.ResLocalCNES, nil)
	req.Header.Add("Authorization", Bearer)
	response, err = client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}, err
	}
	json.NewDecoder(response.Body).Decode(&fullResponse.ResLocal)
	//TODO 4.2 Este não foi registrado

	return 0, nil, nil
}

func getBeneficiario(client *http.Client, req *http.Request, response *http.Response, fullResponse *fromMedResponse, body fromMed) (int, *router.Response, error) {
	var err error
	// 5. request get res/benf
	// const url = `${this.base}/pessoas/${mpi}`
	req, _ = http.NewRequest("GET", ResApiURL+"/pessoas/"+body.ResBeneficiarioMPI, nil)
	req.Header.Add("Authorization", Bearer)
	response, err = client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-read-all"], Error: "unexpected-read-all"}, err
	}
	json.NewDecoder(response.Body).Decode(&fullResponse.ResBeneficiario)

	return 0, nil, nil
}

// Update -
func Update(ctx *router.Context) (int, *router.Response) {
	id := ctx.Params["id"].(int)

	var u put
	b, _ := ioutil.ReadAll(ctx.Body)
	defer ctx.Body.Close()
	json.Unmarshal(b, &u)

	errMsg := validate.NewValidator(u)
	if errMsg != "" {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: errMsg, Error: "unexpected-create"}
	}
	if ctx.ExtraBody != nil {
		extraInfo := ctx.ExtraBody.(put)

		if extraInfo.IsActive != u.IsActive {
			u.IsActive = extraInfo.IsActive
		}
	}

	err := verifyLocals(&u)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-create-local"], Error: "unexpected-create-local"}
	}

	res, err := db.UpdateUser(id, u)
	if err != nil {
		return http.StatusInternalServerError, &router.Response{Status: "error", Message: messages["unexpected-update"], Error: "unexpected-update"}
	}

	return http.StatusOK, &router.Response{Status: "OK", Data: res.(user), Message: messages["success-update"]}
}

// Delete -
func Delete(ctx *router.Context) (int, *router.Response) {
	ctx.ExtraBody = put{Deleted: 1, IsActive: 0}
	return Update(ctx)
}
