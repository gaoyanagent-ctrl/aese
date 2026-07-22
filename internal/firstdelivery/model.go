package firstdelivery

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
	Scale    int    `json:"scale"`
}
type Shipment struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
	Lot      string `json:"lot"`
	Accepted int    `json:"accepted"`
}
type Cost struct{ Material, Labor, Energy, ScrapQuality, ExpediteLogistics, Overhead string }
type Frame struct {
	Step                       int        `json:"step"`
	Phase                      string     `json:"phase"`
	SimTime                    string     `json:"sim_time"`
	Title                      string     `json:"title"`
	CausationID                string     `json:"causation_id"`
	Demand                     int        `json:"demand"`
	Planned                    int        `json:"planned"`
	Supplied                   int        `json:"supplied"`
	Good                       int        `json:"good"`
	Scrap                      int        `json:"scrap"`
	Inventory                  int        `json:"inventory"`
	Shipped                    int        `json:"shipped"`
	Accepted                   int        `json:"accepted"`
	Shipments                  []Shipment `json:"shipments"`
	WorldProgress              int        `json:"world_progress"`
	IAOSProgress               int        `json:"iaos_progress"`
	Discrepancy                string     `json:"discrepancy"`
	Cash                       Money      `json:"cash"`
	ContractLiability          Money      `json:"contract_liability"`
	InvoiceGross               Money      `json:"invoice_gross"`
	AR                         Money      `json:"ar"`
	Collected                  Money      `json:"collected"`
	Revenue                    Money      `json:"revenue"`
	ActualCost                 Money      `json:"actual_cost"`
	GrossMargin                Money      `json:"gross_margin"`
	Knowledge                  []string   `json:"knowledge"`
	FirstCommercialCycleClosed bool       `json:"first_commercial_cycle_closed"`
	IAOSCursor                 int64      `json:"iaos_cursor"`
}
type Trace struct {
	SchemaVersion            string  `json:"schema_version"`
	Campaign                 string  `json:"campaign"`
	WorldRunID               string  `json:"world_run_id"`
	Timezone                 string  `json:"timezone"`
	PolicyVersion            string  `json:"policy_version"`
	M12TerminalHash          string  `json:"m12_terminal_hash"`
	ReleaseCompatibility     string  `json:"release_compatibility"`
	OpeningSaleableInventory int     `json:"opening_saleable_inventory"`
	Frames                   []Frame `json:"frames"`
}

func m(v string) Money { return Money{v, "CNY", 2} }
func BuildTrace() Trace {
	b := Frame{0, "eligible", "2027-06-16T09:00:00+08:00", "消费 M12 量产资格并完成财务结转", "m12-terminal", 0, 0, 0, 0, 0, 0, 0, 0, []Shipment{}, 0, 0, "none", m("9300000.00"), m("2000000.00"), m("0.00"), m("0.00"), m("0.00"), m("0.00"), m("0.00"), m("0.00"), []string{}, false, 500}
	fs := []Frame{b}
	add := func(p, title, c string, w, i int) {
		f := fs[len(fs)-1]
		f.Step = len(fs)
		f.Phase = p
		f.SimTime = time.Date(2027, time.June, 16+len(fs)*3, 10, 0, 0, 0, time.FixedZone("CST", 8*3600)).Format(time.RFC3339)
		f.Title = title
		f.CausationID = c
		f.WorldProgress = w
		f.IAOSProgress = i
		fs = append(fs, f)
	}
	add("ordered", "首张 10,000 件正式订单与 2,000 件追加要求受治理确认", "order-confirmed", 8, 8)
	fs[1].Demand = 12000
	add("planned", "ATP/MRP 基于零成品库存、发布 BOM 与能力生成", "mrp-approved", 15, 15)
	fs[2].Planned = 12000
	add("supplied", "采购、运输、IQC 与可追溯物料放行", "materials-released", 28, 28)
	fs[3].Supplied = 12360
	add("producing", "正式生产形成 12,000 良品与 360 报废", "production-completed", 52, 58)
	fs[4].Good = 12000
	fs[4].Scrap = 360
	fs[4].Inventory = 12000
	add("shipment_1", "第一批 9,000 件实际发运并由客户接受", "shipment-a-accepted", 65, 70)
	fs[5].Shipments = []Shipment{{"SHIP-GENESIS-0001-A", 9000, "accepted", "LOT-GENESIS-BCP-A01-001", 9000}}
	fs[5].Shipped = 9000
	fs[5].Accepted = 9000
	fs[5].Inventory = 3000
	add("shipment_2", "第二批 2,700 件接受，形成 300 件短缺", "shipment-b-accepted", 78, 92)
	fs[6].Shipments = append(fs[5].Shipments, Shipment{"SHIP-GENESIS-0001-B", 2700, "accepted", "LOT-GENESIS-BCP-A01-002", 2700})
	fs[6].Shipped = 11700
	fs[6].Accepted = 11700
	fs[6].Inventory = 300
	fs[6].Discrepancy = "delivery_short_300_open"
	add("recovery", "计划/质量/经营角色获知后，备选供应、维修与加班获批", "recovery-approved", 82, 92)
	fs[7].Knowledge = []string{"planning:short-300", "quality:intensified-inspection", "business:expedite-cost"}
	add("delivered", "第三批 300 件实际交付并接受，数量差异关闭", "shipment-c-accepted", 100, 100)
	fs[8].Shipments = append(fs[6].Shipments, Shipment{"SHIP-GENESIS-0001-C", 300, "accepted", "LOT-GENESIS-BCP-A01-003", 300})
	fs[8].Shipped = 12000
	fs[8].Accepted = 12000
	fs[8].Inventory = 0
	fs[8].Discrepancy = "closed"
	add("invoiced", "仅按已接受 12,000 件开票并形成应收", "invoice-issued", 100, 100)
	fs[9].InvoiceGross = m("16272000.00")
	fs[9].AR = m("14272000.00")
	fs[9].Revenue = m("14400000.00")
	fs[9].ContractLiability = m("0.00")
	add("collected", "客户银行实际到账并受治理核销应收", "bank-receipt-settled", 100, 100)
	fs[10].Collected = m("14272000.00")
	fs[10].AR = m("0.00")
	fs[10].Cash = m("14572000.00")
	add("cost_closed", "实际成本与标准差异归集，项目毛利可解释", "actual-cost-closed", 100, 100)
	fs[11].ActualCost = m("10200000.00")
	fs[11].GrossMargin = m("4200000.00")
	add("first_commercial_cycle_closed", "首张订单交付、开票、回款和成本商业闭环完成", "commercial-cycle-closed", 100, 100)
	fs[12].FirstCommercialCycleClosed = true
	return Trace{"1.0", "first-delivery", "world-run-genesis-first-delivery-001", "Asia/Shanghai", "first-delivery-policy-v1", "sha256:m12-industrialization-terminal", "hctm-o2d-stable-codes-compatible", 0, fs}
}
func Validate(t Trace) error {
	want := []string{"eligible", "ordered", "planned", "supplied", "producing", "shipment_1", "shipment_2", "recovery", "delivered", "invoiced", "collected", "cost_closed", "first_commercial_cycle_closed"}
	if t.OpeningSaleableInventory != 0 || len(t.Frames) != len(want) {
		return errors.New("opening/frame invalid")
	}
	for i, f := range t.Frames {
		if f.Phase != want[i] {
			return fmt.Errorf("phase %d", i)
		}
		if _, e := time.Parse(time.RFC3339, f.SimTime); e != nil {
			return e
		}
	}
	x := t.Frames[12]
	if !x.FirstCommercialCycleClosed || x.Accepted != 12000 || x.AR.Value != "0.00" || x.GrossMargin.Value != "4200000.00" {
		return errors.New("commercial gate incomplete")
	}
	return nil
}
func Hash(t Trace) (string, error) {
	b, e := json.Marshal(t)
	if e != nil {
		return "", e
	}
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:]), nil
}
func Snapshot(t Trace) ([]byte, error) {
	if e := Validate(t); e != nil {
		return nil, e
	}
	return json.Marshal(t)
}
func Restore(b []byte) (Trace, error) {
	var t Trace
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	e := d.Decode(&t)
	if e == nil {
		e = Validate(t)
	}
	return t, e
}
