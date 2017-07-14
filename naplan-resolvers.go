package naprrql

import (
	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	resolvers["NaplanData/score_summaries"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjectForKeys(getObjectReferences("NAPTestScoreSummary"))
	}

	resolvers["NaplanData/students"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjectForKeys(getObjectReferences("StudentPersonal"))
	}

	resolvers["RegistrationRecord/OtherIdList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		otherIDs := []interface{}{}
		if napRegistrationRecord, ok := params.Source.(xml.RegistrationRecord); ok {
			return napRegistrationRecord.OtherIdList.OtherId, nil
		}
		return otherIDs, nil
	}

	resolvers["NaplanData/events"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return getObjectForKeys(getObjectReferences("NAPEventStudentLink"))
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

	return resolvers
}
