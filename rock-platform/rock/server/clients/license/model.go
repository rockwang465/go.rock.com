package license

//
//type CAStatus struct {
//	Disable    bool     `json:"disable,omitempty"`
//	IsActive   bool     `json:"is_active,omitempty"`
//	AloneTime  int32    `json:"alone_time,omitempty"`
//	Status     string   `json:"status,omitempty"`
//	Cluster    string   `json:"cluster,omitempty"`
//	UpdatedAt  string   `json:"updated_at,omitempty"`
//	CaJson     string   `json:"ca_json,omitempty"`
//	Product    string   `json:"product,omitempty"`
//	DongleId   string   `json:"dongle_id,omitempty"`
//	ExpiredAt  string   `json:"expired_at,omitempty"`
//	DongleTime string   `json:"dongle_time,omitempty"`
//	Company    string   `json:"company,omitempty"`
//	Slaves     []string `json:"slaves,omitempty"`
//	FeatureIds []string `json:"feature_ids,omitempty"`
//}
//
//// C2v content
//type C2vResp struct {
//	C2V string `json:"c2v"`
//}
//
//// Fingerprint content
//type FingerprintResp struct {
//	FingerPrint string `json:"finger_print"`
//}
//
//// Online active request body
//type onlineActiveReq struct {
//	Action string `json:"action"`
//}
//
//// Offline active request body
//type offlineActiveReq struct {
//	V2C string `json:"v2c"`
//}
//
//// ActiveResp is status from ca
//type ActiveResp struct {
//	StatusCode    string `json:"status_code"`
//	StatusMessage string `json:"status_message"`
//}
//
//// Get client licenses resp
//type ClientLicResp struct {
//	Licenses []string `json:"licenses"`
//}
