package entities

import (
	"time"
)

//MMAtomicService .
type MMAtomicService struct{}

//UserInfoService .
var MMService = MMAtomicService{}

func (*MMAtomicService) SoldierSave(s *Soldiers) error {
    engine.Insert(s)

    return nil
}

func (*MMAtomicService) GetSoldierByOfficeId(officd_id int) []Soldiers {
    sql := "select * from Soldiers where serve_office_id = ?"
    var soldiers []Soldiers
    err := engine.SQL(sql, officd_id).Find(&soldiers) 
    checkErr(err)
    return soldiers
}

func (*MMAtomicService) GetSoldierByOrgId(org_id int) []Soldiers {
    sql := "select * from OrgSoldierRelationships where serve_org_id = ?"
    var orgsoldierrelationships []OrgSoldierRelationships
    err := engine.SQL(sql, org_id).Find(&orgsoldierrelationships)
    checkErr(err)
    sql = "select * from Soldiers where soldier_id = ?"
    var soldiers []Soldiers
    for _, value := range orgsoldierrelationships {
        var soldier Soldiers
        _, err = engine.SQL(sql, value.Soldier_id).Get(&soldier)
        checkErr(err)
        soldiers = append(soldiers, soldier)
    }
    return soldiers
}

func (*MMAtomicService) GetResSoldierDataByBMId(bm_id int) []ResSoldiersData {
    sql := "select * from CommonNotifications where cn_bm_id = ?"
    var commonnotifications []CommonNotifications
    err := engine.SQL(sql, bm_id).Find(&commonnotifications)
    checkErr(err)
    sql = "select * from Soldiers where soldier_id = ?"
    sql1 := "select * from CmNtReceipts where cn_id = ?"
    var resSDs []ResSoldiersData
    for _, value := range commonnotifications {
        var soldier Soldiers
        _, err = engine.SQL(sql, value.Recv_soldier_id).Get(&soldier)
        checkErr(err)
        var cnr CmNtReceipts
        _, err = engine.SQL(sql1, value.Cn_id).Get(&cnr)
        checkErr(err)
        resSD := NewResSD(soldier, cnr)
        resSDs = append(resSDs, *resSD)
    }
    return resSDs
}

func (*MMAtomicService) GetAllBM() []BroadcastMessages {
    sql := "select * from BroadcastMessages"
    var BMs []BroadcastMessages
    err := engine.SQL(sql).Find(&BMs)
    checkErr(err)
    return BMs
}


func (*MMAtomicService) GetBMById(bm_id int) *BroadcastMessages {
    sql := "select * from BroadcastMessages where bm_id = ?"
    var bm BroadcastMessages
    find, err := engine.SQL(sql, bm_id).Get(&bm)
    checkErr(err)
    if find {
        return &bm
    } else {
        return nil
    }
}

func (*MMAtomicService) AddBM(s *SendBMData) *BroadcastMessages {
    session := engine.NewSession()
    defer session.Close()

    err := session.Begin()
    checkErr(err)

    sql := "INSERT INTO BroadcastMessages (title, detail, send_time, bm_type, wechat_notice, sms_notice, voice_notice) VALUES (?,?,?,?,?,?,?);"

    //var ss BroadcastMessages
    res , err := engine.Exec(sql, s.Title, s.Detail, time.Now().Format("2006-01-02 15:04:05"), s.Bm_type, s.Wechat_notice, s.Sms_notice, s.Voice_notice)
    checkErr(err)

    if err == nil {
        session.Commit()
        sql = "SELECT * FROM BroadcastMessages WHERE BroadcastMessages.bm_id=?;"
        var ss BroadcastMessages
        id, err := res.LastInsertId()
        checkErr(err)
        _, err = engine.SQL(sql, id).Get(&ss)
        checkErr(err)
        return &ss
    } else {
        session.Rollback()
    }
    return nil
}

func (*MMAtomicService) AddBMO(b *BcMsgOffices)  {
    session := engine.NewSession()
    defer session.Close()

    err := session.Begin()
    checkErr(err)

    sql := "INSERT INTO BcMsgOffices (msg_id, msg_office_id) VALUES (?,?);"

    _ , err = engine.Exec(sql, b.Msg_id, b.Msg_office_id)
    checkErr(err)

    if err == nil {
        session.Commit()
    } else {
        session.Rollback()
    }
}

func (*MMAtomicService) ValidOfficeId(office_id int) bool {
    sql := "select * from Offices where office_id = ?"
    var office Offices
    find, err := engine.SQL(sql, office_id).Get(&office)
    checkErr(err)
    return find
}

func (*MMAtomicService) ValidOrgId(org_id int) bool {
    sql := "select * from Organizations where org_id = ?"
    var org Organizations
    find, err := engine.SQL(sql, org_id).Get(&org)
    checkErr(err)
    return find
}

func (*MMAtomicService) AddBMOrg(b *BcMsgOrgs) *BcMsgOrgs {
    session := engine.NewSession()
    defer session.Close()

    err := session.Begin()
    checkErr(err)

    sql := "INSERT INTO BcMsgOrgs (msg_id, msg_org_id) VALUES (?,?);"

    res, err := engine.Exec(sql, b.Msg_id, b.Msg_org_id)
    checkErr(err)

    if err == nil {
        session.Commit()
        sql = "SELECT * FROM BcMsgOrgs WHERE bmo_id =?;"
        var bb BcMsgOrgs
        id, err := res.LastInsertId()
        checkErr(err)
        _, err = engine.SQL(sql, id).Get(&bb)
        checkErr(err)
        return &bb
    } else {
        session.Rollback()
    }
    return nil
}

func (*MMAtomicService) AddCN(c *CommonNotifications) *CommonNotifications {
    session := engine.NewSession()
    defer session.Close()

    err := session.Begin()
    checkErr(err)

    sql := "INSERT INTO CommonNotifications (cn_bm_id, recv_soldier_id) VALUES (?,?);"

    res, err := engine.Exec(sql, c.Cn_bm_id, c.Recv_soldier_id)
    checkErr(err)

    if err == nil {
        session.Commit()
        sql = "SELECT * FROM CommonNotifications WHERE cn_id =?;"
        var cc CommonNotifications
        id, err := res.LastInsertId()
        checkErr(err)
        _, err = engine.SQL(sql, id).Get(&cc)
        checkErr(err)
        return &cc
    } else {
        session.Rollback()
    }
    return nil
}

func (*MMAtomicService) AddCNR(c *CmNtReceipts) *CmNtReceipts {
    session := engine.NewSession()
    defer session.Close()

    err := session.Begin()
    checkErr(err)

    sql := "INSERT INTO CmNtReceipts (cn_id, rec_content) VALUES (?,?);"

    res, err := engine.Exec(sql, c.Cn_id, "")
    checkErr(err)

    if err == nil {
        session.Commit()
        sql = "SELECT * FROM CmNtReceipts WHERE cnr_id =?;"
        var cc CmNtReceipts
        id, err := res.LastInsertId()
        checkErr(err)
        _, err = engine.SQL(sql, id).Get(&cc)
        checkErr(err)
        return &cc
    } else {
        session.Rollback()
    }
    return nil
}