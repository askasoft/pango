package freshdesk

// ---------------------------------------------------
// Automation

type ListAutomationRulesOption = PageOption

func (fd *Freshdesk) ListAutomationRules(aType AutomationType, laro *ListAutomationRulesOption) ([]*AutomationRule, bool, error) {
	url := fd.endpoint("/automations/%d/rules", aType)
	rules := []*AutomationRule{}
	next, err := fd.doList(url, laro, &rules)
	return rules, next, err
}

func (fd *Freshdesk) IterAutomationRules(aType AutomationType, laro *ListAutomationRulesOption, iarf func(*AutomationRule) error) error {
	if laro == nil {
		laro = &ListAutomationRulesOption{}
	}
	if laro.Page < 1 {
		laro.Page = 1
	}
	if laro.PerPage < 1 {
		laro.PerPage = 100
	}

	for {
		ars, next, err := fd.ListAutomationRules(aType, laro)
		if err != nil {
			return err
		}
		for _, ar := range ars {
			if err = iarf(ar); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		laro.Page++
	}
	return nil
}

func (fd *Freshdesk) GetAutomationRule(aType AutomationType, rid int64) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	rule := &AutomationRule{}
	err := fd.doGet(url, rule)
	return rule, err
}

func (fd *Freshdesk) DeleteAutomationRule(aType AutomationType, rid int64) error {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateAutomationRule(aType AutomationType, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules", aType)
	result := &AutomationRule{}
	if err := fd.doPost(url, rule, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateAutomationRule(aType AutomationType, rid int64, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	result := &AutomationRule{}
	if err := fd.doPut(url, rule, result); err != nil {
		return nil, err
	}
	return result, nil
}
