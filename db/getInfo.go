package db

import "log"

type Equip struct {
	Id    int        `json:"eid"`
	Name  string     `json:"name"`
	Type  int        `json:"type"`
	Lv     int        `json:"level"`
	Color  int        `json:"color"`
	Atk    int        `json:"attack"`
	Def    int        `json:"def"`
	Holy   int        `json:"holy"`
	Str    int        `json:"str"`
	Con    int        `json:"con"`
	Hpm    int        `json:"hpmax"`
	Cri    int        `json:"cri"`
	Miss   int        `json:"miss"`
	Suit   int        `json:"suit"`
	Effect string     `json:"effect"`
	Money  int        `json:"money"`
	Descr  string     `json:"descr"`
}

func getEquipInfo(eid int) (p_equip *Equip) {
	p_equip = new(Equip)
	p_equip.Id = eid
	row := gDB.QueryRow("SELECT _name,_type,lv,color,atk,def,holy,str,con,hpm,cri,miss,suit,effect,money,description FROM s_equipments WHERE _id=?", eid)
	e := row.Scan(
		&p_equip.Name, 	&p_equip.Type,
		&p_equip.Lv, 	&p_equip.Color,
		&p_equip.Atk, 	&p_equip.Def,
		&p_equip.Holy, 	&p_equip.Str,
		&p_equip.Con, 	&p_equip.Hpm,
		&p_equip.Cri, 	&p_equip.Miss,
		&p_equip.Suit, 	&p_equip.Effect,
		&p_equip.Money, &p_equip.Descr)
	if e != nil {
		log.Fatal(e)
	}
	return
}
