package service

import (
	"io/ioutil"
	"strings"
	"encoding/json"
	"github.com/gorilla/mux"
    "net/http"
    "strconv"
    "fmt"
    "github.com/MBControlGroup/MBCG-BE-MM/entities"
    "github.com/MBControlGroup/MBCG-BE-MM/token"
    "github.com/unrolled/render"
    //"github.com/dgrijalva/jwt-go/request"
    //"github.com/dgrijalva/jwt-go"

)

func getAllBMsHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") 

        cookie, err := req.Cookie("token")
	    if err != nil || cookie.Value == ""{
	        formatter.JSON(w, 302, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{302, "fail", "失败", nil})
	        return;
	    }

	    user_id, err := token.Valid(cookie.Value)

	    if err != nil {
	        formatter.JSON(w, 302, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{302, "fail", "失败", nil})
	        return;
        }

        fmt.Println(user_id)

        var page_info entities.PageInfo
        page_info.Page_data_count = -1
        page_info.Page_num = -1
        err = json.NewDecoder(req.Body).Decode(&page_info)
        checkErr(err)

        if page_info.Page_data_count == -1 || page_info.Page_num == -1 {
            formatter.JSON(w, 400, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{400, "fail", "失败", nil})
	        return;
        }
 
        BMList := entities.MMService.GetAllBM()

        var ret []entities.ResBMData

        begin := page_info.Page_data_count*(page_info.Page_num-1)
        for key, value := range BMList {
            if(key >= begin&&key <= begin+page_info.Page_data_count-1) {
                temp := entities.NewResBMData(value)
                ret = append(ret, *temp)
            }
        }

        count_data := len(BMList)
        count_page := count_data/(page_info.Page_data_count)
        if count_data%(page_info.Page_data_count) != 0 {
            count_page = count_page + 1
        }

        formatter.JSON(w, http.StatusOK, struct{ 
            Code int `json:"code"`;
            Enmsg string `json:"enmsg"`;
            Cnmsg string `json:"cnmsg"`; 
            Data struct{Count_page int `json:"count_page"`; Count_data int `json:"count_data"`; Data []entities.ResBMData `json:"data"`} `json:"data"`
            }{
            200, 
            "ok", 
            "成功", 
            struct{Count_page int `json:"count_page"`; Count_data int `json:"count_data"`; Data []entities.ResBMData `json:"data"`}{count_page, count_data, ret}})
            
    }
}

func sendShortMessage(sendBMData entities.SendBMData, soldiers []entities.Soldiers) {
    client := &http.Client{}
    var var4 entities.Vars4Template
    var4.TemplateNum = sendBMData.TemplateNum
    var4.Var1 = sendBMData.Var1
    var4.Var2 = sendBMData.Var2
    var4.Var3 = sendBMData.Var3
    var4.Var4 = sendBMData.Var4
    for _, soldiers := range soldiers {
        var4.Num = soldiers.Phone_num
        myJson, err := json.Marshal(var4)
        checkErr(err)
        request, err := http.NewRequest("POST", "http://127.0.0.1:9400/sendInterfaceTemplateSms",strings.NewReader(string(myJson)))
        checkErr(err)
        defer request.Body.Close()

        response, err := client.Do(request)
        checkErr(err)
        defer response.Body.Close()

        body, err := ioutil.ReadAll(response.Body)

        fmt.Println(string(body))
        checkErr(err)
    }
    
}

func sendWebcall(soldiers []entities.Soldiers) {
    client := &http.Client{}
    var webcall entities.WebCallInfo
    for _, soldiers := range soldiers {
        webcall.Exten = soldiers.Phone_num

        myJson, err := json.Marshal(webcall)
        checkErr(err)
        request, err := http.NewRequest("POST", "http://127.0.0.1:9400/webCall",strings.NewReader(string(myJson)))
        checkErr(err)
        defer request.Body.Close()

        response, err := client.Do(request)
        checkErr(err)
        defer response.Body.Close()

        body, err := ioutil.ReadAll(response.Body)

        fmt.Println(string(body))
        checkErr(err)
    }
}

func sendBMsHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") 

        cookie, err := req.Cookie("token")
	    if err != nil || cookie.Value == ""{
	        formatter.JSON(w, 302, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{302, "fail", "失败", nil})
	        return;
	    }

	    user_id, err := token.Valid(cookie.Value)

	    if err != nil {
	        formatter.JSON(w, 302, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{302, "fail", "失败", nil})
	        return;
        }

        fmt.Println(user_id)

        var sendBMData entities.SendBMData
        err = json.NewDecoder(req.Body).Decode(&sendBMData)
        checkErr(err)

        
        bm := entities.MMService.AddBM(&sendBMData)

        req.ParseForm()

        var office_soldiers []entities.Soldiers
        if len(req.Form["office_id"]) != 0 {
            msg_office_id, err := strconv.Atoi(req.Form["office_id"][0])
            checkErr(err)

            if(!entities.MMService.ValidOfficeId(msg_office_id)) {
                formatter.JSON(w, 400, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{400, "fail", "失败", nil})
	            return;
            }

            var bmo entities.BcMsgOffices
            bmo.Msg_id = bm.Bm_id
            bmo.Msg_office_id = msg_office_id
            entities.MMService.AddBMO(&bmo)

            office_id, err := strconv.Atoi(req.Form["office_id"][0])
            checkErr(err)
            office_soldiers = entities.MMService.GetSoldierByOfficeId(office_id)
            fmt.Println(bm.Bm_id)
            for _, soldier := range office_soldiers {
                fmt.Println(soldier)
                var cn entities.CommonNotifications
                cn.Cn_bm_id = bm.Bm_id
                cn.Recv_soldier_id = soldier.Soldier_id
                cn = *entities.MMService.AddCN(&cn)
                var cnr entities.CmNtReceipts
                cnr.Cn_id = cn.Cn_id
                entities.MMService.AddCNR(&cnr)
            }
            sendShortMessage(sendBMData, office_soldiers)
            sendWebcall(office_soldiers)
        }

        if len(req.Form["org_id"]) != 0 {
            msg_org_id, err := strconv.Atoi(req.Form["org_id"][0])
            checkErr(err)

            if(!entities.MMService.ValidOrgId(msg_org_id)) {
                formatter.JSON(w, 400, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{400, "fail", "失败", nil})
	            return;
            }

            var bmorg entities.BcMsgOrgs
            bmorg.Msg_id = bm.Bm_id
            bmorg.Msg_org_id = msg_org_id
            entities.MMService.AddBMOrg(&bmorg)

            var org_soldiers []entities.Soldiers

            org_id, err := strconv.Atoi(req.Form["org_id"][0])
            checkErr(err)
            org_soldiers = entities.MMService.GetSoldierByOrgId(org_id)
            for _, soldier := range org_soldiers {
                find := false
                for _, soldierAdded := range office_soldiers {
                    if soldier.Soldier_id == soldierAdded.Soldier_id {
                        find = true
                        break
                    }
                }
                if find {
                    continue
                }
                var cn entities.CommonNotifications
                cn.Cn_bm_id = bm.Bm_id
                cn.Recv_soldier_id = soldier.Soldier_id
                cn = *entities.MMService.AddCN(&cn)
                var cnr entities.CmNtReceipts
                cnr.Cn_id = cn.Cn_id
                entities.MMService.AddCNR(&cnr)
            }
            sendShortMessage(sendBMData, org_soldiers)
            sendWebcall(org_soldiers)
        }
        
        resBMData := entities.NewResBMData(*bm)

        formatter.JSON(w, http.StatusOK, struct{ 
            Code int `json:"code"`;
            Enmsg string `json:"enmsg"`;
            Cnmsg string `json:"cnmsg"`; 
            Data entities.ResBMData `json:"data"`
            }{
            200, "ok", "成功", *resBMData})
        
        //formatter.JSON(w, http.StatusOK, task)
    }
}

func getBMHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") 
        
        cookie, err := req.Cookie("token")
	    if err != nil || cookie.Value == ""{
	        formatter.JSON(w, 302, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{302, "fail", "失败", nil})
	        return;
	    }

	    user_id, err := token.Valid(cookie.Value)

	    if err != nil {
	        formatter.JSON(w, 302, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{302, "fail", "失败", nil})
	        return;
        }

        fmt.Println(user_id)

        var page_info entities.PageInfo
        page_info.Page_data_count = -1
        page_info.Page_num = -1
        err = json.NewDecoder(req.Body).Decode(&page_info)
        checkErr(err)

        if page_info.Page_data_count == -1 || page_info.Page_num == -1 {
            formatter.JSON(w, 400, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{400, "fail", "失败", nil})
	        return;
        }

        vars := mux.Vars(req)
        bm_id, _ := strconv.Atoi(vars["bm_id"])

        BM := entities.MMService.GetBMById(bm_id)
        if BM == nil {
            formatter.JSON(w, 400, struct{ Code int `json:"code"`;Enmsg string `json:"enmsg"`;Cnmsg string `json:"cnmsg"`; Data interface{} `json:"data"`}{400, "fail", "失败", nil})
	        return;
        }

        resBMData := entities.NewResBMData(*BM)

        resSoldierData := entities.MMService.GetResSoldierDataByBMId(bm_id)

        var ret []entities.ResSoldiersData

        begin := page_info.Page_data_count*(page_info.Page_num-1)
        for key, value := range resSoldierData {
            if(key >= begin&&key <= begin+page_info.Page_data_count-1) {
                ret = append(ret, value)
            }
        }

        count_data := len(resSoldierData)
        count_page := count_data/page_info.Page_data_count
        if count_data%(page_info.Page_data_count) != 0 {
            count_page = count_page + 1
        }

        formatter.JSON(w, http.StatusOK, struct{ 
            Code int `json:"code"`;
            Enmsg string `json:"enmsg"`;
            Cnmsg string `json:"cnmsg"`; 
            Data struct{Count_page int `json:"count_page"`; Count_data int `json:"count_data"`; Data entities.ResBMData `json:"data"`; Soldiers []entities.ResSoldiersData `json:"soldiers"`} `json:"data"`
            }{
            200, "ok", "成功", 
            struct{Count_page int `json:"count_page"`; Count_data int `json:"count_data"`; Data entities.ResBMData `json:"data"`; Soldiers []entities.ResSoldiersData `json:"soldiers"`}{
            count_page, count_data, *resBMData, ret}})
    }
}