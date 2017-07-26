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

	resolvers["NaplanData/tests_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTest")), nil
	}

	resolvers["NaplanData/tests"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTest"))
	}

	resolvers["NaplanData/testlets_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTestlet")), nil
	}

	resolvers["NaplanData/testlets"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTestlet"))
	}

	resolvers["NaplanData/testitems_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPTestItem")), nil
	}

	resolvers["NaplanData/testitems"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPTestItem"))
	}

	resolvers["NaplanData/codeframes_count"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return len(getIdentifiers("NAPCodeFrame")), nil
	}

	resolvers["NaplanData/codeframes"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjects(getIdentifiers("NAPCodeFrame"))
	}

	//
	// addition to spec that allows the original Item to be available when
	// reviewing item responses, e.g. to compare item correct response, item type etc.
	//
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

	return resolvers
}
