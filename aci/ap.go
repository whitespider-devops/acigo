package aci

import (
	"bytes"
	"fmt"
)

func jsonAPAdd(tenant, name, descr string) string {
	return jsonAP(tenant, name, "created", descr)
}

func jsonAPDel(tenant, name string) string {
	return jsonAP(tenant, name, "deleted", "")
}

func rnAP(ap string) string {
	return "ap-" + ap
}

func dnAP(tenant, ap string) string {
	return rnTenant(tenant) + "/" + rnAP(ap)
}

func jsonAP(tenant, name, action, descr string) string {
	ap := rnAP(name)
	dn := dnAP(tenant, name)

	prefix := fmt.Sprintf(`{"fvAp":{"attributes":{"dn":"uni/%s","name":"%s"`, dn, name)

	var mid string
	if descr != "" {
		mid = fmt.Sprintf(`,"descr":"%s"`, descr)
	}

	suffix := fmt.Sprintf(`,"rn":"%s","status":"%s"}}}`, ap, action)

	return prefix + mid + suffix
}

func apiAP(tenant, name string) string {
	return "/api/node/mo/uni/" + dnAP(tenant, name) + ".json"
}

// ApplicationProfileAdd creates a new application profile in a tenant.
func (c *Client) ApplicationProfileAdd(tenant, name, descr string) ([]map[string]interface{}, error) {

	api := apiAP(tenant, name)

	key := "fvAp"

	j := jsonAPAdd(tenant, name, descr)

	url := c.getURL(api)

	url += "?rsp-subtree=modified" // demand response

	c.debugf("ApplicationProfileAdd: url=%s json=%s", url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return nil, errPost
	}

	c.debugf("ApplicationProfileAdd: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "ApplicationProfileAdd")
}

// ApplicationProfileDel deletes an existing application profile from a tenant.
func (c *Client) ApplicationProfileDel(tenant, name string) ([]map[string]interface{}, error) {

	api := apiAP(tenant, name)

	key := "fvAp"

	j := jsonAPDel(tenant, name)

	url := c.getURL(api)

	c.debugf("ApplicationProfileAdd: url=%s json=%s", url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return nil, errPost
	}

	c.debugf("ApplicationProfileDel: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "ApplicationProfileDel")
}

// ApplicationProfileList retrieves application profiles from a tenant.
func (c *Client) ApplicationProfileList(tenant string) ([]map[string]interface{}, error) {

	key := "fvAp"

	api := "/api/node/mo/uni/tn-" + tenant + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("ApplicationProfileList: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, errGet
	}

	c.debugf("ApplicationProfileList: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "ApplicationProfileList")
}
