package naprrql

import (
	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	resolvers["NaplanData/score_summaries"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTestScoreSummary"))
	}

	resolvers["NaplanData/score_summaries_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTestScoreSummary")), nil
	}

	resolvers["NaplanData/students"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("StudentPersonal"))
	}

	resolvers["NaplanData/students_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("StudentPersonal")), nil
	}

	resolvers["RegistrationRecord/OtherIdList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		otherIDs := []interface{}{}
		if napRegistrationRecord, ok := params.Source.(xml.RegistrationRecord); ok {
			return napRegistrationRecord.OtherIdList.OtherId, nil
		}
		return otherIDs, nil
	}

	resolvers["NaplanData/events"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPEventStudentLink"))
	}

	resolvers["NaplanData/events_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPEventStudentLink")), nil
	}

	resolvers["NAPEvent/TestDisruptionList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		disruptionEvents := []interface{}{}
		if napEvent, ok := params.Source.(xml.NAPEvent); ok {
			return napEvent.TestDisruptionList.TestDisruption, nil
		}
		return disruptionEvents, nil
	}

	resolvers["Adjustment/PNPCodeList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		pnpCodes := []interface{}{}
		if adjustment, ok := params.Source.(xml.Adjustment); ok {
			return adjustment.PNPCodelist.PNPCode, nil
		}
		return pnpCodes, nil
	}

	resolvers["NaplanData/responses"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPStudentResponseSet"))
	}

	resolvers["NaplanData/responses_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPStudentResponseSet")), nil
	}

	resolvers["NAPResponseSet/DomainScore"] = func(params *graphql.ResolveParams) (interface{}, error) {
		domainScore := []interface{}{}
		if response, ok := params.Source.(xml.NAPResponseSet); ok {
			return response.DomainScore, nil
		}
		return domainScore, nil
	}

	resolvers["NAPResponseSet/TestletList"] = func(params *graphql.ResolveParams) (interface{}, error) {

		testletList := []interface{}{}
		if response, ok := params.Source.(xml.NAPResponseSet); ok {
			return response.TestletList.Testlet, nil
		}
		return testletList, nil
	}

	resolvers["NAPResponseSet_Testlet/ItemResponseList"] = func(params *graphql.ResolveParams) (interface{}, error) {

		itemList := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet_Testlet); ok {
			return napResponse.ItemResponseList.ItemResponse, nil
		}
		return itemList, nil

	}

	resolvers["NAPResponseSet_ItemResponse/Item"] = func(params *graphql.ResolveParams) (interface{}, error) {

		linkedItem := make([]string, 0)
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet_ItemResponse); ok {
			linkedItem = append(linkedItem, napResponse.ItemRefID)
			obj, err := getObjects(linkedItem)
			return obj[0], err
		}
		return linkedItem, nil

	}

	resolvers["NAPResponseSet_ItemResponse/SubscoreList"] = func(params *graphql.ResolveParams) (interface{}, error) {

		subscoreList := []interface{}{}
		// log.Printf("params: %#v\n\n", params)
		if napResponse, ok := params.Source.(xml.NAPResponseSet_ItemResponse); ok {
			return napResponse.SubscoreList.Subscore, nil
		}
		return subscoreList, nil

	}

	return resolvers
}
