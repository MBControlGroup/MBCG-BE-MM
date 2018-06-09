package entities

import (
    "time"
)

type ResSoldiersData struct {
    Soldier_id      int `json:"soldier_id"`
    Rank            string `json:"rank"`
    Id_num          string `json:"id_num"`
    Name            string `json:"name"`
    Phone_num       string `json:"phone_num"`
    Wechat_openid   string `json:"wechat_openid"`
    Commander_id    int `json:"commander_id"`
    Serve_office_id int `json:"serve_office_id"`
    Im_user_id      int `json:"im_user_id"`
    Rec_content     string `json:"rec_content"`
}

type OrgSoldierRelationships struct {
    Osr_id          int   
    Serve_org_id    int
    Soldier_id      int
}

type PageInfo struct {
    Page_data_count int `json:"pc"`
    Page_num        int `json:"pn"`
}

type SendBMData struct {
    Title       string `json:"title"`
    Detail      string `json:"detail"`
    Bm_type     string `json:"bm_type"`
    Wechat_notice   bool `json:"wechat_notice"`
    Sms_notice      bool `json:"sms_notice"`
    Voice_notice    bool `json:"voice_notice"`
}

type ResBMData struct {
    Bm_id       int `json:"bm_id"`
    Title       string `json:"title"`
    Detail      string `json:"detail"`
    Bm_type     string `json:"bm_type"`
    Send_method []string `json:"send_method"`
    Send_time   *time.Time `json:"send_time"`
}   

func NewResBMData(value BroadcastMessages) *ResBMData {
    var temp ResBMData
    temp.Bm_id = value.Bm_id
    temp.Title = value.Title
    temp.Detail = value.Detail
    temp.Bm_type = value.Bm_type
    temp.Send_time = value.Send_time
    if value.Wechat_notice {
        temp.Send_method = append(temp.Send_method, "微信")
    }
    if value.Sms_notice {
        temp.Send_method = append(temp.Send_method, "短信")
    }
    if value.Voice_notice {
        temp.Send_method = append(temp.Send_method, "语音")
    }
    return &temp
}

func NewResSD(s Soldiers, c CmNtReceipts) *ResSoldiersData {
    var resSD ResSoldiersData
    resSD.Soldier_id = s.Soldier_id
    resSD.Rank = s.Rank
    resSD.Id_num = s.Id_num
    resSD.Name = s.Name
    resSD.Phone_num = s.Phone_num
    resSD.Wechat_openid = s.Wechat_openid
    resSD.Commander_id = s.Commander_id
    resSD.Serve_office_id = s.Serve_office_id
    resSD.Im_user_id = s.Im_user_id
    resSD.Rec_content = c.Rec_content
    return &resSD
}

type Soldiers struct {
    Soldier_id      int   
    Rank            string
    Id_num          string
    Name            string
    Phone_num       string
    Wechat_openid   string
    Commander_id    int
    Serve_office_id int
    Im_user_id      int
}

type BroadcastMessages struct {
    Bm_id           int   
    Title           string
    Detail          string
    Send_time       *time.Time
    Bm_type         string
    Wechat_notice   bool
    Sms_notice      bool
    Voice_notice    bool
}

type BcMsgOffices struct {
    Bmo_id          int
    Msg_id          int
    Msg_office_id   int
}

type BcMsgOrgs struct {
    Bmo_id          int   
    Msg_id          int
    Msg_org_id      int
}

type Organizations struct {
    Org_id          int   
    Serve_office_id int
    Name            string
}

type Offices struct {
    Office_id           int   
    Office_level        string
    Higher_office_id    int
    Name                string
}

type CommonNotifications struct {
    Cn_id           int   
    Cn_bm_id        int
    Recv_soldier_id int
}

type CmNtReceipts struct {
    Cnr_id          int   
    Cn_id           int
    Rec_content     string
}

